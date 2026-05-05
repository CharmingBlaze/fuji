package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"fuji/internal/lexer"
	"fuji/internal/parser"
)

// emitObject emits LLVM IR for object expressions.
func (g *Generator) emitObject(e *parser.ObjectExpr) (value.Value, error) {
	nKeys := len(e.Keys) + len(e.ComputedKeys)
	if nKeys < 1 {
		nKeys = 1 // runtime table expects non-zero capacity for key/value slots
	}
	count := constant.NewInt(types.I32, int64(nKeys))
	obj := g.block.NewCall(g.runtimeAllocObj, count)

	// Set properties
	for i, key := range e.Keys {
		keyVal := g.emitStringLiteral(key.Lexeme)
		val, err := g.emitExpr(e.Values[i])
		if err != nil {
			return nil, err
		}
		g.block.NewCall(g.runtimeObjSet, obj, keyVal, g.emitAsFujiI64(val))
	}

	// Set computed keys (if any)
	for i, keyExpr := range e.ComputedKeys {
		keyVal, err := g.emitExpr(keyExpr)
		if err != nil {
			return nil, err
		}
		valIdx := len(e.Keys) + i
		val, err := g.emitExpr(e.Values[valIdx])
		if err != nil {
			return nil, err
		}
		g.block.NewCall(g.runtimeObjSet, obj, keyVal, g.emitAsFujiI64(val))
	}

	return obj, nil
}

// emitIndex emits LLVM IR for index expressions (property or array access).
func (g *Generator) emitIndex(e *parser.IndexExpr) (value.Value, error) {
	obj, err := g.emitExpr(e.Object)
	if err != nil {
		return nil, err
	}

	key, err := g.emitExpr(e.Index)
	if err != nil {
		return nil, err
	}

	// For now, we'll use a simple heuristic:
	// If the key is a string (represented as an object in our Value system), use object_get.
	// Otherwise, assume it's an array index.
	// In a production compiler, we'd have a unified 'FUJI_get' runtime function.

	// Check if key is a string (this is a bit hacky without type info, but works for literals)
	// For now, let's just use FUJI_object_get if the index expression was a string literal
	if _, ok := e.Index.(*parser.LiteralExpr); ok {
		// If it was a literal, we can be more certain
		// But in general, we should use a runtime helper FUJI_get(obj, key)
	}

	// Use FUJI_object_get for everything for now, or implement a basic check
	// Actually, let's use FUJI_object_get as the default and fallback to array if it fails?
	// No, the runtime should handle it. Let's use a placeholder 'FUJI_get' if it existed.
	// Since it doesn't, let's use FUJI_object_get for string keys and array_get for others.

	return g.block.NewCall(g.runtimeObjGet, g.emitAsFujiI64(obj), g.emitAsFujiI64(key)), nil
}

// emitAssign emits LLVM IR for assignment expressions.
func (g *Generator) emitAssign(e *parser.AssignExpr) (value.Value, error) {
	if e.Token.Type != lexer.TokenEqual {
		return g.emitCompoundAssign(e, e.Token.Lexeme)
	}

	val, err := g.emitExpr(e.Value)
	if err != nil {
		return nil, err
	}

	switch left := e.Left.(type) {
	case *parser.IdentifierExpr:
		name := left.Name.Lexeme
		if slot, ok := g.locals[name]; ok {
			boxed := g.emitAsFujiI64(val)
			if g.localIsCell != nil && g.localIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return val, nil
		}
		return nil, fmt.Errorf("undefined variable: %s", name)
	case *parser.IndexExpr:
		// Property or array assignment
		obj, err := g.emitExpr(left.Object)
		if err != nil {
			return nil, err
		}
		key, err := g.emitExpr(left.Index)
		if err != nil {
			return nil, err
		}
		// If the value is a double, convert it to i64 (NaN-boxed)
		if val.Type().String() == "double" {
			val = g.block.NewBitCast(val, types.I64)
		}
		return g.block.NewCall(g.runtimeSet, g.emitAsFujiI64(obj), g.emitAsFujiI64(key), g.emitAsFujiI64(val)), nil
	default:
		return nil, fmt.Errorf("unsupported assignment target: %T", left)
	}
}

// emitCompoundAssign emits LLVM IR for compound assignment expressions (+=, -=, *=, /=).
func (g *Generator) emitCompoundAssign(e *parser.AssignExpr, op string) (value.Value, error) {
	// Get current value
	current, err := g.emitExpr(e.Left)
	if err != nil {
		return nil, err
	}

	// Get new value
	newVal, err := g.emitExpr(e.Value)
	if err != nil {
		return nil, err
	}

	// Perform operation on unboxed numbers, then re-box (same semantics as emitInfix for + - * /).
	leftI := g.emitAsFujiI64(current)
	rightI := g.emitAsFujiI64(newVal)
	ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
	rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
	var result value.Value
	switch op {
	case "+=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFAdd(ld, rd))
	case "-=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFSub(ld, rd))
	case "*=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFMul(ld, rd))
	case "/=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFDiv(ld, rd))
	case "%=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFRem(ld, rd))
	default:
		return nil, fmt.Errorf("unsupported compound operator: %s", op)
	}

	// Assign result
	switch left := e.Left.(type) {
	case *parser.IdentifierExpr:
		name := left.Name.Lexeme
		if slot, ok := g.locals[name]; ok {
			boxed := g.emitAsFujiI64(result)
			if g.localIsCell != nil && g.localIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return result, nil
		}
		return nil, fmt.Errorf("undefined variable: %s", name)
	case *parser.IndexExpr:
		obj, err := g.emitExpr(left.Object)
		if err != nil {
			return nil, err
		}
		key, err := g.emitExpr(left.Index)
		if err != nil {
			return nil, err
		}
		g.block.NewCall(g.runtimeSet, g.emitAsFujiI64(obj), g.emitAsFujiI64(key), g.emitAsFujiI64(result))
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported assignment target: %T", left)
	}
}

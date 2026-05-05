package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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
	objSlot := g.entryAlloca(types.I64)
	g.block.NewStore(obj, objSlot)
	g.shadowStoreTemp(objSlot)

	// Set properties
	for i, key := range e.Keys {
		keyVal := g.emitStringLiteral(key.Lexeme)
		keySlot := g.entryAlloca(types.I64)
		g.shadowStoreTemp(keySlot)
		g.block.NewStore(g.emitAsFujiI64(keyVal), keySlot)
		val, err := g.emitExpr(e.Values[i])
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		keyLive := g.block.NewLoad(types.I64, keySlot)
		g.block.NewCall(g.runtimeObjSet, objLive, keyLive, g.emitAsFujiI64(val))
	}

	// Set computed keys (if any)
	for i, keyExpr := range e.ComputedKeys {
		keyVal, err := g.emitExpr(keyExpr)
		if err != nil {
			return nil, err
		}
		keySlot := g.entryAlloca(types.I64)
		g.shadowStoreTemp(keySlot)
		g.block.NewStore(g.emitAsFujiI64(keyVal), keySlot)
		valIdx := len(e.Keys) + i
		val, err := g.emitExpr(e.Values[valIdx])
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		keyLive := g.block.NewLoad(types.I64, keySlot)
		g.block.NewCall(g.runtimeObjSet, objLive, keyLive, g.emitAsFujiI64(val))
	}

	return g.block.NewLoad(types.I64, objSlot), nil
}

// emitIndex emits LLVM IR for index expressions (property or array access).
func (g *Generator) emitIndex(e *parser.IndexExpr) (value.Value, error) {
	obj, err := g.emitExpr(e.Object)
	if err != nil {
		return nil, err
	}
	objSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(objSlot)
	g.block.NewStore(g.emitAsFujiI64(obj), objSlot)

	if !e.Optional {
		key, err := g.emitExpr(e.Index)
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		return g.block.NewCall(g.runtimeArrayGet, objLive, g.emitAsFujiI64(key)), nil
	}

	nilTag := constant.NewInt(types.I64, llvmNilTagged)
	objI := g.block.NewLoad(types.I64, objSlot)
	isNil := g.block.NewICmp(enum.IPredEQ, objI, nilTag)

	g.tempN++
	suf := fmt.Sprintf(".opt%d", g.tempN)
	skip := g.currentFn.NewBlock("optnil" + suf)
	cont := g.currentFn.NewBlock("optget" + suf)
	merge := g.currentFn.NewBlock("optmg" + suf)
	g.block.NewCondBr(isNil, skip, cont)

	g.block = skip
	skip.NewBr(merge)

	g.block = cont
	key, err := g.emitExpr(e.Index)
	if err != nil {
		return nil, err
	}
	got := g.block.NewCall(g.runtimeArrayGet, objI, g.emitAsFujiI64(key))
	cont.NewBr(merge)

	g.block = merge
	return merge.NewPhi(
		ir.NewIncoming(nilTag, skip),
		ir.NewIncoming(got, cont),
	), nil
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
		if slot, ok := g.moduleGlobals[name]; ok {
			boxed := g.emitAsFujiI64(val)
			if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return val, nil
		}
		if global, ok := g.globals[name]; ok {
			g.block.NewStore(g.emitAsFujiI64(val), global)
			return val, nil
		}
		return nil, g.undefinedVarError(name, left.Name.Line, left.Name.Col)
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
		if slot, ok := g.moduleGlobals[name]; ok {
			boxed := g.emitAsFujiI64(result)
			if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return result, nil
		}
		if global, ok := g.globals[name]; ok {
			g.block.NewStore(g.emitAsFujiI64(result), global)
			return result, nil
		}
		return nil, g.undefinedVarError(name, left.Name.Line, left.Name.Col)
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

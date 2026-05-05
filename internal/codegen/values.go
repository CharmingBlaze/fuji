package codegen

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"fuji/internal/parser"
)

// emitLiteral emits LLVM IR for literal expressions.
func (g *Generator) emitLiteral(e *parser.LiteralExpr) (value.Value, error) {
	switch v := e.Value.(type) {
	case nil:
		// parser.LiteralExpr for `null` — must match NIL_VAL (see shadow.llvmNilTagged, runtime/value.h).
		return constant.NewInt(types.I64, llvmNilTagged), nil
	case float64:
		return g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, v)), nil
	case int:
		// Integer literals must be NaN-boxed numbers (same representation as runtime NUMBER_VAL).
		return g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, float64(v))), nil
	case bool:
		// Must match runtime TRUE_VAL / FALSE_VAL (see emitBoxBoolNaN, value.h).
		if v {
			return constant.NewInt(types.I64, 0x7ffc000000000003), nil
		}
		return constant.NewInt(types.I64, 0x7ffc000000000002), nil
	case string:
		return g.emitStringLiteral(v), nil
	default:
		return nil, fmt.Errorf("unsupported literal type: %T", v)
	}
}

// emitIdentifier emits LLVM IR for identifier references.
func (g *Generator) emitIdentifier(e *parser.IdentifierExpr) (value.Value, error) {
	name := e.Name.Lexeme

	if slot, ok := g.locals[name]; ok {
		if g.localIsCell != nil && g.localIsCell[name] {
			return g.block.NewCall(g.runtimeCellRead, slot), nil
		}
		return g.block.NewLoad(types.I64, slot), nil
	}

	// Check if it's a global variable
	if global, ok := g.globals[name]; ok {
		return g.block.NewLoad(types.I64, global), nil
	}

	// Check if it's a function
	if fn, ok := g.funcs[name]; ok {
		return fn, nil
	}

	return nil, fmt.Errorf("undefined variable: %s", name)
}

// emitStringLiteral creates a string object using the runtime.
func (g *Generator) emitStringLiteral(s string) value.Value {
	// Create a string literal as a global constant
	arr := constant.NewCharArrayFromString(s)
	global := g.mod.NewGlobalDef("", arr)
	global.Immutable = true
	global.Linkage = enum.LinkagePrivate

	// Get pointer to the string data
	zero := constant.NewInt(types.I32, 0)
	ptr := g.block.NewGetElementPtr(arr.Type(), global, zero, zero)

	// Call runtime to create string object
	length := constant.NewInt(types.I32, int64(len(s)))
	strObj := g.block.NewCall(g.runtimeAllocStr, length, ptr)

	return strObj
}

// emitInfix emits LLVM IR for infix expressions.
func (g *Generator) emitInfix(e *parser.InfixExpr) (value.Value, error) {
	switch e.Operator {
	case "&&":
		return g.emitLogicalAnd(e)
	case "||":
		return g.emitLogicalOr(e)
	}

	left, err := g.emitExpr(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := g.emitExpr(e.Right)
	if err != nil {
		return nil, err
	}

	leftI := g.emitAsFujiI64(left)
	rightI := g.emitAsFujiI64(right)

	switch e.Operator {
	case "+":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		return g.block.NewCall(g.runtimeBoxNumber, g.block.NewFAdd(ld, rd)), nil
	case "-":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		return g.block.NewCall(g.runtimeBoxNumber, g.block.NewFSub(ld, rd)), nil
	case "*":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		return g.block.NewCall(g.runtimeBoxNumber, g.block.NewFMul(ld, rd)), nil
	case "/":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		return g.block.NewCall(g.runtimeBoxNumber, g.block.NewFDiv(ld, rd)), nil
	case "%":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		return g.block.NewCall(g.runtimeBoxNumber, g.block.NewFRem(ld, rd)), nil
	case "==", "===":
		eq := g.block.NewCall(g.runtimeValuesEqual, leftI, rightI)
		cmp := g.block.NewICmp(enum.IPredEQ, eq, constant.NewInt(types.I64, 1))
		return g.emitBoxBoolNaN(cmp), nil
	case "!=", "!==":
		eq := g.block.NewCall(g.runtimeValuesEqual, leftI, rightI)
		cmp := g.block.NewICmp(enum.IPredEQ, eq, constant.NewInt(types.I64, 0))
		return g.emitBoxBoolNaN(cmp), nil
	case "<":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		cmp := g.block.NewFCmp(enum.FPredOLT, ld, rd)
		return g.emitBoxBoolNaN(cmp), nil
	case "<=":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		cmp := g.block.NewFCmp(enum.FPredOLE, ld, rd)
		return g.emitBoxBoolNaN(cmp), nil
	case ">":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		cmp := g.block.NewFCmp(enum.FPredOGT, ld, rd)
		return g.emitBoxBoolNaN(cmp), nil
	case ">=":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		cmp := g.block.NewFCmp(enum.FPredOGE, ld, rd)
		return g.emitBoxBoolNaN(cmp), nil
	case "&", "|", "^":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		l32 := g.block.NewFPToSI(ld, types.I32)
		r32 := g.block.NewFPToSI(rd, types.I32)
		var res32 value.Value
		switch e.Operator {
		case "&":
			res32 = g.block.NewAnd(l32, r32)
		case "|":
			res32 = g.block.NewOr(l32, r32)
		case "^":
			res32 = g.block.NewXor(l32, r32)
		}
		rf := g.block.NewSIToFP(res32, types.Double)
		return g.block.NewCall(g.runtimeBoxNumber, rf), nil
	case "<<", ">>":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		l32 := g.block.NewFPToSI(ld, types.I32)
		r32 := g.block.NewFPToSI(rd, types.I32)
		shiftMask := constant.NewInt(types.I32, 31)
		rM := g.block.NewAnd(r32, shiftMask)
		var res32 value.Value
		if e.Operator == "<<" {
			res32 = g.block.NewShl(l32, rM)
		} else {
			res32 = g.block.NewAShr(l32, rM)
		}
		rf := g.block.NewSIToFP(res32, types.Double)
		return g.block.NewCall(g.runtimeBoxNumber, rf), nil
	case ">>>":
		ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
		rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
		l32 := g.block.NewFPToSI(ld, types.I32)
		r32 := g.block.NewFPToSI(rd, types.I32)
		rM := g.block.NewAnd(r32, constant.NewInt(types.I32, 31))
		ures := g.block.NewLShr(l32, rM)
		rf := g.block.NewUIToFP(ures, types.Double)
		return g.block.NewCall(g.runtimeBoxNumber, rf), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %s", e.Operator)
	}
}

func (g *Generator) emitLogicalAnd(e *parser.InfixExpr) (value.Value, error) {
	resultSlot := g.block.NewAlloca(types.I64)
	g.shadowStoreTemp(resultSlot)
	leftVal, err := g.emitExpr(e.Left)
	if err != nil {
		return nil, err
	}
	g.block.NewStore(g.emitAsFujiI64(leftVal), resultSlot)

	g.tempN++
	rhsBlock := g.currentFn.NewBlock(fmt.Sprintf("land.rhs.%d", g.tempN))
	mergeBlock := g.currentFn.NewBlock(fmt.Sprintf("land.merge.%d", g.tempN))
	g.block.NewCondBr(g.emitTruthy(leftVal), rhsBlock, mergeBlock)

	g.block = rhsBlock
	rightVal, err := g.emitExpr(e.Right)
	if err != nil {
		return nil, err
	}
	g.block.NewStore(g.emitAsFujiI64(rightVal), resultSlot)
	g.block.NewBr(mergeBlock)

	g.block = mergeBlock
	return g.block.NewLoad(types.I64, resultSlot), nil
}

func (g *Generator) emitLogicalOr(e *parser.InfixExpr) (value.Value, error) {
	resultSlot := g.block.NewAlloca(types.I64)
	g.shadowStoreTemp(resultSlot)
	leftVal, err := g.emitExpr(e.Left)
	if err != nil {
		return nil, err
	}
	g.block.NewStore(g.emitAsFujiI64(leftVal), resultSlot)

	g.tempN++
	rhsBlock := g.currentFn.NewBlock(fmt.Sprintf("lor.rhs.%d", g.tempN))
	mergeBlock := g.currentFn.NewBlock(fmt.Sprintf("lor.merge.%d", g.tempN))
	g.block.NewCondBr(g.emitTruthy(leftVal), mergeBlock, rhsBlock)

	g.block = rhsBlock
	rightVal, err := g.emitExpr(e.Right)
	if err != nil {
		return nil, err
	}
	g.block.NewStore(g.emitAsFujiI64(rightVal), resultSlot)
	g.block.NewBr(mergeBlock)

	g.block = mergeBlock
	return g.block.NewLoad(types.I64, resultSlot), nil
}

// emitArray emits LLVM IR for array expressions.
func (g *Generator) emitArray(e *parser.ArrayExpr) (value.Value, error) {
	n := len(e.Elements)
	if n < 1 {
		n = 1 // allow empty [] literal; runtime needs at least one slot before push fills
	}
	capacity := constant.NewInt(types.I32, int64(n))
	arr := g.block.NewCall(g.runtimeAllocArray, capacity)

	// Initialize array elements
	for i, elem := range e.Elements {
		elemVal, err := g.emitExpr(elem)
		if err != nil {
			return nil, err
		}
		index := constant.NewInt(types.I64, int64(i))
		g.block.NewCall(g.runtimeArraySet, arr, index, g.emitAsFujiI64(elemVal))
	}

	return arr, nil
}

// emitThis emits LLVM IR for the this keyword.
func (g *Generator) emitThis(_ *parser.ThisExpr) (value.Value, error) {
	// Return the current this value from the local 'this' variable
	if slot, ok := g.locals["this"]; ok {
		return g.block.NewLoad(types.I64, slot), nil
	}
	// Fallback to nil if 'this' is not found (should not happen in well-formed code)
	return constant.NewInt(types.I64, 0), nil
}

// emitSlice emits LLVM IR for slice expressions.
func (g *Generator) emitSlice(e *parser.SliceExpr) (value.Value, error) {
	obj, err := g.emitExpr(e.Object)
	if err != nil {
		return nil, err
	}
	objI := g.emitAsFujiI64(obj)

	var startV value.Value
	if e.Start != nil {
		startV, err = g.emitExpr(e.Start)
		if err != nil {
			return nil, err
		}
	} else {
		startV = constant.NewFloat(types.Double, 0)
	}
	startI := g.emitAsFujiI64(startV)

	var endV value.Value
	if e.End != nil {
		endV, err = g.emitExpr(e.End)
		if err != nil {
			return nil, err
		}
	} else {
		endV = g.block.NewCall(g.runtimeArrayLen, objI)
	}
	endI := g.emitAsFujiI64(endV)

	zero := constant.NewInt(types.I32, 0)
	arrTy := types.NewArray(3, types.I64)
	slot := g.block.NewAlloca(arrTy)
	g.block.NewStore(objI, g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, 0)))
	g.block.NewStore(startI, g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, 1)))
	g.block.NewStore(endI, g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, 2)))
	argvPtr := g.block.NewGetElementPtr(arrTy, slot, zero, zero)
	return g.block.NewCall(g.runtimeArraySlice, constant.NewInt(types.I32, 3), argvPtr), nil
}

// emitTemplate emits LLVM IR for template expressions.
// Static text and holes are lowered to fuji_format(fmt, ...values) using "{}" placeholders.
func (g *Generator) emitTemplate(e *parser.TemplateExpr) (value.Value, error) {
	if len(e.Parts) == 0 {
		return g.emitStringLiteral(""), nil
	}
	var fmtBuilder strings.Builder
	var holeExprs []parser.Expr
	for _, part := range e.Parts {
		if lit, ok := part.(*parser.LiteralExpr); ok {
			if s, ok := lit.Value.(string); ok {
				fmtBuilder.WriteString(s)
				continue
			}
		}
		fmtBuilder.WriteString("{}")
		holeExprs = append(holeExprs, part)
	}

	fmtVal := g.emitStringLiteral(fmtBuilder.String())
	argv := []value.Value{fmtVal}
	for _, ex := range holeExprs {
		v, err := g.emitExpr(ex)
		if err != nil {
			return nil, err
		}
		argv = append(argv, v)
	}
	return g.emitArgvRuntime(g.runtimeFormat, argv), nil
}

// emitRange emits LLVM IR for range expressions.
// Semantics match fuji_range: half-open integer range [from, to) of numbers.
func (g *Generator) emitRange(e *parser.RangeExpr) (value.Value, error) {
	fromVal, err := g.emitExpr(e.From)
	if err != nil {
		return nil, err
	}

	toVal, err := g.emitExpr(e.To)
	if err != nil {
		return nil, err
	}

	zero := constant.NewInt(types.I32, 0)
	arrTy := types.NewArray(2, types.I64)
	slot := g.block.NewAlloca(arrTy)
	g.block.NewStore(g.emitAsFujiI64(fromVal), g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, 0)))
	g.block.NewStore(g.emitAsFujiI64(toVal), g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, 1)))
	argvPtr := g.block.NewGetElementPtr(arrTy, slot, zero, zero)
	return g.block.NewCall(g.runtimeRange, constant.NewInt(types.I32, 2), argvPtr), nil
}

// emitTuple emits LLVM IR for tuple expressions.
func (g *Generator) emitTuple(e *parser.TupleExpr) (value.Value, error) {
	if len(e.Elements) == 0 {
		return nil, fmt.Errorf("native codegen: empty tuple is not supported")
	}
	return nil, fmt.Errorf("native codegen: tuple expressions are not supported yet")
}

// emitImport emits LLVM IR for import expressions.
func (g *Generator) emitImport(_ *parser.ImportExpr) (value.Value, error) {
	return nil, fmt.Errorf("native codegen: dynamic import() is not supported; use static @path imports at file scope")
}

// emitSwitchExpr emits LLVM IR for switch expressions (expression-level switch).
func (g *Generator) emitSwitchExpr(e *parser.SwitchExpr) (value.Value, error) {
	resultSlot := g.block.NewAlloca(types.I64)
	g.shadowStoreTemp(resultSlot)

	subj, err := g.emitExpr(e.Subject)
	if err != nil {
		return nil, err
	}
	subjSlot := g.block.NewAlloca(types.I64)
	g.shadowStoreTemp(subjSlot)
	g.block.NewStore(g.emitAsFujiI64(subj), subjSlot)

	mergeB := g.block.Parent.NewBlock("swexpr.merge")

	for i, c := range e.Cases {
		missB := g.block.Parent.NewBlock(fmt.Sprintf("swexpr.miss.%d", i))
		bodyB := g.block.Parent.NewBlock(fmt.Sprintf("swexpr.body.%d", i))

		subjL := g.block.NewLoad(types.I64, subjSlot)
		cv, err := g.emitExpr(c.Value)
		if err != nil {
			return nil, err
		}
		cmp := g.block.NewICmp(enum.IPredEQ, subjL, g.emitAsFujiI64(cv))
		g.block.NewCondBr(cmp, bodyB, missB)

		g.block = bodyB
		bodyVal, err := g.emitExpr(c.Body)
		if err != nil {
			return nil, err
		}
		g.block.NewStore(g.emitAsFujiI64(bodyVal), resultSlot)
		g.block.NewBr(mergeB)

		g.block = missB
	}

	defVal, err := g.emitExpr(e.Default)
	if err != nil {
		return nil, err
	}
	g.block.NewStore(g.emitAsFujiI64(defVal), resultSlot)
	g.block.NewBr(mergeB)

	g.block = mergeB
	return g.block.NewLoad(types.I64, resultSlot), nil
}

// emitSwitchCaseExpr emits LLVM IR for switch case expressions.
func (g *Generator) emitSwitchCaseExpr(e *parser.SwitchCaseExpr) (value.Value, error) {
	return g.emitExpr(e.Body)
}

// emitUpdate emits LLVM IR for increment/decrement operators (++, --).
func (g *Generator) emitUpdate(e *parser.UpdateExpr) (value.Value, error) {
	current, err := g.emitExpr(e.Operand)
	if err != nil {
		return nil, err
	}
	curI := g.emitAsFujiI64(current)
	ld := g.block.NewCall(g.runtimeUnboxNumber, curI)
	var rd value.Value
	switch e.Operator.Lexeme {
	case "++":
		rd = g.block.NewFAdd(ld, constant.NewFloat(types.Double, 1))
	case "--":
		rd = g.block.NewFSub(ld, constant.NewFloat(types.Double, 1))
	default:
		return nil, fmt.Errorf("unsupported update operator: %s", e.Operator.Lexeme)
	}
	result := g.block.NewCall(g.runtimeBoxNumber, rd)

	if ident, ok := e.Operand.(*parser.IdentifierExpr); ok {
		name := ident.Name.Lexeme
		if slot, ok := g.locals[name]; ok {
			boxed := g.emitAsFujiI64(result)
			if g.localIsCell != nil && g.localIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
		}
	}

	if e.IsPrefix {
		return result, nil
	}
	return current, nil
}

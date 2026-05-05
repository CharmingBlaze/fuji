package codegen

import (
	"fmt"
	"strings"

	"fuji/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// tryEmitMethodCall lowers argv runtime calls and LLVM-backed `.map()` / `.filter()` / `.find()` / `.reduce()`.
func (g *Generator) tryEmitMethodCall(member *parser.IndexExpr, recvVal value.Value, call *parser.CallExpr) (value.Value, bool, error) {
	lit, ok := member.Index.(*parser.LiteralExpr)
	if !ok {
		return nil, false, nil
	}
	name, ok := lit.Value.(string)
	if !ok {
		return nil, false, nil
	}
	name = strings.ToLower(name)

	switch name {
	case "concat":
		args := []value.Value{recvVal}
		for _, arg := range call.Arguments {
			v, err := g.emitExpr(arg)
			if err != nil {
				return nil, true, err
			}
			args = append(args, v)
		}
		return g.emitArgvRuntime(g.runtimeArrayConcat, args), true, nil

	case "split", "trim", "toupper", "tolower", "replace", "replaceall", "startswith", "endswith":
		v, err := g.emitStringMethod(name, recvVal, call.Arguments)
		return v, true, err

	case "join", "sort", "reverse", "push", "pop":
		v, err := g.emitArrayOnlyMethod(name, recvVal, call.Arguments)
		return v, true, err

	case "slice":
		v, err := g.emitSliceAmbiguous(recvVal, call.Arguments)
		return v, true, err
	case "indexof", "includes":
		v, err := g.emitIndexOrIncludesAmbiguous(name, recvVal, call.Arguments)
		return v, true, err

	case "map":
		return g.emitArrayMethodMap(recvVal, call.Arguments)
	case "filter":
		return g.emitArrayMethodFilter(recvVal, call.Arguments)
	case "find":
		return g.emitArrayMethodFind(recvVal, call.Arguments)
	case "reduce":
		return g.emitArrayMethodReduce(recvVal, call.Arguments)
	}

	return nil, false, nil
}

func (g *Generator) emitStringMethod(name string, recv value.Value, args []parser.Expr) (value.Value, error) {
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }

	switch name {
	case "split":
		if len(args) != 1 {
			return nil, fmt.Errorf("split expects 1 argument (delimiter)")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringSplit, []value.Value{loadRecv(), a0}), nil
	case "trim":
		return g.emitArgvRuntime(g.runtimeStringTrim, []value.Value{loadRecv()}), nil
	case "toupper":
		return g.emitArgvRuntime(g.runtimeStringUpper, []value.Value{loadRecv()}), nil
	case "tolower":
		return g.emitArgvRuntime(g.runtimeStringLower, []value.Value{loadRecv()}), nil
	case "replace":
		if len(args) != 2 {
			return nil, fmt.Errorf("replace expects 2 arguments")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := g.emitExpr(args[1])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringReplace, []value.Value{loadRecv(), a0, a1}), nil
	case "replaceall":
		if len(args) != 2 {
			return nil, fmt.Errorf("replaceAll expects 2 arguments")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := g.emitExpr(args[1])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringReplaceAll, []value.Value{loadRecv(), a0, a1}), nil
	case "startswith":
		if len(args) != 1 {
			return nil, fmt.Errorf("startsWith expects 1 argument")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringStartsWith, []value.Value{loadRecv(), a0}), nil
	case "endswith":
		if len(args) != 1 {
			return nil, fmt.Errorf("endsWith expects 1 argument")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringEndsWith, []value.Value{loadRecv(), a0}), nil
	default:
		return nil, fmt.Errorf("unknown string method %q", name)
	}
}

func (g *Generator) emitArrayOnlyMethod(name string, recv value.Value, args []parser.Expr) (value.Value, error) {
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }

	switch name {
	case "join":
		if len(args) != 1 {
			return nil, fmt.Errorf("join expects 1 argument (separator)")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeArrayJoin, []value.Value{loadRecv(), a0}), nil
	case "sort":
		return g.emitArgvRuntime(g.runtimeArraySort, []value.Value{loadRecv()}), nil
	case "reverse":
		return g.emitArgvRuntime(g.runtimeArrayReverse, []value.Value{loadRecv()}), nil
	case "push":
		if len(args) != 1 {
			return nil, fmt.Errorf("push expects 1 argument")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		g.block.NewCall(g.runtimeArrayPush, loadRecv(), g.emitAsFujiI64(a0))
		return loadRecv(), nil
	case "pop":
		if len(args) != 0 {
			return nil, fmt.Errorf("pop expects 0 arguments")
		}
		return g.block.NewCall(g.runtimeArrayPop, loadRecv()), nil
	default:
		return nil, fmt.Errorf("unknown array method %q", name)
	}
}

func (g *Generator) emitSliceAmbiguous(recv value.Value, args []parser.Expr) (value.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice expects 2 arguments (start, end)")
	}
	a0, err := g.emitExpr(args[0])
	if err != nil {
		return nil, err
	}
	a1, err := g.emitExpr(args[1])
	if err != nil {
		return nil, err
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }
	return g.emitArgvNilTaggedFallback(g.emitArgvRuntime(g.runtimeArraySlice, []value.Value{loadRecv(), a0, a1}), func() value.Value {
		return g.emitArgvRuntime(g.runtimeStringSlice, []value.Value{loadRecv(), a0, a1})
	}), nil
}

func (g *Generator) emitIndexOrIncludesAmbiguous(name string, recv value.Value, args []parser.Expr) (value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects 1 argument", name)
	}
	a0, err := g.emitExpr(args[0])
	if err != nil {
		return nil, err
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }
	argv := []value.Value{loadRecv(), a0}
	switch name {
	case "indexof":
		va := g.emitArgvRuntime(g.runtimeArrayIndexOf, argv)
		return g.emitArgvNilTaggedFallback(va, func() value.Value {
			return g.emitArgvRuntime(g.runtimeStringIndexOf, argv)
		}), nil
	case "includes":
		va := g.emitArgvRuntime(g.runtimeArrayIncludes, argv)
		return g.emitArgvNilTaggedFallback(va, func() value.Value {
			// Strings use contains-like helpers via argv wrappers only for booleans — approximate via indexOf >= 0.
			io := g.emitArgvRuntime(g.runtimeStringIndexOf, argv)
			negOne := constant.NewFloat(types.Double, -1)
			boxNeg := g.block.NewCall(g.runtimeBoxNumber, negOne)
			cmp := g.block.NewFCmp(enum.FPredOGT, g.block.NewCall(g.runtimeUnboxNumber, g.emitAsFujiI64(io)), g.block.NewCall(g.runtimeUnboxNumber, boxNeg))
			return g.emitBoxBoolNaN(cmp)
		}), nil
	default:
		return nil, fmt.Errorf("unknown ambiguous method %q", name)
	}
}

// emitArgvNilTaggedFallback returns primary if it is non-nil-tagged; otherwise evaluates alt() on an alternate block.
func (g *Generator) emitArgvNilTaggedFallback(primary value.Value, alt func() value.Value) value.Value {
	entry := g.block
	nilTag := constant.NewInt(types.I64, llvmNilTagged)
	isNil := entry.NewICmp(enum.IPredEQ, g.emitAsFujiI64(primary), nilTag)

	g.tempN++
	suf := fmt.Sprintf(".nf%d", g.tempN)
	okB := g.currentFn.NewBlock("nilfb.ok" + suf)
	altB := g.currentFn.NewBlock("nilfb.alt" + suf)
	mergeB := g.currentFn.NewBlock("nilfb.mg" + suf)
	entry.NewCondBr(isNil, altB, okB)

	g.block = okB
	okB.NewBr(mergeB)

	g.block = altB
	vAlt := alt()
	altTerm := g.block
	altTerm.NewBr(mergeB)

	g.block = mergeB
	return mergeB.NewPhi(
		ir.NewIncoming(primary, okB),
		ir.NewIncoming(vAlt, altTerm),
	)
}

func (g *Generator) emitArrayLenAsDouble(recv value.Value) value.Value {
	lv := g.block.NewCall(g.runtimeLen, g.emitAsFujiI64(recv))
	return g.block.NewCall(g.runtimeUnboxNumber, g.emitAsFujiI64(lv))
}

func (g *Generator) emitArrayMethodMap(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, true, fmt.Errorf("map expects 1 callback")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	nThis := constant.NewInt(types.I64, 0)

	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)
	capI := g.block.NewFPToSI(lenF, types.I32)
	out := g.block.NewCall(g.runtimeAllocArray, capI)
	outSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(outSlot)
	g.block.NewStore(g.emitAsFujiI64(out), outSlot)

	g.tempN++
	suf := fmt.Sprintf(".map%d", g.tempN)
	hdr := g.currentFn.NewBlock("map.hdr" + suf)
	body := g.currentFn.NewBlock("map.body" + suf)
	step := g.currentFn.NewBlock("map.step" + suf)
	done := g.currentFn.NewBlock("map.done" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(constant.NewFloat(types.Double, 0), idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, done)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsFujiI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsFujiI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	mv, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{elLive})
	if err != nil {
		return nil, true, err
	}
	idxInt := g.block.NewFPToSI(idxNow, types.I64)
	outLive := g.block.NewLoad(types.I64, outSlot)
	g.block.NewCall(g.runtimeArraySet, outLive, idxInt, g.emitAsFujiI64(mv))
	g.block.NewBr(step)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = done
	return g.block.NewLoad(types.I64, outSlot), true, nil
}

func (g *Generator) emitArrayMethodFilter(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, true, fmt.Errorf("filter expects 1 callback")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	out := g.block.NewCall(g.runtimeAllocArray, constant.NewInt(types.I32, 1))
	outSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(outSlot)
	g.block.NewStore(g.emitAsFujiI64(out), outSlot)
	nThis := constant.NewInt(types.I64, 0)

	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)

	g.tempN++
	suf := fmt.Sprintf(".filt%d", g.tempN)
	hdr := g.currentFn.NewBlock("filt.hdr" + suf)
	body := g.currentFn.NewBlock("filt.body" + suf)
	pushB := g.currentFn.NewBlock("filt.push" + suf)
	skipB := g.currentFn.NewBlock("filt.skip" + suf)
	step := g.currentFn.NewBlock("filt.step" + suf)
	done := g.currentFn.NewBlock("filt.done" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(constant.NewFloat(types.Double, 0), idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, done)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsFujiI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsFujiI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	predV, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{elLive})
	if err != nil {
		return nil, true, err
	}
	brOk := g.emitTruthy(predV)
	g.block.NewCondBr(brOk, pushB, skipB)

	g.block = pushB
	outLive := g.block.NewLoad(types.I64, outSlot)
	pushVal := g.block.NewLoad(types.I64, elSlot)
	g.block.NewCall(g.runtimeArrayPush, outLive, pushVal)
	g.block.NewBr(step)

	g.block = skipB
	skipB.NewBr(step)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = done
	return g.block.NewLoad(types.I64, outSlot), true, nil
}

func (g *Generator) emitArrayMethodFind(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, true, fmt.Errorf("find expects 1 callback")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	nThis := constant.NewInt(types.I64, 0)
	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)

	g.tempN++
	suf := fmt.Sprintf(".find%d", g.tempN)
	merge := g.currentFn.NewBlock("find.merge" + suf)
	hdr := g.currentFn.NewBlock("find.hdr" + suf)
	body := g.currentFn.NewBlock("find.body" + suf)
	foundHit := g.currentFn.NewBlock("find.hit" + suf)
	step := g.currentFn.NewBlock("find.step" + suf)
	miss := g.currentFn.NewBlock("find.miss" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(constant.NewFloat(types.Double, 0), idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, miss)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsFujiI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsFujiI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	predV, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{elLive})
	if err != nil {
		return nil, true, err
	}
	brOk := g.emitTruthy(predV)
	g.block.NewCondBr(brOk, foundHit, step)

	g.block = foundHit
	foundHit.NewBr(merge)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = miss
	nilV := constant.NewInt(types.I64, llvmNilTagged)
	miss.NewBr(merge)

	g.block = merge
	return merge.NewPhi(
		ir.NewIncoming(elLive, foundHit),
		ir.NewIncoming(nilV, miss),
	), true, nil
}

func (g *Generator) emitArrayMethodReduce(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 && len(args) != 2 {
		return nil, true, fmt.Errorf("reduce expects 1 or 2 arguments")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsFujiI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	nThis := constant.NewInt(types.I64, 0)

	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)
	zero := constant.NewFloat(types.Double, 0)
	isEmpty := g.block.NewFCmp(enum.FPredOEQ, lenF, zero)

	g.tempN++
	suf := fmt.Sprintf(".red%d", g.tempN)
	emptyB := g.currentFn.NewBlock("red.empty" + suf)
	startB := g.currentFn.NewBlock("red.setup" + suf)
	mergeFinal := g.currentFn.NewBlock("red.merge" + suf)

	g.block.NewCondBr(isEmpty, emptyB, startB)

	var emptyOut value.Value
	g.block = emptyB
	if len(args) == 1 {
		emptyOut = constant.NewInt(types.I64, llvmNilTagged)
	} else {
		initV, err := g.emitExpr(args[1])
		if err != nil {
			return nil, true, err
		}
		emptyOut = g.emitAsFujiI64(initV)
	}
	emptyB.NewBr(mergeFinal)

	g.block = startB
	var accSlot value.Value
	var startIdx value.Value
	if len(args) == 1 {
		zBox := g.block.NewCall(g.runtimeBoxNumber, zero)
		firstEl := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsFujiI64(zBox))
		accSlot = g.entryAlloca(types.I64)
		g.shadowStoreTemp(accSlot)
		g.block.NewStore(g.emitAsFujiI64(firstEl), accSlot)
		startIdx = constant.NewFloat(types.Double, 1)
	} else {
		initV, err := g.emitExpr(args[1])
		if err != nil {
			return nil, true, err
		}
		accSlot = g.entryAlloca(types.I64)
		g.shadowStoreTemp(accSlot)
		g.block.NewStore(g.emitAsFujiI64(initV), accSlot)
		startIdx = zero
	}

	hdr := g.currentFn.NewBlock("red.hdr" + suf)
	body := g.currentFn.NewBlock("red.body" + suf)
	step := g.currentFn.NewBlock("red.step" + suf)
	done := g.currentFn.NewBlock("red.done" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(startIdx, idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, done)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsFujiI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsFujiI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	accLoad := g.block.NewLoad(types.I64, accSlot)
	rv, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{accLoad, elLive})
	if err != nil {
		return nil, true, err
	}
	g.block.NewStore(g.emitAsFujiI64(rv), accSlot)
	g.block.NewBr(step)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = done
	loopOut := g.block.NewLoad(types.I64, accSlot)
	done.NewBr(mergeFinal)

	g.block = mergeFinal
	return mergeFinal.NewPhi(
		ir.NewIncoming(emptyOut, emptyB),
		ir.NewIncoming(loopOut, done),
	), true, nil
}

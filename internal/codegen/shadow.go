package codegen

import (
	"fuji/internal/parser"
	"fuji/internal/sema"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// LLVM i64 representation of NIL_VAL (see runtime/src/value.h).
const llvmNilTagged int64 = 0x7ffc000000000001

func (g *Generator) emitShadowPop() {
	if !g.shadowPushed {
		return
	}
	g.block.NewCall(g.runtimePopFrame)
	g.shadowPushed = false
	g.shadowFramePtr = nil
	g.shadowFrameArrTy = nil
	g.shadowTempNext = 0
}

func (g *Generator) beginShadowFrame(layout *sema.ShadowLayout, thisSlot value.Value, orderedParamNames []string) {
	ptrPtrTy := types.NewPointer(types.NewPointer(types.I64))
	if layout == nil || layout.Total <= 0 {
		g.shadowFramePtr = nil
		g.shadowFrameArrTy = nil
		g.shadowTempNext = 0
		g.block.NewCall(g.runtimePushFrame, constant.NewNull(ptrPtrTy), constant.NewInt(types.I32, 0))
		g.shadowPushed = true
		return
	}

	n := layout.Total
	elemTy := types.NewPointer(types.I64)
	arrTy := types.NewArray(uint64(n), elemTy)
	frame := g.block.NewAlloca(arrTy)
	nilSlot := g.block.NewAlloca(types.I64)
	g.block.NewStore(constant.NewInt(types.I64, llvmNilTagged), nilSlot)

	zero := constant.NewInt(types.I32, 0)
	for i := 0; i < n; i++ {
		ep := g.block.NewGetElementPtr(arrTy, frame, zero, constant.NewInt(types.I32, int64(i)))
		g.block.NewStore(nilSlot, ep)
	}

	g.shadowFrameArrTy = arrTy
	g.shadowFramePtr = frame

	thisEp := g.block.NewGetElementPtr(arrTy, frame, zero, constant.NewInt(types.I32, 0))
	g.block.NewStore(thisSlot, thisEp)
	for i, name := range orderedParamNames {
		pSlot, ok := g.locals[name]
		if !ok {
			continue
		}
		pe := g.block.NewGetElementPtr(arrTy, frame, zero, constant.NewInt(types.I32, int64(1+i)))
		g.block.NewStore(pSlot, pe)
	}

	row0 := g.block.NewGetElementPtr(arrTy, frame, zero, zero)
	g.block.NewCall(g.runtimePushFrame, row0, constant.NewInt(types.I32, int64(n)))
	g.shadowPushed = true
	g.shadowTempNext = layout.TempBase
}

func (g *Generator) shadowStoreIndex(idx int, slotPtr value.Value) {
	if g.shadowFramePtr == nil || g.shadowFrameArrTy == nil {
		return
	}
	zero := constant.NewInt(types.I32, 0)
	ep := g.block.NewGetElementPtr(g.shadowFrameArrTy, g.shadowFramePtr, zero, constant.NewInt(types.I32, int64(idx)))
	g.block.NewStore(slotPtr, ep)
}

func (g *Generator) shadowStoreLet(d *parser.LetDecl, slotPtr value.Value) {
	if g.shadowLayout == nil {
		return
	}
	idx, ok := g.shadowLayout.LetIndex[d]
	if !ok {
		return
	}
	g.shadowStoreIndex(idx, slotPtr)
}

func (g *Generator) shadowStoreForOf(s *parser.ForOfStmt, slotPtr value.Value) {
	if g.shadowLayout == nil {
		return
	}
	idx, ok := g.shadowLayout.ForOfIndex[s]
	if !ok {
		return
	}
	g.shadowStoreIndex(idx, slotPtr)
}

// shadowStoreTemp registers an ad-hoc i64* alloca that may hold a GC-managed
// Value across a call or allocation (short-circuit &&/||, switch-expr).
func (g *Generator) shadowStoreTemp(slotPtr value.Value) {
	if g.shadowLayout == nil || g.shadowFramePtr == nil {
		return
	}
	if g.shadowTempNext >= g.shadowLayout.Total {
		return
	}
	idx := g.shadowTempNext
	g.shadowTempNext++
	g.shadowStoreIndex(idx, slotPtr)
}

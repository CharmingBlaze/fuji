package codegen

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// entryAlloca places stack allocations in the function entry block so they are
// allocated once per function call (not once per loop iteration).
func (g *Generator) entryAlloca(allocType types.Type) value.Value {
	if g.currentFn == nil || len(g.currentFn.Blocks) == 0 {
		return g.block.NewAlloca(allocType)
	}
	entry := g.currentFn.Blocks[0]
	return entry.NewAlloca(allocType)
}

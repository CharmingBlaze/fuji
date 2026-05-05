package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"fuji/internal/parser"
)

// emitStmt emits LLVM IR for statements.
func (g *Generator) emitStmt(stmt parser.Stmt) error {
	switch s := stmt.(type) {
	case *parser.BlockStmt:
		return g.emitBlockStmt(s)
	case *parser.ExpressionStmt:
		_, err := g.emitExpr(s.Expr)
		return err
	case *parser.ReturnStmt:
		return g.emitReturnStmt(s)
	case *parser.IfStmt:
		return g.emitIfStmt(s)
	case *parser.WhileStmt:
		return g.emitWhileStmt(s)
	case *parser.DoWhileStmt:
		return g.emitDoWhileStmt(s)
	case *parser.ForOfStmt:
		return g.emitForOfStmt(s)
	case *parser.ForStmt:
		return g.emitForStmt(s)
	case *parser.ForInStmt:
		return g.emitForInStmt(s)
	case *parser.BreakStmt:
		return g.emitBreakStmt(s)
	case *parser.ContinueStmt:
		return g.emitContinueStmt(s)
	case *parser.SwitchStmt:
		return g.emitSwitchStmt(s)
	case *parser.DeleteStmt:
		return g.emitDeleteStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// emitBlockStmt emits LLVM IR for block statements.
func (g *Generator) emitBlockStmt(s *parser.BlockStmt) error {
	for _, decl := range s.Declarations {
		if err := g.emitDecl(decl); err != nil {
			return err
		}
	}
	return nil
}

// emitReturnStmt emits LLVM IR for return statements.
func (g *Generator) emitReturnStmt(s *parser.ReturnStmt) error {
	if s.Value != nil {
		val, err := g.emitExpr(s.Value)
		if err != nil {
			return err
		}
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(g.emitAsFujiI64(val))
	} else {
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(constant.NewInt(types.I64, 0))
	}
	return nil
}

// emitIfStmt emits LLVM IR for if statements.
func (g *Generator) emitIfStmt(s *parser.IfStmt) error {
	cond, err := g.emitExpr(s.Condition)
	if err != nil {
		return err
	}

	g.tempN++
	thenBlock := g.block.Parent.NewBlock(fmt.Sprintf("then.%d", g.tempN))
	elseBlock := g.block.Parent.NewBlock(fmt.Sprintf("else.%d", g.tempN))
	mergeBlock := g.block.Parent.NewBlock(fmt.Sprintf("merge.%d", g.tempN))

	g.block.NewCondBr(g.emitTruthy(cond), thenBlock, elseBlock)

	g.block = thenBlock
	if err := g.emitStmt(s.Then); err != nil {
		return err
	}
	if g.block.Term == nil {
		g.block.NewBr(mergeBlock)
	}

	g.block = elseBlock
	if s.Else != nil {
		if err := g.emitStmt(s.Else); err != nil {
			return err
		}
	}
	if g.block.Term == nil {
		g.block.NewBr(mergeBlock)
	}

	g.block = mergeBlock
	return nil
}

// emitWhileStmt emits LLVM IR for while loops.
func (g *Generator) emitWhileStmt(s *parser.WhileStmt) error {
	condBlock := g.block.Parent.NewBlock("while.cond")
	bodyBlock := g.block.Parent.NewBlock("while.body")
	afterBlock := g.block.Parent.NewBlock("while.after")

	// Push loop context onto stack
	ctx := loopContext{condBlock: condBlock, incBlock: condBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	cond, err := g.emitExpr(s.Condition)
	if err != nil {
		return err
	}
	g.block.NewCondBr(g.emitTruthy(cond), bodyBlock, afterBlock)

	g.block = bodyBlock
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(condBlock)

	g.block = afterBlock

	// Pop loop context from stack
	g.loopStack = g.loopStack[:len(g.loopStack)-1]

	return nil
}

// emitDoWhileStmt emits LLVM IR for do-while loops.
func (g *Generator) emitDoWhileStmt(s *parser.DoWhileStmt) error {
	bodyBlock := g.block.Parent.NewBlock("dowhile.body")
	condBlock := g.block.Parent.NewBlock("dowhile.cond")
	afterBlock := g.block.Parent.NewBlock("dowhile.after")

	// Push loop context onto stack
	ctx := loopContext{condBlock: condBlock, incBlock: condBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(bodyBlock)

	g.block = bodyBlock
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(condBlock)

	g.block = condBlock
	cond, err := g.emitExpr(s.Condition)
	if err != nil {
		return err
	}
	g.block.NewCondBr(g.emitTruthy(cond), bodyBlock, afterBlock)

	g.block = afterBlock

	// Pop loop context from stack
	g.loopStack = g.loopStack[:len(g.loopStack)-1]

	return nil
}

// emitForOfStmt emits LLVM IR for for-of loops over arrays and tables (slot order).
func (g *Generator) emitForOfStmt(s *parser.ForOfStmt) error {
	iterable, err := g.emitExpr(s.Iterable)
	if err != nil {
		return err
	}
	iterI := g.emitAsFujiI64(iterable)

	lenVal := g.block.NewCall(g.runtimeForOfLength, iterI)

	idxSlot := g.entryAlloca(types.I64)
	zeroNum := g.emitAsFujiI64(constant.NewFloat(types.Double, 0))
	g.block.NewStore(zeroNum, idxSlot)

	valSlot := g.entryAlloca(types.I64)
	var keySlot value.Value
	if s.ValueVar != nil {
		keySlot = g.entryAlloca(types.I64)
		g.locals[s.VarName.Lexeme] = keySlot
		g.locals[s.ValueVar.Lexeme] = valSlot
	} else {
		g.locals[s.VarName.Lexeme] = valSlot
	}

	g.tempN++
	condBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.cond.%d", g.tempN))
	bodyBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.body.%d", g.tempN))
	incBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.inc.%d", g.tempN))
	afterBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.after.%d", g.tempN))

	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	idxV := g.block.NewLoad(types.I64, idxSlot)
	idxD := g.block.NewBitCast(idxV, types.Double)
	lenD := g.block.NewBitCast(lenVal, types.Double)
	cmp := g.block.NewFCmp(enum.FPredOLT, idxD, lenD)
	g.block.NewCondBr(cmp, bodyBlock, afterBlock)

	g.block = bodyBlock
	if s.ValueVar != nil {
		key := g.block.NewCall(g.runtimeForOfKeyAt, iterI, idxV)
		val := g.block.NewCall(g.runtimeForOfValueAt, iterI, idxV)
		g.block.NewStore(key, keySlot)
		g.block.NewStore(val, valSlot)
	} else {
		val := g.block.NewCall(g.runtimeForOfValueAt, iterI, idxV)
		g.block.NewStore(val, valSlot)
	}
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	g.block = incBlock
	idxV2 := g.block.NewLoad(types.I64, idxSlot)
	idxD2 := g.block.NewBitCast(idxV2, types.Double)
	nextD := g.block.NewFAdd(idxD2, constant.NewFloat(types.Double, 1))
	nextI := g.block.NewBitCast(nextD, types.I64)
	g.block.NewStore(nextI, idxSlot)
	g.block.NewBr(condBlock)

	g.block = afterBlock

	g.loopStack = g.loopStack[:len(g.loopStack)-1]
	return nil
}

// emitForStmt emits LLVM IR for C-style for loops.
func (g *Generator) emitForStmt(s *parser.ForStmt) error {
	// Create blocks
	condBlock := g.block.Parent.NewBlock("for.cond")
	bodyBlock := g.block.Parent.NewBlock("for.body")
	incBlock := g.block.Parent.NewBlock("for.inc")
	afterBlock := g.block.Parent.NewBlock("for.after")

	// Push loop context onto stack
	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	// Emit initialization
	for _, init := range s.Inits {
		if err := g.emitDecl(init); err != nil {
			return err
		}
	}

	// Jump to condition block
	g.block.NewBr(condBlock)

	// Condition block
	g.block = condBlock
	if s.Condition != nil {
		cond, err := g.emitExpr(s.Condition)
		if err != nil {
			return err
		}
		g.block.NewCondBr(g.emitTruthy(cond), bodyBlock, afterBlock)
	} else {
		// No condition means always enter loop
		g.block.NewBr(bodyBlock)
	}

	// Body block
	g.block = bodyBlock
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	// Increment block
	g.block = incBlock
	for _, inc := range s.Increments {
		_, err := g.emitExpr(inc)
		if err != nil {
			return err
		}
	}
	g.block.NewBr(condBlock)

	// After block
	g.block = afterBlock

	// Pop loop context from stack
	g.loopStack = g.loopStack[:len(g.loopStack)-1]

	return nil
}

// emitForInStmt emits LLVM IR for for-in loops (keys in insertion order for tables).
func (g *Generator) emitForInStmt(s *parser.ForInStmt) error {
	if s.ValueVar != nil {
		return fmt.Errorf("native codegen: for-in with two variables is not supported")
	}
	if s.KeyVar == nil {
		return fmt.Errorf("native codegen: for-in missing loop variable")
	}
	iterable, err := g.emitExpr(s.Iterable)
	if err != nil {
		return err
	}
	iterI := g.emitAsFujiI64(iterable)

	lenVal := g.block.NewCall(g.runtimeForOfLength, iterI)

	idxSlot := g.entryAlloca(types.I64)
	zeroNum := g.emitAsFujiI64(constant.NewFloat(types.Double, 0))
	g.block.NewStore(zeroNum, idxSlot)

	keySlot := g.entryAlloca(types.I64)
	g.locals[s.KeyVar.Lexeme] = keySlot

	g.tempN++
	condBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.cond.%d", g.tempN))
	bodyBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.body.%d", g.tempN))
	incBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.inc.%d", g.tempN))
	afterBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.after.%d", g.tempN))

	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	idxV := g.block.NewLoad(types.I64, idxSlot)
	idxD := g.block.NewBitCast(idxV, types.Double)
	lenD := g.block.NewBitCast(lenVal, types.Double)
	cmp := g.block.NewFCmp(enum.FPredOLT, idxD, lenD)
	g.block.NewCondBr(cmp, bodyBlock, afterBlock)

	g.block = bodyBlock
	key := g.block.NewCall(g.runtimeForOfKeyAt, iterI, idxV)
	g.block.NewStore(key, keySlot)
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	g.block = incBlock
	idxV2 := g.block.NewLoad(types.I64, idxSlot)
	idxD2 := g.block.NewBitCast(idxV2, types.Double)
	nextD := g.block.NewFAdd(idxD2, constant.NewFloat(types.Double, 1))
	nextI := g.block.NewBitCast(nextD, types.I64)
	g.block.NewStore(nextI, idxSlot)
	g.block.NewBr(condBlock)

	g.block = afterBlock

	g.loopStack = g.loopStack[:len(g.loopStack)-1]
	return nil
}

// emitBreakStmt emits LLVM IR for break statements.
func (g *Generator) emitBreakStmt(_ *parser.BreakStmt) error {
	// Break statement - jump to the loop's after block
	if len(g.loopStack) == 0 {
		return fmt.Errorf("break statement outside of loop")
	}
	ctx := g.loopStack[len(g.loopStack)-1]
	g.block.NewBr(ctx.afterBlock)
	return nil
}

// emitContinueStmt emits LLVM IR for continue statements.
func (g *Generator) emitContinueStmt(_ *parser.ContinueStmt) error {
	// Continue statement - jump to the loop's increment or condition block
	if len(g.loopStack) == 0 {
		return fmt.Errorf("continue statement outside of loop")
	}
	ctx := g.loopStack[len(g.loopStack)-1]
	g.block.NewBr(ctx.incBlock)
	return nil
}

// emitSwitchStmt emits LLVM IR for switch statements.
func (g *Generator) emitSwitchStmt(s *parser.SwitchStmt) error {
	// Lower switch to compare/branch chain (NaN-boxed Fuji values compared as i64).
	subject, err := g.emitExpr(s.Subject)
	if err != nil {
		return err
	}
	subjI := g.emitAsFujiI64(subject)

	// Create merge block
	mergeBlock := g.block.Parent.NewBlock("switch.merge")

	swCtx := loopContext{condBlock: mergeBlock, incBlock: mergeBlock, afterBlock: mergeBlock}
	g.loopStack = append(g.loopStack, swCtx)
	defer func() { g.loopStack = g.loopStack[:len(g.loopStack)-1] }()

	// For each case, create a conditional branch
	caseBlocks := make([]*ir.Block, len(s.Cases)+1)
	caseBlocks[len(s.Cases)] = g.block.Parent.NewBlock("switch.default")

	for i := range s.Cases {
		caseBlocks[i] = g.block.Parent.NewBlock(fmt.Sprintf("switch.case.%d", i))
	}

	// Emit comparisons and branches
	for i, caseStmt := range s.Cases {
		caseVal, err := g.emitExpr(caseStmt.Value)
		if err != nil {
			return err
		}

		cmp := g.block.NewICmp(enum.IPredEQ, subjI, g.emitAsFujiI64(caseVal))
		nextBlock := caseBlocks[i+1]
		if i == len(s.Cases)-1 {
			nextBlock = caseBlocks[len(s.Cases)]
		}
		g.block.NewCondBr(cmp, caseBlocks[i], nextBlock)
		g.block = nextBlock
	}

	// Emit default case
	g.block = caseBlocks[len(s.Cases)]
	for _, decl := range s.Default {
		if err := g.emitDecl(decl); err != nil {
			return err
		}
	}
	g.block.NewBr(mergeBlock)

	// Emit case blocks
	for i, caseStmt := range s.Cases {
		g.block = caseBlocks[i]

		// Emit case body declarations
		for _, decl := range caseStmt.Body {
			if err := g.emitDecl(decl); err != nil {
				return err
			}
		}

		// Fall-through to next case or merge (classic C behavior)
		if i < len(s.Cases)-1 {
			g.block.NewBr(caseBlocks[i+1])
		} else {
			g.block.NewBr(mergeBlock)
		}
	}

	g.block = mergeBlock
	return nil
}

func (g *Generator) emitDeleteStmt(s *parser.DeleteStmt) error {
	ix, ok := s.Target.(*parser.IndexExpr)
	if !ok {
		return fmt.Errorf("delete expects a property access expression such as obj[\"key\"]")
	}
	obj, err := g.emitExpr(ix.Object)
	if err != nil {
		return err
	}
	key, err := g.emitExpr(ix.Index)
	if err != nil {
		return err
	}
	g.block.NewCall(g.runtimeObjRemove, g.emitAsFujiI64(obj), g.emitAsFujiI64(key))
	return nil
}

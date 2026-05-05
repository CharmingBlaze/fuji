package sema

import (
	"fuji/internal/parser"
	"reflect"
)

// ParamCellKey identifies a function parameter for escape-driven heap-cell lowering.
type ParamCellKey struct {
	Func uintptr
	Idx  int
}

// NewParamCellKey builds a stable key for a parameter on a *parser.FuncDecl or *parser.FuncExpr.
func NewParamCellKey(owner interface{}, idx int) ParamCellKey {
	return ParamCellKey{Func: reflect.ValueOf(owner).Pointer(), Idx: idx}
}

type NativeEmitContext struct {
	Bundle       *parser.ProgramBundle
	locals       map[parser.Expr]parser.Decl
	capturedVars map[parser.Decl]bool
	paramDecls   map[*parser.FuncDecl][]parser.Decl
	currentFn    *parser.FuncDecl

	FuncCaptures map[*parser.FuncDecl][]*parser.LetDecl
	ExprCaptures map[*parser.FuncExpr][]*parser.LetDecl

	EscapingDecls map[*parser.LetDecl]bool
	StackDecls    map[*parser.LetDecl]bool
	ParamIsCell   map[ParamCellKey]bool
	letOwner      map[*parser.LetDecl]interface{}

	// FreeVarsExpr / FreeVarsDecl: outer bindings referenced by nested functions (for codegen closure seeding).
	FreeVarsExpr map[*parser.FuncExpr][]string
	FreeVarsDecl map[*parser.FuncDecl][]string

	// FuncExprEnclosing maps each function expression to its enclosing function (FuncExpr or FuncDecl)
	// while walking, so captures can be propagated to intermediate closures (nested lambdas).
	FuncExprEnclosing map[*parser.FuncExpr]interface{}

	ShadowFuncDecl map[*parser.FuncDecl]*ShadowLayout
	ShadowFuncExpr map[*parser.FuncExpr]*ShadowLayout
	ShadowEntry    *ShadowLayout
}

func PrepareNativeBundle(bundle *parser.ProgramBundle) (*NativeEmitContext, error) {
	if err := parser.FlattenEntryIncludes(bundle); err != nil {
		return nil, err
	}
	parser.InjectNativeMathPrelude(bundle)
	ctx := &NativeEmitContext{
		Bundle:         bundle,
		locals:         make(map[parser.Expr]parser.Decl),
		capturedVars:   make(map[parser.Decl]bool),
		paramDecls:     make(map[*parser.FuncDecl][]parser.Decl),
		FuncCaptures:   make(map[*parser.FuncDecl][]*parser.LetDecl),
		ExprCaptures:   make(map[*parser.FuncExpr][]*parser.LetDecl),
		EscapingDecls:  make(map[*parser.LetDecl]bool),
		StackDecls:     make(map[*parser.LetDecl]bool),
		ParamIsCell:    make(map[ParamCellKey]bool),
		letOwner:       make(map[*parser.LetDecl]interface{}),
		FreeVarsExpr:      make(map[*parser.FuncExpr][]string),
		FreeVarsDecl:      make(map[*parser.FuncDecl][]string),
		FuncExprEnclosing: make(map[*parser.FuncExpr]interface{}),
		ShadowFuncDecl:    make(map[*parser.FuncDecl]*ShadowLayout),
		ShadowFuncExpr: make(map[*parser.FuncExpr]*ShadowLayout),
	}
	prepareNativeAnalysis(ctx, bundle)
	prepareShadowLayouts(ctx, bundle)
	return ctx, nil
}

func ValidateNativeEmitSupport(ctx *NativeEmitContext) error {
	_ = ctx
	return nil
}

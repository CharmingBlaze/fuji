package sema

import (
	"fmt"
	"fuji/internal/parser"
)

// Analyzer performs semantic analysis on the AST.
type Analyzer struct {
	currentScope *Scope
	scopes       []*Scope
	errors       []error
}

// Scope represents a lexical scope with symbol bindings.
type Scope struct {
	parent  *Scope
	symbols map[string]parser.Decl
}

// NewScope creates a new scope.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]parser.Decl),
	}
}

// Define adds a symbol to the scope.
func (s *Scope) Define(name string, decl parser.Decl) {
	s.symbols[name] = decl
}

// Resolve looks up a symbol in the scope chain.
func (s *Scope) Resolve(name string) (parser.Decl, bool) {
	if decl, ok := s.symbols[name]; ok {
		return decl, true
	}
	if s.parent != nil {
		return s.parent.Resolve(name)
	}
	return nil, false
}

// NewAnalyzer creates a new semantic analyzer.
func NewAnalyzer() *Analyzer {
	globalScope := NewScope(nil)
	return &Analyzer{
		currentScope: globalScope,
		scopes:       []*Scope{globalScope},
		errors:       []error{},
	}
}

// Analyze performs semantic analysis on a program.
func (a *Analyzer) Analyze(prog *parser.Program) error {
	for _, decl := range prog.Declarations {
		if err := a.analyzeDecl(decl); err != nil {
			a.errors = append(a.errors, err)
		}
	}
	if len(a.errors) > 0 {
		return a.errors[0]
	}
	return nil
}

// Errors returns all errors found during analysis.
func (a *Analyzer) Errors() []error {
	return a.errors
}

func (a *Analyzer) enterScope() {
	newScope := NewScope(a.currentScope)
	a.scopes = append(a.scopes, newScope)
	a.currentScope = newScope
}

func (a *Analyzer) exitScope() {
	if len(a.scopes) > 1 {
		a.scopes = a.scopes[:len(a.scopes)-1]
		a.currentScope = a.scopes[len(a.scopes)-1]
	}
}

func (a *Analyzer) analyzeDecl(decl parser.Decl) error {
	switch d := decl.(type) {
	case *parser.LetDecl:
		return a.analyzeLetDecl(d)
	case *parser.FuncDecl:
		return a.analyzeFuncDecl(d)
	case *parser.IncludeDecl:
		return nil
	case parser.Stmt:
		return a.analyzeStmt(d)
	}
	return nil
}

func (a *Analyzer) analyzeLetDecl(d *parser.LetDecl) error {
	name := d.Name.Lexeme
	a.currentScope.Define(name, d)
	if d.Init != nil {
		return a.analyzeExpr(d.Init)
	}
	return nil
}

func (a *Analyzer) analyzeFuncDecl(d *parser.FuncDecl) error {
	name := d.Name.Lexeme
	a.currentScope.Define(name, d)

	a.enterScope()

	for _, param := range d.Params {
		a.currentScope.Define(param.Name, d)
		if param.Default != nil {
			if err := a.analyzeExpr(param.Default); err != nil {
				return err
			}
		}
	}

	if err := a.analyzeStmt(d.Body); err != nil {
		a.exitScope()
		return err
	}

	a.exitScope()
	return nil
}

func (a *Analyzer) analyzeStmt(stmt parser.Stmt) error {
	switch s := stmt.(type) {
	case *parser.BlockStmt:
		return a.analyzeBlockStmt(s)
	case *parser.ExpressionStmt:
		return a.analyzeExpr(s.Expr)
	case *parser.ReturnStmt:
		if s.Value != nil {
			return a.analyzeExpr(s.Value)
		}
		return nil
	case *parser.IfStmt:
		if err := a.analyzeExpr(s.Condition); err != nil {
			return err
		}
		if err := a.analyzeStmt(s.Then); err != nil {
			return err
		}
		if s.Else != nil {
			return a.analyzeStmt(s.Else)
		}
		return nil
	case *parser.WhileStmt:
		if err := a.analyzeExpr(s.Condition); err != nil {
			return err
		}
		return a.analyzeStmt(s.Body)
	case *parser.BreakStmt, *parser.ContinueStmt:
		return nil
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

func (a *Analyzer) analyzeBlockStmt(s *parser.BlockStmt) error {
	a.enterScope()
	defer a.exitScope()

	for _, decl := range s.Declarations {
		if err := a.analyzeDecl(decl); err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) analyzeExpr(expr parser.Expr) error {
	switch e := expr.(type) {
	case *parser.IdentifierExpr:
		name := e.Name.Lexeme
		if _, ok := a.currentScope.Resolve(name); !ok {
			return fmt.Errorf("undefined variable '%s'", name)
		}
		return nil
	case *parser.LiteralExpr:
		return nil
	case *parser.PrefixExpr:
		return a.analyzeExpr(e.Right)
	case *parser.InfixExpr:
		if err := a.analyzeExpr(e.Left); err != nil {
			return err
		}
		return a.analyzeExpr(e.Right)
	case *parser.LogicalExpr:
		if err := a.analyzeExpr(e.Left); err != nil {
			return err
		}
		return a.analyzeExpr(e.Right)
	case *parser.CallExpr:
		if err := a.analyzeExpr(e.Function); err != nil {
			return err
		}
		for _, arg := range e.Arguments {
			if err := a.analyzeExpr(arg); err != nil {
				return err
			}
		}
		return nil
	case *parser.AssignExpr:
		if err := a.analyzeExpr(e.Value); err != nil {
			return err
		}
		if ident, ok := e.Left.(*parser.IdentifierExpr); ok {
			name := ident.Name.Lexeme
			if _, ok := a.currentScope.Resolve(name); !ok {
				return fmt.Errorf("undefined variable '%s'", name)
			}
			return nil
		}
		return fmt.Errorf("invalid assignment target")
	case *parser.GroupingExpr:
		return a.analyzeExpr(e.Expr)
	case *parser.ImportExpr:
		return nil
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

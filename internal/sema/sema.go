package sema

import (
	"fmt"

	"fuji/internal/diagnostic"
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

func (s *Scope) VisibleNames() []string {
	out := []string{}
	seen := map[string]bool{}
	for cur := s; cur != nil; cur = cur.parent {
		for name := range cur.symbols {
			if seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

// NewAnalyzer creates a new semantic analyzer.
func NewAnalyzer() *Analyzer {
	builtinRoot := NewScope(nil)
	seedGlobalBuiltins(builtinRoot)
	globalScope := NewScope(builtinRoot)
	return &Analyzer{
		currentScope: globalScope,
		scopes:       []*Scope{builtinRoot, globalScope},
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
	case *parser.FuncExpr:
		return a.analyzeFuncExpr(d)
	case *parser.IncludeDecl:
		return nil
	case parser.Stmt:
		return a.analyzeStmt(d)
	}
	return nil
}

func (a *Analyzer) analyzeLetDecl(d *parser.LetDecl) error {
	name := d.Name.Lexeme
	if _, ok := a.currentScope.symbols[name]; ok {
		return &diagnostic.DiagnosticError{
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("duplicate binding '%s' in the same scope", name),
		}
	}
	a.currentScope.Define(name, d)
	if d.Init != nil {
		return a.analyzeExpr(d.Init)
	}
	return nil
}

func (a *Analyzer) analyzeFuncDecl(d *parser.FuncDecl) error {
	name := d.Name.Lexeme
	if _, ok := a.currentScope.symbols[name]; ok {
		return &diagnostic.DiagnosticError{
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("duplicate function '%s' in the same scope", name),
		}
	}
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

func (a *Analyzer) analyzeFuncExpr(e *parser.FuncExpr) error {
	a.enterScope()
	defer a.exitScope()
	for _, param := range e.Params {
		a.currentScope.Define(param.Name, e)
		if param.Default != nil {
			if err := a.analyzeExpr(param.Default); err != nil {
				return err
			}
		}
	}
	return a.analyzeStmt(e.Body)
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
	case *parser.DoWhileStmt:
		if err := a.analyzeStmt(s.Body); err != nil {
			return err
		}
		return a.analyzeExpr(s.Condition)
	case *parser.ForStmt:
		for _, ini := range s.Inits {
			if err := a.analyzeDecl(ini); err != nil {
				return err
			}
		}
		if s.Condition != nil {
			if err := a.analyzeExpr(s.Condition); err != nil {
				return err
			}
		}
		for _, inc := range s.Increments {
			if err := a.analyzeExpr(inc); err != nil {
				return err
			}
		}
		return a.analyzeStmt(s.Body)
	case *parser.ForInStmt:
		if err := a.analyzeExpr(s.Iterable); err != nil {
			return err
		}
		return a.analyzeStmt(s.Body)
	case *parser.ForOfStmt:
		if err := a.analyzeExpr(s.Iterable); err != nil {
			return err
		}
		return a.analyzeStmt(s.Body)
	case *parser.SwitchStmt:
		if err := a.analyzeExpr(s.Subject); err != nil {
			return err
		}
		for _, c := range s.Cases {
			if err := a.analyzeExpr(c.Value); err != nil {
				return err
			}
			for _, cd := range c.Body {
				if err := a.analyzeDecl(cd); err != nil {
					return err
				}
			}
		}
		for _, cd := range s.Default {
			if err := a.analyzeDecl(cd); err != nil {
				return err
			}
		}
		return nil
	case *parser.DeleteStmt:
		return a.analyzeExpr(s.Target)
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
			hint := ""
			if s, ok := diagnostic.BestSuggestion(name, a.currentScope.VisibleNames(), 2); ok {
				hint = fmt.Sprintf("did you mean '%s'?", s)
			}
			return &diagnostic.DiagnosticError{
				Line:    e.Name.Line,
				Col:     e.Name.Col,
				Message: fmt.Sprintf("undefined variable '%s'", name),
				Hint:    hint,
			}
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
				hint := ""
				if s, ok := diagnostic.BestSuggestion(name, a.currentScope.VisibleNames(), 2); ok {
					hint = fmt.Sprintf("did you mean '%s'?", s)
				}
				return &diagnostic.DiagnosticError{
					Line:    ident.Name.Line,
					Col:     ident.Name.Col,
					Message: fmt.Sprintf("undefined variable '%s'", name),
					Hint:    hint,
				}
			}
			return nil
		}
		if ix, ok := e.Left.(*parser.IndexExpr); ok {
			if err := a.analyzeExpr(ix.Object); err != nil {
				return err
			}
			return a.analyzeExpr(ix.Index)
		}
		return &diagnostic.DiagnosticError{
			Line:    e.Token.Line,
			Col:     e.Token.Col,
			Message: "invalid assignment target",
			Hint:    "left side of '=' must be a variable or index expression",
		}
	case *parser.GroupingExpr:
		return a.analyzeExpr(e.Expr)
	case *parser.ImportExpr:
		return nil
	case *parser.IndexExpr:
		if err := a.analyzeExpr(e.Object); err != nil {
			return err
		}
		return a.analyzeExpr(e.Index)
	case *parser.SpreadExpr:
		return a.analyzeExpr(e.Expr)
	case *parser.TemplateExpr:
		for _, p := range e.Parts {
			if err := a.analyzeExpr(p); err != nil {
				return err
			}
		}
		return nil
	case *parser.ThisExpr:
		return nil
	case *parser.ArrayExpr:
		for _, el := range e.Elements {
			if err := a.analyzeExpr(el); err != nil {
				return err
			}
		}
		return nil
	case *parser.ObjectExpr:
		for _, v := range e.Values {
			if err := a.analyzeExpr(v); err != nil {
				return err
			}
		}
		for _, ck := range e.ComputedKeys {
			if ck == nil {
				continue
			}
			if err := a.analyzeExpr(ck); err != nil {
				return err
			}
		}
		return nil
	case *parser.FuncExpr:
		return a.analyzeFuncExpr(e)
	case *parser.RangeExpr:
		if err := a.analyzeExpr(e.From); err != nil {
			return err
		}
		return a.analyzeExpr(e.To)
	case *parser.UpdateExpr:
		return a.analyzeExpr(e.Operand)
	case *parser.TupleExpr:
		for _, el := range e.Elements {
			if err := a.analyzeExpr(el); err != nil {
				return err
			}
		}
		return nil
	case *parser.IfExpr:
		if err := a.analyzeExpr(e.Condition); err != nil {
			return err
		}
		if err := a.analyzeExpr(e.Then); err != nil {
			return err
		}
		if e.Else != nil {
			return a.analyzeExpr(e.Else)
		}
		return nil
	case *parser.SwitchExpr:
		if err := a.analyzeExpr(e.Subject); err != nil {
			return err
		}
		for _, c := range e.Cases {
			if err := a.analyzeExpr(c.Value); err != nil {
				return err
			}
			if err := a.analyzeExpr(c.Body); err != nil {
				return err
			}
		}
		if e.Default != nil {
			return a.analyzeExpr(e.Default)
		}
		return nil
	case *parser.SliceExpr:
		if err := a.analyzeExpr(e.Object); err != nil {
			return err
		}
		if e.Start != nil {
			if err := a.analyzeExpr(e.Start); err != nil {
				return err
			}
		}
		if e.End != nil {
			return a.analyzeExpr(e.End)
		}
		return nil
	case *parser.TernaryExpr:
		if err := a.analyzeExpr(e.Condition); err != nil {
			return err
		}
		if err := a.analyzeExpr(e.Then); err != nil {
			return err
		}
		return a.analyzeExpr(e.Else)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

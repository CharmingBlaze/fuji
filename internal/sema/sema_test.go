package sema

import (
	"fuji/internal/lexer"
	"fuji/internal/parser"
	"testing"
)

func parseForTest(t *testing.T, source string) *parser.Program {
	t.Helper()
	l := lexer.NewLexer(source)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}
	p := parser.NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}
	return program
}

func TestSemaDuplicateBindingSameScope(t *testing.T) {
	src := `let x = 1;
let X = 2;
`
	program := parseForTest(t, src)
	err := NewAnalyzer().Analyze(program)
	if err == nil {
		t.Fatal("expected duplicate binding error")
	}
}

func TestSemaBasicScoping(t *testing.T) {
	source := `
		let x = 10;
		func add(a, b) {
			let y = a + b;
			return y;
		}
		add(x, 2);
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
	if len(analyzer.Errors()) != 0 {
		t.Fatalf("Unexpected errors: %v", analyzer.Errors())
	}
}

func TestSemaUndefinedVariable(t *testing.T) {
	source := `
		let x = 10;
		y + 5;
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err == nil {
		t.Fatal("Expected error for undefined variable")
	}
	if len(analyzer.Errors()) == 0 {
		t.Fatal("Expected errors to be collected")
	}
}

func TestSemaBlockScoping(t *testing.T) {
	source := `
		let x = 10;
		if (x > 5) {
			let y = x;
			y = y + 1;
		}
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
}

func TestSemaFunctionParams(t *testing.T) {
	source := `
		func f(a, b = 2, ...rest) {
			return a;
		}
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
}

func TestSemaAssignmentTarget(t *testing.T) {
	source := `
		let x = 10;
		x = 20;
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
}

package codegen

import (
	"strings"
	"testing"

	"fuji/internal/lexer"
	"fuji/internal/parser"
	"fuji/internal/sema"
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

func TestCodegenBasic(t *testing.T) {
	source := `
		func add(a, b) {
			return a + b;
		}
		let result = add(1, 2);
	`
	program := parseForTest(t, source)

	analyzer := sema.NewAnalyzer()
	if err := analyzer.Analyze(program); err != nil {
		t.Fatalf("Sema failed: %v", err)
	}

	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("PrepareNativeBundle failed: %v", err)
	}
	gen := NewGenerator(ctx)

	mod, err := gen.Generate(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("Codegen failed: %v", err)
	}
	if mod == nil {
		t.Fatal("Expected non-nil module")
	}
}

func TestCodegenControlFlow(t *testing.T) {
	source := `
		func test(x) {
			if (x > 0) {
				return x;
			}
			return 0;
		}
	`
	program := parseForTest(t, source)

	analyzer := sema.NewAnalyzer()
	if err := analyzer.Analyze(program); err != nil {
		t.Fatalf("Sema failed: %v", err)
	}

	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("PrepareNativeBundle failed: %v", err)
	}
	gen := NewGenerator(ctx)

	mod, err := gen.Generate(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("Codegen failed: %v", err)
	}
	if mod == nil {
		t.Fatal("Expected non-nil module")
	}
}

func TestNativeExternDeclaresGlueSymbol(t *testing.T) {
	src := `
// fuji:extern Foo FUJI_shim_Foo 1
let Foo = 0;
func main() {
	Foo(1);
}
`
	program := parseForTest(t, src)
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("PrepareNativeBundle: %v", err)
	}
	gen := NewGenerator(ctx)
	mod, err := gen.Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "@FUJI_shim_Foo") {
		t.Fatalf("expected LLVM to declare @FUJI_shim_Foo, got:\n%s", ir)
	}
}

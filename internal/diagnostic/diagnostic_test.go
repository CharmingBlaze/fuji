package diagnostic

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExtractPositionParser(t *testing.T) {
	err := fmt.Errorf("[line 2:5] error at 'x': expected thing")
	pos, ok := ExtractPosition(err)
	if !ok || pos.Line != 2 || pos.Col != 5 {
		t.Fatalf("got (%v, %v) ok=%v", pos.Line, pos.Col, ok)
	}
}

func TestExtractPositionLexer(t *testing.T) {
	err := errors.New("unexpected character at 3:7: !")
	pos, ok := ExtractPosition(err)
	if !ok || pos.Line != 3 || pos.Col != 7 {
		t.Fatalf("got (%v, %v) ok=%v", pos.Line, pos.Col, ok)
	}
}

func TestSnippetCaret(t *testing.T) {
	src := "one\ntwo x\nthree"
	sn := Snippet(src, Position{Line: 2, Col: 5})
	if sn == "" {
		t.Fatal("empty snippet")
	}
	if !strings.Contains(sn, "2 | two x") || !strings.Contains(sn, "^") {
		t.Fatalf("snippet:\n%s", sn)
	}
}

func TestSourceContextError(t *testing.T) {
	src := "let x = ;\n"
	cause := fmt.Errorf("[line 1:9] error at ';': boom")
	e := WrapParse("/tmp/a.fuji", src, cause)
	s := e.Error()
	if !strings.Contains(s, "/tmp/a.fuji") || !strings.Contains(s, "^") {
		t.Fatal(s)
	}
}

func TestIdentifierSpanAt(t *testing.T) {
	if got := IdentifierSpanAt(`  player.x`, 3); got != len("player") {
		t.Fatalf("span=%d want %d", got, len("player"))
	}
	if IdentifierSpanAt(`  x`, 99) != 1 {
		t.Fatal("oob col")
	}
}

func TestFormatErrorDiagnosticSnippet(t *testing.T) {
	path := t.TempDir() + "/t.fuji"
	src := "let a = 1;\nlet b = mistyped + 2;\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	err := &DiagnosticError{
		File:    path,
		Line:    2,
		Col:     9,
		Message: "undefined variable 'mistyped'",
		Hint:    "did you mean 'a'?",
	}
	out := FormatError(err)
	if !strings.Contains(out, "-->") || !strings.Contains(out, "mistyped") || !strings.Contains(out, "^") {
		t.Fatalf("FormatError:\n%s", out)
	}
}

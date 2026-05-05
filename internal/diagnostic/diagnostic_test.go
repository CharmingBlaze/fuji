package diagnostic

import (
	"errors"
	"fmt"
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

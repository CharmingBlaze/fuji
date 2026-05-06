package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiagnoseValidProgram(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ok.fuji")
	src := "let x = 1;\nprint(x);\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) != 0 {
		t.Fatalf("want no diagnostics, got %#v", d)
	}
}

func TestDiagnoseSemaUndefined(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.fuji")
	if err := os.WriteFile(p, []byte("notARealBinding;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want diagnostics")
	}
	msg := strings.ToLower(d[0].Message)
	if !strings.Contains(msg, "undefined") {
		t.Fatalf("message: %q", d[0].Message)
	}
}

func TestDiagnoseOverlayReplacesOnDiskSource(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "entry.fuji")
	if err := os.WriteFile(p, []byte("brokenOnDisk;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	overlay := "let ok = 1;\nprint(ok);\n"
	d := Diagnose(p, overlay)
	if len(d) != 0 {
		t.Fatalf("overlay should supply valid sema input, got %#v", d)
	}
}

func TestDiagnoseCallArity(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "arity.fuji")
	src := "func f(a) { return a; }\nf(1, 2);\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want arity diagnostic")
	}
	if !strings.Contains(strings.ToLower(d[0].Message), "too many") {
		t.Fatalf("message: %q", d[0].Message)
	}
}

func TestDiagnoseMultipleSemaErrorsAggregated(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "multi.fuji")
	src := "aaa;\nbbb;\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want diagnostics")
	}
	msg := strings.ToLower(d[0].Message)
	if !strings.Contains(msg, "aaa") {
		t.Fatalf("want first undefined in message: %q", d[0].Message)
	}
	if !strings.Contains(msg, "bbb") {
		t.Fatalf("want second undefined in message: %q", d[0].Message)
	}
}

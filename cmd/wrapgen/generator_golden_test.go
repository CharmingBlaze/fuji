package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMinimalHeader(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "minimal.h")
	src := filepath.Join("testdata", "minimal.h")
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	if err := os.WriteFile(header, b, 0o644); err != nil {
		t.Fatalf("write temp header: %v", err)
	}

	out := filepath.Join(dir, "out")
	cfg := &WrapGenConfig{
		LibraryName:   "minimal",
		InputHeaders:  []string{header},
		OutputDir:     out,
		Documentation: false,
		BuildSystem:   false,
		IncludeTests:  false,
	}
	wg := NewWrapperGenerator(cfg)
	api, err := wg.ParseHeaders()
	if err != nil {
		t.Fatalf("ParseHeaders: %v", err)
	}
	if len(api.Functions) < 2 {
		t.Fatalf("expected at least 2 functions, got %d", len(api.Functions))
	}
	var names []string
	for _, f := range api.Functions {
		n := strings.TrimSpace(f.Name)
		if n != "" {
			names = append(names, n)
		}
	}
	joined := strings.Join(names, " ")
	if !strings.Contains(joined, "wg_reset") || !strings.Contains(joined, "wg_add") {
		t.Fatalf("expected wg_reset and wg_add in %#v", names)
	}
}

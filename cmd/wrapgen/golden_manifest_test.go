package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWrapgenGoldenManifestValid(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	p := filepath.Join(root, "testdata", "wrapgen_golden", "manifest.json")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var m struct {
		RaylibPinned struct {
			HeaderURL         string `json:"header_url"`
			HeaderSHA256      string `json:"header_sha256"`
			RaylibFujiSHA256  string `json:"raylib_fuji_sha256"`
			WrapperCSHA256    string `json:"wrapper_c_sha256"`
			MinFunctions      int    `json:"min_functions"`
			MaxFunctions      int    `json:"max_functions"`
		} `json:"raylib_pinned"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json: %v", err)
	}
	r := m.RaylibPinned
	if r.HeaderURL == "" || r.HeaderSHA256 == "" || r.RaylibFujiSHA256 == "" || r.WrapperCSHA256 == "" {
		t.Fatalf("missing golden fields: %+v", r)
	}
	if r.MinFunctions <= 0 || r.MaxFunctions < r.MinFunctions {
		t.Fatalf("invalid function bounds: min=%d max=%d", r.MinFunctions, r.MaxFunctions)
	}
}

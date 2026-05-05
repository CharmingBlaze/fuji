package api

import (
	"path/filepath"
	"regexp"
	"strconv"

	"fuji/internal/parser"
	"fuji/internal/sema"
)

var reLineCol = regexp.MustCompile(`\[line (\d+):(\d+)\]`)
var reLine = regexp.MustCompile(`\[line (\d+)\]`)

func diagnosticsFromError(path string, err error) []Diagnostic {
	if err == nil {
		return nil
	}
	abs, _ := filepath.Abs(path)
	msg := err.Error()
	d := Diagnostic{Path: abs, Severity: "error", Message: msg, Line: 1, Col: 1}
	if m := reLineCol.FindStringSubmatch(msg); len(m) == 3 {
		d.Line, _ = strconv.Atoi(m[1])
		d.Col, _ = strconv.Atoi(m[2])
		return []Diagnostic{d}
	}
	if m := reLine.FindStringSubmatch(msg); len(m) == 2 {
		d.Line, _ = strconv.Atoi(m[1])
		d.Col = 1
		return []Diagnostic{d}
	}
	return []Diagnostic{d}
}

// Diagnose runs load + semantic preparation for path; optional overlay replaces on-disk source for that file only.
func Diagnose(path, overlay string) []Diagnostic {
	abs, err := filepath.Abs(path)
	if err != nil {
		return []Diagnostic{{Path: path, Line: 1, Col: 1, Message: err.Error(), Severity: "error"}}
	}
	overlays := map[string]string{}
	if overlay != "" {
		overlays[abs] = overlay
	}
	bundle, err := parser.LoadProgramWithOverlays(path, overlays)
	if err != nil {
		return diagnosticsFromError(path, err)
	}
	_, err = sema.PrepareNativeBundle(bundle)
	if err != nil {
		return diagnosticsFromError(path, err)
	}
	return []Diagnostic{}
}

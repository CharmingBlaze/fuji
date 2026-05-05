package main

import (
	"os"
	"path/filepath"
)

// toolDisplayName is the binary basename (e.g. fujiwrap.exe) for CLI messages ("Run … -help").
func toolDisplayName() string {
	return filepath.Base(os.Args[0])
}

// generatedByBrand is the canonical name stamped into .fuji / wrapper.c / docs output.
const generatedByBrand = "fujiwrap"

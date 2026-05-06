package nativebuild

import "fmt"

// BuildOptions configures the native link pipeline (llc optimisation level, IR passes, etc.).
type BuildOptions struct {
	// NoOpt forces unoptimised llc (-O0) and skips the OptimiseIR step (for debugging and A/B).
	NoOpt bool
	// Debug enables DWARF debug info emission (for debugger symbolization/stepping).
	Debug bool
	// OptLevel is passed to llc as -O1 .. -O3 when NoOpt is false. Zero selects -O2.
	OptLevel int
}

func (o BuildOptions) llcOptFlag() string {
	if o.NoOpt {
		return "-O0"
	}
	if o.OptLevel >= 1 && o.OptLevel <= 3 {
		return fmt.Sprintf("-O%d", o.OptLevel)
	}
	return "-O2"
}

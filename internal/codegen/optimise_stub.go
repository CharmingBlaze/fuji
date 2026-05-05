//go:build !cgo || !llvm14

package codegen

// OptimiseIR is a no-op when CGo is disabled or when the llvm14 build tag is not set.
// For IR optimisation via LLVM, build with CGo enabled and -tags llvm14, with LLVM
// development libraries on the linker path (see CONTRIBUTING.md).
func OptimiseIR(irText string) (string, error) {
	return irText, nil
}

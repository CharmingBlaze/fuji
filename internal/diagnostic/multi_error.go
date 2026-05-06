package diagnostic

import (
	"errors"
	"fmt"
	"strings"
)

// MultiError groups several independent diagnostics (e.g. multiple sema issues in one compile).
// It implements Unwrap() []error so errors.As can match the first *DiagnosticError in the list.
type MultiError struct {
	// Label is optional context for the summary line (e.g. entry file path).
	Label string
	List  []error
}

func (m *MultiError) Error() string {
	if m == nil || len(m.List) == 0 {
		return "no errors"
	}
	var b strings.Builder
	if m.Label != "" {
		fmt.Fprintf(&b, "%d errors in %s:\n", len(m.List), m.Label)
	} else {
		fmt.Fprintf(&b, "%d errors:\n", len(m.List))
	}
	for i, e := range m.List {
		if i > 0 {
			b.WriteByte('\n')
		}
		line := formatMultiErrorItem(i+1, e)
		b.WriteString(line)
	}
	return b.String()
}

func formatMultiErrorItem(idx int, e error) string {
	var d *DiagnosticError
	if errors.As(e, &d) && d != nil && d.File != "" && d.Line > 0 && d.Col > 0 {
		return fmt.Sprintf("  [%d] %s:%d:%d — %s", idx, d.File, d.Line, d.Col, d.Message)
	}
	return fmt.Sprintf("  [%d] %s", idx, e.Error())
}

// Unwrap implements multi-error unwrap for errors.As / errors.Join semantics (Go 1.20+).
func (m *MultiError) Unwrap() []error {
	if m == nil {
		return nil
	}
	return m.List
}

package api

// Diagnostic is a single editor / LSP-style diagnostic for a .fuji file.
type Diagnostic struct {
	Path     string `json:"path"`
	Line     int    `json:"line"`
	Col      int    `json:"col"`
	Message  string `json:"message"`
	Hint     string `json:"hint,omitempty"`
	Severity string `json:"severity"`
}

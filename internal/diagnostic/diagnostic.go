package diagnostic

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reParserLineCol = regexp.MustCompile(`\[line (\d+):(\d+)\]`)
	reUnexpected    = regexp.MustCompile(`unexpected character at (\d+):(\d+):`)
	reHexLit        = regexp.MustCompile(`invalid hex literal at (\d+):(\d+)`)
	reBinLit        = regexp.MustCompile(`invalid binary literal at (\d+):(\d+)`)
	reExponent      = regexp.MustCompile(`invalid exponent in number at (\d+):(\d+)`)
	reUnterminated  = regexp.MustCompile(`unterminated string at (\d+)`)
)

// Position is a 1-based line and column in source text (matching lexer/parser output).
type Position struct {
	Line int
	Col  int
}

// DiagnosticError is a structured compiler error with optional hint text.
type DiagnosticError struct {
	File    string
	Line    int
	Col     int
	Message string
	Hint    string
}

func (e *DiagnosticError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.File == "" || e.Line <= 0 || e.Col <= 0 {
		return e.Message
	}
	return fmt.Sprintf("%s:%d:%d: %s", e.File, e.Line, e.Col, e.Message)
}

// ExtractPosition parses line/column from lexer or parser error messages.
func ExtractPosition(err error) (pos Position, ok bool) {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if p, ok := extractFromString(e.Error()); ok {
			return p, true
		}
	}
	return Position{}, false
}

func extractFromString(msg string) (Position, bool) {
	if m := reParserLineCol.FindStringSubmatch(msg); len(m) == 3 {
		line, err1 := strconv.Atoi(m[1])
		col, err2 := strconv.Atoi(m[2])
		if err1 == nil && err2 == nil {
			return Position{Line: line, Col: col}, true
		}
	}
	tryPairs := []*regexp.Regexp{reUnexpected, reHexLit, reBinLit, reExponent}
	for _, re := range tryPairs {
		if m := re.FindStringSubmatch(msg); len(m) == 3 {
			line, err1 := strconv.Atoi(m[1])
			col, err2 := strconv.Atoi(m[2])
			if err1 == nil && err2 == nil {
				return Position{Line: line, Col: col}, true
			}
		}
	}
	if m := reUnterminated.FindStringSubmatch(msg); len(m) == 2 {
		if line, err := strconv.Atoi(m[1]); err == nil {
			return Position{Line: line, Col: 1}, true
		}
	}
	return Position{}, false
}

// Snippet renders a caret line under the given 1-based line/column (best effort).
func Snippet(src string, pos Position) string {
	if pos.Line < 1 || pos.Col < 1 || src == "" {
		return ""
	}
	text := strings.ReplaceAll(src, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	if pos.Line > len(lines) {
		return ""
	}
	lineText := lines[pos.Line-1]
	prefix := fmt.Sprintf("%d | ", pos.Line)
	var b strings.Builder
	b.WriteString(prefix)
	b.WriteString(lineText)
	b.WriteByte('\n')
	pad := len(prefix) + pos.Col - 1
	if pad < 0 {
		pad = 0
	}
	b.WriteString(strings.Repeat(" ", pad))
	b.WriteString("^\n")
	return b.String()
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// IdentifierSpanAt returns a best-effort byte span for an identifier starting at 1-based column col.
func IdentifierSpanAt(line string, col int) int {
	if col < 1 || col > len(line) {
		return 1
	}
	i := col - 1
	if !isIdentByte(line[i]) {
		return 1
	}
	j := i
	for j < len(line) && isIdentByte(line[j]) {
		j++
	}
	w := j - i
	if w < 1 {
		return 1
	}
	return w
}

// RustStyleSnippet appends a rustc-like source excerpt when path is readable and line exists.
func RustStyleSnippet(b *strings.Builder, path string, line, col int, hint string) bool {
	if path == "" || line < 1 || col < 1 {
		return false
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	src := normalizeNewlines(string(raw))
	lines := strings.Split(src, "\n")
	if line > len(lines) {
		return false
	}
	lineText := lines[line-1]
	span := IdentifierSpanAt(lineText, col)
	ln := strconv.Itoa(line)
	w := len(ln)
	fmt.Fprintf(b, "\n%*d | %s", w, line, lineText)
	under := strings.Repeat(" ", col-1) + strings.Repeat("^", span)
	if hint != "" {
		under += " " + hint
	}
	fmt.Fprintf(b, "\n%*s | %s", w, "", under)
	return true
}

// SourceContextError attaches file contents so compile failures can print a snippet.
type SourceContextError struct {
	Path   string
	Source string
	Phase  string // "lexer" or "parse"
	Cause  error
}

func (e *SourceContextError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s: %s: %v", e.Path, e.Phase, e.Cause)
	if pos, ok := ExtractPosition(e.Cause); ok {
		if sn := Snippet(e.Source, pos); sn != "" {
			b.WriteByte('\n')
			b.WriteString(sn)
		}
	}
	return b.String()
}

func (e *SourceContextError) Unwrap() error { return e.Cause }

// WrapLexer wraps a lexer failure with path and source for richer CLI output.
func WrapLexer(path, src string, err error) error {
	return &SourceContextError{Path: path, Source: src, Phase: "lexer", Cause: err}
}

// WrapParse wraps a parser failure with path and source.
func WrapParse(path, src string, err error) error {
	return &SourceContextError{Path: path, Source: src, Phase: "parse", Cause: err}
}

// LevenshteinDistance computes edit distance between two strings.
func LevenshteinDistance(a, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			del := prev[j] + 1
			ins := cur[j-1] + 1
			sub := prev[j-1] + cost
			cur[j] = del
			if ins < cur[j] {
				cur[j] = ins
			}
			if sub < cur[j] {
				cur[j] = sub
			}
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

// BestSuggestion returns the closest candidate (<= maxDist edits).
func BestSuggestion(target string, candidates []string, maxDist int) (string, bool) {
	best := ""
	bestDist := maxDist + 1
	for _, c := range candidates {
		if c == "" || c == target {
			continue
		}
		d := LevenshteinDistance(strings.ToLower(target), strings.ToLower(c))
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	if best == "" || bestDist > maxDist {
		return "", false
	}
	return best, true
}

// FormatError renders rich user-facing output when diagnostic data exists.
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	var multi *MultiError
	if errors.As(err, &multi) && multi != nil && len(multi.List) > 0 {
		var b strings.Builder
		if multi.Label != "" {
			fmt.Fprintf(&b, "%d errors in %s\n\n", len(multi.List), multi.Label)
		} else {
			fmt.Fprintf(&b, "%d errors\n\n", len(multi.List))
		}
		for i, child := range multi.List {
			if i > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString(FormatError(child))
		}
		return b.String()
	}
	var d *DiagnosticError
	if errors.As(err, &d) && d != nil {
		var b strings.Builder
		b.WriteString(d.Message)
		if d.File != "" && d.Line > 0 && d.Col > 0 {
			fmt.Fprintf(&b, "\n --> %s:%d:%d", d.File, d.Line, d.Col)
			if RustStyleSnippet(&b, d.File, d.Line, d.Col, d.Hint) {
				return b.String()
			}
		}
		if d.Hint != "" {
			b.WriteString("\n  hint: ")
			b.WriteString(d.Hint)
		}
		return b.String()
	}
	var s *SourceContextError
	if errors.As(err, &s) && s != nil {
		if pos, ok := ExtractPosition(s.Cause); ok {
			var b strings.Builder
			b.WriteString(s.Cause.Error())
			fmt.Fprintf(&b, "\n  --> %s:%d:%d", s.Path, pos.Line, pos.Col)
			if sn := Snippet(s.Source, pos); sn != "" {
				b.WriteString("\n")
				b.WriteString(sn)
			}
			return b.String()
		}
	}
	return err.Error()
}

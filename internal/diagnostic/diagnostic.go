package diagnostic

import (
	"errors"
	"fmt"
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

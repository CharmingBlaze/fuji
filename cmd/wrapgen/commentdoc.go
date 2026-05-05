package main

import (
	"strings"
)

// docCommentBeforeFunction returns text from // or /* */ comments immediately
// above the first occurrence of funcName( in raw source.
func docCommentBeforeFunction(raw, funcName string) string {
	idx := strings.Index(raw, funcName+"(")
	if idx < 0 {
		return ""
	}
	prefix := raw[:idx]
	lines := strings.Split(prefix, "\n")
	var block []string
	for i := len(lines) - 1; i >= 0; i-- {
		s := strings.TrimSpace(lines[i])
		if s == "" {
			if len(block) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(s, "//") {
			t := strings.TrimSpace(strings.TrimPrefix(s, "//"))
			block = append([]string{t}, block...)
			continue
		}
		if strings.HasPrefix(s, "/*") {
			end := strings.Index(s, "*/")
			if end > 2 {
				return strings.TrimSpace(s[2:end])
			}
		}
		break
	}
	return strings.TrimSpace(strings.Join(block, " "))
}

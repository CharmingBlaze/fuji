package parser

import (
	"strings"

	"fuji/internal/lexer"
)

// normalizeIdentLexeme folds an identifier to ASCII lowercase for names,
// property keys, and lookups (Fuji is case-insensitive for identifiers).
func normalizeIdentLexeme(t lexer.Token) lexer.Token {
	t.Lexeme = strings.ToLower(t.Lexeme)
	return t
}

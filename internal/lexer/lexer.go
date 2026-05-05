package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	source    []byte
	tokens    []Token
	start     int
	current   int
	line      int
	lineStart int
}

func NewLexer(source string) *Lexer {
	src := []byte(source)
	// Handle BOM
	if len(src) >= 3 && src[0] == 0xEF && src[1] == 0xBB && src[2] == 0xBF {
		src = src[3:]
	}
	return &Lexer{
		source: src,
		line:   1,
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	for !l.isAtEnd() {
		l.start = l.current
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}

	l.addToken(TokenEOF)
	return l.tokens, nil
}

func (l *Lexer) scanToken() error {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(TokenLParen)
	case ')':
		l.addToken(TokenRParen)
	case '{':
		l.addToken(TokenLBrace)
	case '}':
		l.addToken(TokenRBrace)
	case '[':
		l.addToken(TokenLBracket)
	case ']':
		l.addToken(TokenRBracket)
	case ',':
		l.addToken(TokenComma)
	case '.':
		if l.match('.') {
			if l.match('.') {
				l.addToken(TokenTripleDot)
			} else {
				l.addToken(TokenDotDot)
			}
		} else {
			l.addToken(TokenDot)
		}
	case ';':
		l.addToken(TokenSemicolon)
	case ':':
		l.addToken(TokenColon)
	case '?':
		l.addToken(TokenQuestion)
	case '-':
		if l.match('-') {
			l.addToken(TokenMinusMinus)
		} else if l.match('=') {
			l.addToken(TokenMinusEqual)
		} else {
			l.addToken(TokenMinus)
		}
	case '+':
		if l.match('+') {
			l.addToken(TokenPlusPlus)
		} else if l.match('=') {
			l.addToken(TokenPlusEqual)
		} else {
			l.addToken(TokenPlus)
		}
	case '*':
		if l.match('=') {
			l.addToken(TokenStarEqual)
		} else {
			l.addToken(TokenStar)
		}
	case '/':
		if l.match('/') {
			start := l.current
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
			comment := string(l.source[start:l.current])
			if strings.HasPrefix(comment, " fuji:") {
				l.addTokenWithLexeme(TokenComment, "fuji:"+strings.TrimPrefix(comment, " fuji:"))
			}
		} else if l.match('*') {
			for !l.isAtEnd() {
				if l.peek() == '*' && l.peekNext() == '/' {
					l.advance()
					l.advance()
					break
				}
				if l.peek() == '\n' {
					l.line++
					l.lineStart = l.current + 1
				}
				l.advance()
			}
		} else if l.match('=') {
			l.addToken(TokenSlashEqual)
		} else {
			l.addToken(TokenSlash)
		}
	case '%':
		if l.match('=') {
			l.addToken(TokenPercentEqual)
		} else {
			l.addToken(TokenPercent)
		}
	case '&':
		if l.match('&') {
			l.addToken(TokenAndAnd)
		} else if l.match('=') {
			l.addToken(TokenAndEqual)
		} else {
			l.addToken(TokenAnd)
		}
	case '|':
		if l.match('|') {
			l.addToken(TokenOrOr)
		} else if l.match('=') {
			l.addToken(TokenOrEqual)
		} else {
			l.addToken(TokenOr)
		}
	case '^':
		if l.match('=') {
			l.addToken(TokenCaretEqual)
		} else {
			l.addToken(TokenCaret)
		}
	case '~':
		l.addToken(TokenTilde)
	case '!':
		if l.match('=') {
			l.addToken(TokenBangEqual)
		} else {
			l.addToken(TokenBang)
		}
	case '=':
		if l.match('=') {
			l.addToken(TokenEqualEqual)
		} else if l.match('>') {
			l.addToken(TokenArrow)
		} else {
			l.addToken(TokenEqual)
		}
	case '<':
		if l.match('=') {
			l.addToken(TokenLessEqual)
		} else if l.match('<') {
			if l.match('=') {
				l.addToken(TokenLessLessEqual)
			} else {
				l.addToken(TokenLessLess)
			}
		} else {
			l.addToken(TokenLess)
		}
	case '>':
		if l.match('=') {
			l.addToken(TokenGreaterEqual)
		} else if l.match('>') {
			if l.match('=') {
				l.addToken(TokenGreaterGreaterEqual)
			} else {
				l.addToken(TokenGreaterGreater)
			}
		} else {
			l.addToken(TokenGreater)
		}
	case '"':
		return l.string()
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		l.line++
		l.lineStart = l.current
	case '#':
		return l.directive()
	default:
		if isDigit(c) {
			return l.number()
		} else if isAlpha(c) {
			return l.identifier()
		} else {
			return fmt.Errorf("unexpected character at %d:%d: %c", l.line, l.current-l.lineStart, c)
		}
	}
	return nil
}

func (l *Lexer) string() error {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
			l.lineStart = l.current + 1
		}
		if l.peek() == '\\' {
			l.advance()
		}
		l.advance()
	}

	if l.isAtEnd() {
		return fmt.Errorf("unterminated string at %d", l.line)
	}

	l.advance() // The closing "
	l.addToken(TokenString)
	return nil
}

func (l *Lexer) number() error {
	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // Consume the "."
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	l.addToken(TokenNumber)
	return nil
}

func (l *Lexer) identifier() error {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := string(l.source[l.start:l.current])
	typ := l.lookupKeyword(text)
	l.addToken(typ)
	return nil
}

func (l *Lexer) directive() error {
	for isAlpha(l.peek()) {
		l.advance()
	}
	text := string(l.source[l.start:l.current])
	if text == "#include" {
		l.addToken(TokenInclude)
	} else {
		l.addToken(TokenError)
	}
	return nil
}

func (l *Lexer) lookupKeyword(text string) TokenType {
	switch text {
	case "break":
		return TokenBreak
	case "case":
		return TokenCase
	case "continue":
		return TokenContinue
	case "default":
		return TokenDefault
	case "delete":
		return TokenDelete
	case "do":
		return TokenDo
	case "else":
		return TokenElse
	case "false":
		return TokenFalse
	case "for":
		return TokenFor
	case "func":
		return TokenFunc
	case "if":
		return TokenIf
	case "import":
		return TokenImport
	case "in":
		return TokenIn
	case "let":
		return TokenLet
	case "var":
		return TokenVar
	case "null":
		return TokenNull
	case "of":
		return TokenOf
	case "return":
		return TokenReturn
	case "switch":
		return TokenSwitch
	case "this":
		return TokenThis
	case "true":
		return TokenTrue
	case "while":
		return TokenWhile
	default:
		return TokenIdentifier
	}
}

func (l *Lexer) advance() byte {
	c := l.source[l.current]
	l.current++
	return c
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) addToken(typ TokenType) {
	text := string(l.source[l.start:l.current])
	l.addTokenWithLexeme(typ, text)
}

func (l *Lexer) addTokenWithLexeme(typ TokenType, lexeme string) {
	l.tokens = append(l.tokens, Token{
		Type:   typ,
		Lexeme: lexeme,
		Line:   l.line,
		Col:    l.start - l.lineStart + 1,
	})
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

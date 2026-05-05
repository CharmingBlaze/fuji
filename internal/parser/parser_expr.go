package parser

import (
	"fmt"
	"fuji/internal/lexer"
	"strconv"
	"strings"
)

const (
	_ int = iota
	PrecedenceLowest
	PrecedenceAssign      // =
	PrecedenceOr          // ||
	PrecedenceAnd         // &&
	PrecedenceBitOr       // |
	PrecedenceBitXor      // ^
	PrecedenceBitAnd      // &
	PrecedenceEquals      // == != === !==
	PrecedenceLessGreater // < > <= >=
	PrecedenceSum         // + -
	PrecedenceShift       // << >> >>>
	PrecedenceProduct     // * / %
	PrecedencePrefix      // -X or !X
	PrecedenceCall        // myFunction(X)
	PrecedenceIndex       // array[index]
)

var precedences = map[lexer.TokenType]int{
	lexer.TokenEqual:          PrecedenceAssign,
	lexer.TokenPlusEqual:      PrecedenceAssign,
	lexer.TokenMinusEqual:     PrecedenceAssign,
	lexer.TokenStarEqual:      PrecedenceAssign,
	lexer.TokenSlashEqual:     PrecedenceAssign,
	lexer.TokenPercentEqual:   PrecedenceAssign,
	lexer.TokenAndEqual:       PrecedenceAssign,
	lexer.TokenOrEqual:        PrecedenceAssign,
	lexer.TokenCaretEqual:     PrecedenceAssign,
	lexer.TokenLessLessEqual:  PrecedenceAssign,
	lexer.TokenGreaterGreaterEqual: PrecedenceAssign,
	lexer.TokenEqualEqual:     PrecedenceEquals,
	lexer.TokenStrictEqual:    PrecedenceEquals,
	lexer.TokenBangEqual:      PrecedenceEquals,
	lexer.TokenStrictNotEqual: PrecedenceEquals,
	lexer.TokenLess:           PrecedenceLessGreater,
	lexer.TokenLessEqual:      PrecedenceLessGreater,
	lexer.TokenGreater:        PrecedenceLessGreater,
	lexer.TokenGreaterEqual:   PrecedenceLessGreater,
	lexer.TokenPlus:           PrecedenceSum,
	lexer.TokenMinus:          PrecedenceSum,
	lexer.TokenLessLess:       PrecedenceShift,
	lexer.TokenGreaterGreater: PrecedenceShift,
	lexer.TokenUnsignedShift:  PrecedenceShift,
	lexer.TokenSlash:          PrecedenceProduct,
	lexer.TokenStar:           PrecedenceProduct,
	lexer.TokenPercent:        PrecedenceProduct,
	lexer.TokenAndAnd:         PrecedenceAnd,
	lexer.TokenOrOr:           PrecedenceOr,
	lexer.TokenAnd:            PrecedenceBitAnd,
	lexer.TokenOr:             PrecedenceBitOr,
	lexer.TokenCaret:          PrecedenceBitXor,
	lexer.TokenLParen:         PrecedenceCall,
	lexer.TokenLBracket:       PrecedenceIndex,
	lexer.TokenDot:            PrecedenceIndex,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peek().Type]; ok {
		return p
	}
	return PrecedenceLowest
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.previous().Type]; ok {
		return p
	}
	return PrecedenceLowest
}

func (p *Parser) parseExpression(precedence int) (Expr, error) {
	token := p.peek()
	prefix := p.getPrefixFn(token.Type)
	if prefix == nil {
		return nil, p.error(token, fmt.Sprintf("no prefix parsing function for %v found", token.Type))
	}

	leftExpr, err := prefix()
	if err != nil {
		return nil, err
	}

	// Postfix ++ / -- (bind tighter than binary operators; same tier as calls).
	for (p.check(lexer.TokenPlusPlus) || p.check(lexer.TokenMinusMinus)) && precedence < PrecedenceCall {
		opTok := p.advance()
		leftExpr = &UpdateExpr{Token: opTok, Operator: opTok, Operand: leftExpr, IsPrefix: false}
	}

	for !p.isAtEnd() && precedence < p.peekPrecedence() {
		token = p.peek()
		infix := p.getInfixFn(token.Type)
		if infix == nil {
			return leftExpr, nil
		}

		p.advance()
		leftExpr, err = infix(leftExpr)
		if err != nil {
			return nil, err
		}
	}

	return leftExpr, nil
}

type (
	prefixParseFn func() (Expr, error)
	infixParseFn  func(Expr) (Expr, error)
)

func (p *Parser) getPrefixFn(typ lexer.TokenType) prefixParseFn {
	switch typ {
	case lexer.TokenIdentifier:
		return p.parseIdentifier
	case lexer.TokenNumber:
		return p.parseNumberLiteral
	case lexer.TokenString:
		return p.parseStringLiteral
	case lexer.TokenImport:
		return p.parseImportExpression
	case lexer.TokenTrue, lexer.TokenFalse:
		return p.parseBooleanLiteral
	case lexer.TokenNull:
		return p.parseNullLiteral
	case lexer.TokenBang, lexer.TokenMinus:
		return p.parsePrefixExpression
	case lexer.TokenLParen:
		return p.parseGroupedExpression
	case lexer.TokenLBrace:
		return p.parseObjectLiteral
	case lexer.TokenLBracket:
		return p.parseArrayLiteral
	case lexer.TokenFunc:
		return p.parseFuncExpression
	case lexer.TokenThis:
		return p.parseThisExpression
	case lexer.TokenVar:
		return p.parseReservedVar
	default:
		return nil
	}
}

func (p *Parser) parseReservedVar() (Expr, error) {
	tok := p.previous()
	return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
}

func (p *Parser) getInfixFn(typ lexer.TokenType) infixParseFn {
	switch typ {
	case lexer.TokenPlus, lexer.TokenMinus, lexer.TokenSlash, lexer.TokenStar, lexer.TokenPercent,
		lexer.TokenEqualEqual, lexer.TokenStrictEqual, lexer.TokenBangEqual, lexer.TokenStrictNotEqual,
		lexer.TokenLess, lexer.TokenLessEqual,
		lexer.TokenGreater, lexer.TokenGreaterEqual, lexer.TokenAndAnd, lexer.TokenOrOr,
		lexer.TokenAnd, lexer.TokenOr, lexer.TokenCaret,
		lexer.TokenLessLess, lexer.TokenGreaterGreater, lexer.TokenUnsignedShift:
		return p.parseInfixExpression
	case lexer.TokenLParen:
		return p.parseCallExpression
	case lexer.TokenEqual, lexer.TokenPlusEqual, lexer.TokenMinusEqual, lexer.TokenStarEqual,
		lexer.TokenSlashEqual, lexer.TokenPercentEqual, lexer.TokenAndEqual, lexer.TokenOrEqual,
		lexer.TokenCaretEqual, lexer.TokenLessLessEqual, lexer.TokenGreaterGreaterEqual:
		return p.parseAssignExpression
	case lexer.TokenDot:
		return p.parseDotExpression
	case lexer.TokenLBracket:
		return p.parseIndexExpression
	default:
		return nil
	}
}

func (p *Parser) parseIdentifier() (Expr, error) {
	token := p.advance()
	return &IdentifierExpr{Token: token, Name: token}, nil
}

func (p *Parser) parseNumberLiteral() (Expr, error) {
	token := p.advance()
	lex := token.Lexeme
	var value float64
	if strings.ContainsAny(lex, "eE.") {
		v, err := strconv.ParseFloat(lex, 64)
		if err != nil {
			return nil, p.error(token, "could not parse number literal")
		}
		value = v
	} else {
		i, err := strconv.ParseInt(lex, 0, 64)
		if err != nil {
			return nil, p.error(token, "could not parse number literal")
		}
		value = float64(i)
	}
	return &LiteralExpr{Token: token, Value: value}, nil
}

func (p *Parser) parseStringLiteral() (Expr, error) {
	token := p.advance()
	s := token.Lexeme
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	return &LiteralExpr{Token: token, Value: s}, nil
}

func (p *Parser) parseBooleanLiteral() (Expr, error) {
	token := p.advance()
	return &LiteralExpr{Token: token, Value: token.Type == lexer.TokenTrue}, nil
}

func (p *Parser) parseNullLiteral() (Expr, error) {
	token := p.advance()
	return &LiteralExpr{Token: token, Value: nil}, nil
}

func (p *Parser) parseImportExpression() (Expr, error) {
	token := p.advance()
	path, err := p.consume(lexer.TokenString, "expected string path after import")
	if err != nil {
		return nil, err
	}
	return &ImportExpr{Token: token, Path: path}, nil
}

func (p *Parser) parsePrefixExpression() (Expr, error) {
	token := p.advance()
	right, err := p.parseExpression(PrecedencePrefix)
	if err != nil {
		return nil, err
	}
	return &PrefixExpr{Token: token, Operator: token.Lexeme, Right: right}, nil
}

func (p *Parser) parseGroupedExpression() (Expr, error) {
	p.advance() // (
	expr, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenRParen, "expected ')' after grouped expression"); err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) parseInfixExpression(left Expr) (Expr, error) {
	token := p.previous()
	precedence := p.curPrecedence()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	return &InfixExpr{Token: token, Left: left, Operator: token.Lexeme, Right: right}, nil
}

func (p *Parser) parseCallExpression(left Expr) (Expr, error) {
	token := p.previous()
	args := []Expr{}

	if !p.check(lexer.TokenRParen) {
		for {
			arg, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after arguments"); err != nil {
		return nil, err
	}

	return &CallExpr{Token: token, Function: left, Arguments: args}, nil
}

func (p *Parser) parseAssignExpression(left Expr) (Expr, error) {
	token := p.previous()
	value, err := p.parseExpression(PrecedenceAssign - 1) // Right associative
	if err != nil {
		return nil, err
	}
	return &AssignExpr{Token: token, Left: left, Value: value}, nil
}

func (p *Parser) parseDotExpression(left Expr) (Expr, error) {
	token := p.previous()
	name, err := p.consume(lexer.TokenIdentifier, "expected property name after '.'")
	if err != nil {
		return nil, err
	}
	// Dot access obj.name is equivalent to obj["name"]
	return &IndexExpr{
		Token:  token,
		Object: left,
		Index:  &LiteralExpr{Token: name, Value: name.Lexeme},
	}, nil
}

func (p *Parser) parseIndexExpression(left Expr) (Expr, error) {
	token := p.previous()
	index, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenRBracket, "expected ']' after index"); err != nil {
		return nil, err
	}
	return &IndexExpr{Token: token, Object: left, Index: index}, nil
}

func (p *Parser) parseArrayLiteral() (Expr, error) {
	token := p.advance()
	elements := []Expr{}
	if !p.check(lexer.TokenRBracket) {
		for {
			el, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			elements = append(elements, el)
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}
	if _, err := p.consume(lexer.TokenRBracket, "expected ']' after array literal"); err != nil {
		return nil, err
	}
	return &ArrayExpr{Token: token, Elements: elements}, nil
}

func (p *Parser) parseObjectLiteral() (Expr, error) {
	token := p.advance()
	keys := []lexer.Token{}
	values := []Expr{}

	if !p.check(lexer.TokenRBrace) {
		for {
			key, err := p.consume(lexer.TokenIdentifier, "expected property name")
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)

			if _, err := p.consume(lexer.TokenColon, "expected ':' after property name"); err != nil {
				return nil, err
			}

			val, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			values = append(values, val)

			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	if _, err := p.consume(lexer.TokenRBrace, "expected '}' after object literal"); err != nil {
		return nil, err
	}

	return &ObjectExpr{Token: token, Keys: keys, Values: values}, nil
}

func (p *Parser) parseFuncExpression() (Expr, error) {
	token := p.advance()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'func'"); err != nil {
		return nil, err
	}

	params := []Param{}
	if !p.check(lexer.TokenRParen) {
		for {
			name, err := p.consume(lexer.TokenIdentifier, "expected parameter name")
			if err != nil {
				return nil, err
			}
			params = append(params, Param{Name: name.Lexeme})

			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after parameters"); err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &FuncExpr{Token: token, Params: params, Body: body}, nil
}

func (p *Parser) parseThisExpression() (Expr, error) {
	return &ThisExpr{Token: p.advance()}, nil
}

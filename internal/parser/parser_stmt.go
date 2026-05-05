package parser

import (
	"fmt"
	"strings"

	"fuji/internal/lexer"
)

func (p *Parser) parseDeclaration() (Decl, error) {
	if p.match(lexer.TokenComment) {
		comment := p.previous().Lexeme
		var rest string
		switch {
		case strings.HasPrefix(comment, "fuji:extern "):
			rest = strings.TrimPrefix(comment, "fuji:extern ")
		}
		if rest != "" {
			parts := strings.Fields(rest)
			if len(parts) >= 2 {
				arity := 0
				if len(parts) >= 3 {
					fmt.Sscanf(parts[2], "%d", &arity)
				}
				p.lastDirective = &NativeDirective{
					BindingName: parts[0],
					Symbol:      parts[1],
					Arity:       arity,
				}
			}
		}
		return p.parseDeclaration()
	}
	if p.match(lexer.TokenInclude) {
		return p.parseIncludeDeclaration()
	}
	if p.match(lexer.TokenVar) {
		tok := p.previous()
		return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
	}
	if p.match(lexer.TokenLet) {
		decl, err := p.parseLetDeclaration()
		if err == nil {
			if let, ok := decl.(*LetDecl); ok {
				let.Native = p.lastDirective
			}
			p.lastDirective = nil
		}
		return decl, err
	}
	if p.match(lexer.TokenFunc) {
		decl, err := p.parseFuncDeclaration()
		if err == nil {
			if f, ok := decl.(*FuncDecl); ok {
				f.Native = p.lastDirective
			}
			p.lastDirective = nil
		}
		return decl, err
	}
	p.lastDirective = nil
	return p.parseStatement()
}

func (p *Parser) parseIncludeDeclaration() (Decl, error) {
	token := p.previous()
	path, err := p.consume(lexer.TokenString, "expected include path")
	if err != nil {
		return nil, err
	}
	return &IncludeDecl{Token: token, Path: path}, nil
}

func (p *Parser) parseLetDeclaration() (Decl, error) {
	token := p.previous()
	name, err := p.consume(lexer.TokenIdentifier, "expected variable name")
	if err != nil {
		if p.check(lexer.TokenVar) {
			tok := p.peek()
			p.advance()
			return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
		}
		return nil, err
	}

	var init Expr
	if p.match(lexer.TokenEqual) {
		init, err = p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after variable declaration"); err != nil {
		return nil, err
	}

	return &LetDecl{Token: token, Name: name, Init: init}, nil
}

func (p *Parser) parseFuncDeclaration() (Decl, error) {
	token := p.previous()
	name, err := p.consume(lexer.TokenIdentifier, "expected function name")
	if err != nil {
		if p.check(lexer.TokenVar) {
			tok := p.peek()
			p.advance()
			return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
		}
		return nil, err
	}

	if _, err := p.consume(lexer.TokenLParen, "expected '(' after function name"); err != nil {
		return nil, err
	}

	params := []Param{}
	if !p.check(lexer.TokenRParen) {
		for {
			isRest := p.match(lexer.TokenTripleDot)
			paramName, err := p.consume(lexer.TokenIdentifier, "expected parameter name")
			if err != nil {
				if p.check(lexer.TokenVar) {
					tok := p.peek()
					p.advance()
					return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
				}
				return nil, err
			}
			param := Param{Name: paramName.Lexeme, IsRest: isRest}
			if p.match(lexer.TokenEqual) {
				if isRest {
					return nil, p.error(paramName, "rest parameter cannot have a default")
				}
				param.Default, err = p.parseExpression(PrecedenceLowest)
				if err != nil {
					return nil, err
				}
			}
			params = append(params, param)
			if isRest && !p.check(lexer.TokenRParen) {
				return nil, p.error(paramName, "rest parameter must be last")
			}

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

	return &FuncDecl{Token: token, Name: name, Params: params, Body: body}, nil
}

func (p *Parser) parseStatement() (Stmt, error) {
	if p.match(lexer.TokenIf) {
		return p.parseIfStatement()
	}
	if p.match(lexer.TokenWhile) {
		return p.parseWhileStatement()
	}
	if p.match(lexer.TokenReturn) {
		return p.parseReturnStatement()
	}
	if p.match(lexer.TokenBreak) {
		token := p.previous()
		if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after break"); err != nil {
			return nil, err
		}
		return &BreakStmt{Token: token}, nil
	}
	if p.match(lexer.TokenContinue) {
		token := p.previous()
		if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after continue"); err != nil {
			return nil, err
		}
		return &ContinueStmt{Token: token}, nil
	}
	if p.check(lexer.TokenLBrace) {
		return p.parseBlockStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseBlockStatement() (*BlockStmt, error) {
	token, err := p.consume(lexer.TokenLBrace, "expected '{' at start of block")
	if err != nil {
		return nil, err
	}
	declarations := []Decl{}

	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		decl, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		if _, ok := decl.(Stmt); !ok {
			return nil, p.error(p.previous(), "expected statement in block")
		}
		declarations = append(declarations, decl)
	}

	if _, err := p.consume(lexer.TokenRBrace, "expected '}' after block"); err != nil {
		return nil, err
	}

	return &BlockStmt{Token: token, Declarations: declarations}, nil
}

func (p *Parser) parseIfStatement() (Stmt, error) {
	token := p.previous()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'if'"); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after condition"); err != nil {
		return nil, err
	}

	thenBranch, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(lexer.TokenElse) {
		elseBranch, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{Token: token, Condition: condition, Then: thenBranch, Else: elseBranch}, nil
}

func (p *Parser) parseWhileStatement() (Stmt, error) {
	token := p.previous()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'while'"); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after condition"); err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{Token: token, Condition: condition, Body: body}, nil
}

func (p *Parser) parseReturnStatement() (Stmt, error) {
	token := p.previous()
	var value Expr
	var err error

	if !p.check(lexer.TokenSemicolon) {
		value, err = p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after return value"); err != nil {
		return nil, err
	}

	return &ReturnStmt{Token: token, Value: value}, nil
}

func (p *Parser) parseExpressionStatement() (Stmt, error) {
	token := p.peek()
	expr, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after expression"); err != nil {
		return nil, err
	}

	return &ExpressionStmt{Token: token, Expr: expr}, nil
}

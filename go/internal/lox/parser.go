package lox

import "fmt"

type Parser struct {
	tokens  []Token
	current int

	lox *Lox
}

type ParseError struct {
	error
}

func NewParser(tokens []Token, l *Lox) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		lox:     l,
	}
}

func NewParseError(err error) *ParseError {
	return &ParseError{err}
}

func (p *Parser) Parse() (expr Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse error: %v", r)
		}
	}()

	return p.expression(), nil
}

// Expressions
func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return NewLiteral(false)
	}
	if p.match(TRUE) {
		return NewLiteral(true)
	}
	if p.match(NIL) {
		return NewLiteral(nil)
	}

	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().literal)
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return NewGrouping(expr)
	}

	panic(p.error(p.peek(), "Expect expression."))
}

// ERRORS
func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

func (p *Parser) error(token Token, message string) *ParseError {
	p.lox.error(token.line, message)
	return NewParseError(
		fmt.Errorf("[line %d] Error at '%s': %s", token.line, token.lexeme, message),
	)
}

// UTILS
func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

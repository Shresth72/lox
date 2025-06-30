package lox

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int

	lox *Lox
}

func NewScanner(source string, lox *Lox) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
		lox:     lox,
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	token := NewToken(EOF, "", nil, s.line)
	s.tokens = append(s.tokens, *token)
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)

	case '!':
		s.addMatchToken('=', BANG_EQUAL, BANG)
	case '=':
		s.addMatchToken('=', EQUAL_EQUAL, EQUAL)
	case '<':
		s.addMatchToken('=', LESS_EQUAL, LESS)
	case '>':
		s.addMatchToken('=', GREATER_EQUAL, GREATER)

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}

	case '"':
		s.captureString()

	case ' ', '\r', '\t':
	case '\n':
		s.line++

	default:
		if s.isDigit(c) {
			s.captureNumber()
		} else if s.isAlpha(c) {
			s.captureIdentifier()
		} else {
			s.lox.error(s.line, fmt.Sprintf("Unexpected character: %q", c))
		}
	}
}

func (s *Scanner) captureString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lox.error(s.line, "Unterminated string")
		return
	}
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) captureNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	text := s.source[s.start:s.current]
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		s.lox.error(s.line, "Invalid number format")
		return
	}
	s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) captureIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	s.addToken(IDENTIFIER)
}

func (s *Scanner) addMatchToken(expected byte, first, second TokenType) {
	if s.match(expected) {
		s.addToken(first)
	} else {
		s.addToken(second)
	}
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	token := NewToken(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, *token)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

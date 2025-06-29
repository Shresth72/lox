package lox

import "fmt"

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
	default:
		s.lox.error(s.line, fmt.Sprintf("Unexpected character: %q", c))
	}
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

func (s *Scanner) peek() byte {
  if s.isAtEnd() {
    return '\0'
  }
  return s.source[s.current]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	lexeme := s.source[s.start:s.current]
	token := NewToken(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, *token)
}

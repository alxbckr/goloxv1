package main

import (
	"errors"
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int

	keywords map[string]TokenType
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:   source,
		tokens:   []Token{},
		start:    0,
		current:  0,
		line:     0,
		keywords: map[string]TokenType{"and": AND, "class": CLASS, "else": ELSE, "false": FALSE, "for": FOR, "fun": FUN, "if": IF, "nil": NIL, "or": OR, "print": PRINT, "return": RETURN, "super": SUPER, "this": THIS, "true": TRUE, "var": VAR, "while": WHILE},
	}
}

func (s *Scanner) scanTokens() ([]Token, error) {
	var tokens []Token
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	tokens = append(tokens, *NewToken(EOF, "", "", s.line))
	return tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
		break
	case ')':
		s.addToken(RIGHT_PAREN)
		break
	case '{':
		s.addToken(LEFT_BRACE)
		break
	case '}':
		s.addToken(RIGHT_BRACE)
		break
	case ',':
		s.addToken(COMMA)
		break
	case '.':
		s.addToken(DOT)
		break
	case '-':
		s.addToken(MINUS)
		break
	case '+':
		s.addToken(PLUS)
		break
	case ';':
		s.addToken(SEMICOLON)
		break
	case '*':
		s.addToken(STAR)
		break
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
		break
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
		break
	case '"':
		err := s.string()
		if err != nil {
			return err
		}
		break
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			return &ScannerError{
				Line:  s.line,
				Where: "",
				Err:   errors.New("Unexpected character."),
			}
		}
	}
	return nil
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, "")
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(tokenType, text, literal, s.line))
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

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return &ScannerError{
			Line:  s.line,
			Where: "",
			Err:   errors.New("Unterminated string."),
		}
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
	return nil
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	d, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addTokenWithLiteral(NUMBER, d)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	s.addToken(IDENTIFIER)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

package lox

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
	return t.TokenType.String() + " " + t.Lexeme + " " + fmt.Sprintf("%v", t.Literal)
}

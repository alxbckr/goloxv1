package main

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t *Token) String() string {
	//return t.tokenType.String() + " " + t.lexeme + " " + t.literal
	return t.tokenType.String() + " " + t.lexeme
}

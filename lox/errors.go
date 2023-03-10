package lox

import (
	"fmt"
)

type ScannerError struct {
	Line  int
	Where string
	Error string
}

type LoxError struct {
	Token   Token
	Message string
}

type RuntimeError struct {
	Token   Token
	Message string
}

func NewScannerError(line int, where string, error string) *ScannerError {
	return &ScannerError{
		Line:  line,
		Where: where,
		Error: error,
	}
}

func NewLoxError(token Token, message string) *LoxError {
	return &LoxError{
		Token:   token,
		Message: message,
	}
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (err *LoxError) Error() string {
	line := err.Token.Line
	where := err.Token.Lexeme
	message := err.Message

	if err.Token.TokenType == EOF {
		where = "end"
	}
	return fmt.Sprintf("[line %v] Error at %v: %v\n", line, where, message)
}

func (err *RuntimeError) Error() string {
	line := err.Token.Line
	message := err.Message
	return fmt.Sprintf("%v [line %v]", message, line)
}

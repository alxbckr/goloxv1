package lox

import (
	"fmt"
	"strconv"
)

type ScannerError struct {
	hadError bool
}

type LoxError struct {
	Token   Token
	Message string
}

var scannerError ScannerError

func GetScannerError() *ScannerError {
	return &scannerError
}

func (s *ScannerError) GetHadError() bool {
	return s.hadError
}

func (s *ScannerError) Reset() {
	s.hadError = false
}

func ReportError(line int, where string, error string) {
	fmt.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + error)
	scannerError.hadError = true
}

func NewLoxError(token Token, message string) *LoxError {
	return &LoxError{
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

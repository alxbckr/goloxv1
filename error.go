package main

import (
	"strconv"
)

type ScannerError struct {
	Line  int
	Where string
	Err   error
}

func (e *ScannerError) Error() string {
	return "[line " + strconv.Itoa(e.Line) + "] Error" + e.Where + ": " + e.Err.Error()
}

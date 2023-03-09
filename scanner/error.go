package scanner

import (
	"fmt"
	"strconv"
)

type ScannerError struct {
	hadError bool
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

package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		fmt.Printf("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

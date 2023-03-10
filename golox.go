package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alxbckr/goloxv1/lox"
)

func run(source string) {
	scan := lox.NewScanner(source)
	tokens := scan.ScanTokens()

	// For now, just print the tokens.
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		if line == "" {
			return
		}
		run(line)
		lox.GetScannerError().Reset()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	run(string(bytes))
	// Indicate an error in the exit code.
	if lox.GetScannerError().GetHadError() {
		os.Exit(65)
	}
}

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

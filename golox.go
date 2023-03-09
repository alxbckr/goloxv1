package main

import (
	"bufio"
	"fmt"
	"os"
)

func report(e ScannerError) {
	fmt.Println(e)
}

func run(source string) error {
	scanner := NewScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}

	// For now, just print the tokens.
	for _, token := range tokens {
		fmt.Println(token)
	}
	return &ScannerError{}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		if line == "" {
			return
		}
		err := run(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = run(string(bytes))
	if err != nil {
		fmt.Println(err)
		os.Exit(65)
	}
}

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

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alxbckr/goloxv1/lox"
	"github.com/alxbckr/goloxv1/printer"
)

func run(source string) error {
	scan := lox.NewScanner(source)

	tokens, err := scan.ScanTokens()
	if err != nil {
		return err
	}

	parser := lox.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return err
	}

	fmt.Println(printer.NewAstPrinter().Print(expr))
	return nil
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
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = run(string(bytes))
	// Indicate an error in the exit code.
	if err != nil {
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

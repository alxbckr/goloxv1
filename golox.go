package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alxbckr/goloxv1/lox"
)

var interpreter *lox.Interpreter

func run(source string) error {
	scan := lox.NewScanner(source)

	tokens, err := scan.ScanTokens()
	if err != nil {
		return err
	}

	parser := lox.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	resolver := lox.NewResolver(interpreter)
	err = resolver.ResolveStatements(statements)
	if err != nil {
		return err
	}

	err = interpreter.Interpret(statements)
	if err != nil {
		return err
	}
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
	interpreter = lox.NewInterpreter()
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

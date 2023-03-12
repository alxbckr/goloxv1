package lox

import "fmt"

type Environment struct {
	Values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		Values: map[string]interface{}{},
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Get(name Token) interface{} {
	if v, ok := e.Values[name.Lexeme]; ok {
		return v
	}
	panic(NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'.", name.Lexeme)))
}

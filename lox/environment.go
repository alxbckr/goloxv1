package lox

import "fmt"

type Environment struct {
	enclosing *Environment
	Values    map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		Values:    map[string]interface{}{},
	}
}

func NewEnvironmentWithEnclosing(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		Values:    map[string]interface{}{},
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Assign(name Token, value interface{}) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	panic(NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'.", name.Lexeme)))
}

func (e *Environment) Get(name Token) interface{} {
	if v, ok := e.Values[name.Lexeme]; ok {
		return v
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	panic(NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'.", name.Lexeme)))
}

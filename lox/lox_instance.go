package lox

import "fmt"

type LoxInstance struct {
	Class  LoxClass
	Fields map[string]interface{}
}

func NewLoxInstance(class LoxClass) *LoxInstance {
	return &LoxInstance{
		Class:  class,
		Fields: map[string]interface{}{},
	}
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("%v instance", i.Class.Name)
}

func (i *LoxInstance) Get(name Token) interface{} {
	if f, ok := i.Fields[name.Lexeme]; ok {
		return f
	}

	method := i.Class.FindMethod(name.Lexeme)
	if method != nil {
		return method
	}

	panic(NewRuntimeError(name, fmt.Sprintf("undefined property %v .", name.Lexeme)))
}

func (i *LoxInstance) Set(name Token, value interface{}) {
	i.Fields[name.Lexeme] = value
}

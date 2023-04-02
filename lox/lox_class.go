package lox

import "fmt"

type LoxClass struct {
	Name    string
	Methods map[string]LoxFunction
}

func NewLoxClass(name string, methods map[string]LoxFunction) *LoxClass {
	return &LoxClass{
		Name:    name,
		Methods: methods,
	}
}

func (c LoxClass) String() string {
	return fmt.Sprintf("%v", c.Name)
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return NewLoxInstance(*c)
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if f, ok := c.Methods[name]; ok {
		return &f
	}
	return nil
}

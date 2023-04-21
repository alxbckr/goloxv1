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
	loxInstance := NewLoxInstance(*c)
	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(loxInstance).Call(interpreter, arguments)
	}
	return loxInstance
}

func (c *LoxClass) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if f, ok := c.Methods[name]; ok {
		return &f
	}
	return nil
}

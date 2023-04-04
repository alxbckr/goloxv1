package lox

import "fmt"

type LoxFunction struct {
	Declaration Function
	Closure     *Environment
}

func NewLoxFunction(declaration Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		Declaration: declaration,
		Closure:     closure,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (retVal interface{}) {
	environment := NewEnvironmentWithEnclosing(f.Closure)
	for i, param := range f.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		val := recover()
		if wrapper, ok := val.(*ReturnWrapper); ok && wrapper != nil {
			retVal = wrapper.Value
			return
		}
		panic(val)
	}()

	interpreter.executeBlock(f.Declaration.Body, environment)
	return nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f LoxFunction) String() string {
	return fmt.Sprintf("<fn %v>", f.Declaration.Name.Lexeme)
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironmentWithEnclosing(f.Closure)
	environment.Define("this", instance)
	return NewLoxFunction(f.Declaration, environment)
}

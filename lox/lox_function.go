package lox

import "fmt"

type LoxFunction struct {
	Declaration Function
}

func NewLoxFunction(declaration Function) *LoxFunction {
	return &LoxFunction{
		Declaration: declaration,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (retVal interface{}) {
	environment := NewEnvironmentWithEnclosing(interpreter.globals)
	for i, param := range f.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		val := recover()
		retVal = val.(*ReturnWrapper).Value
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

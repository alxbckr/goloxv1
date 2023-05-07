package lox

import "fmt"

type LoxFunction struct {
	Declaration   Function
	Closure       *Environment
	isInitializer bool
}

func NewLoxFunction(declaration Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		Declaration:   declaration,
		Closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (retVal interface{}) {
	environment := NewEnvironmentWithEnclosing(f.Closure)
	for i, param := range f.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		val := recover()
		if val == nil {
			retVal = nil
			return
		}
		if wrapper, ok := val.(*ReturnWrapper); ok && wrapper != nil {
			if f.isInitializer {
				retVal = f.Closure.GetAt(0, "this")
			} else {
				retVal = wrapper.Value
			}
			return
		}
		panic(val)
	}()

	interpreter.executeBlock(f.Declaration.Body, environment)
	if f.isInitializer {
		return f.Closure.GetAt(0, "this")
	}
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
	return NewLoxFunction(f.Declaration, environment, f.isInitializer)
}

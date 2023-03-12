package lox

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	hadRuntimeError bool
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		hadRuntimeError: false,
	}
}

func (i *Interpreter) Interpret(expression Expr) (err error) {
	defer func() {
		if val := recover(); val != nil {
			runtimeError := val.(*RuntimeError)
			fmt.Println(runtimeError.Error())
			err = runtimeError
			i.hadRuntimeError = true
		}
	}()
	value := i.evaluate(expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}
	// unreachable
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case PLUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		if reflect.TypeOf(left).Kind() == reflect.String && reflect.TypeOf(right).Kind() == reflect.String {
			return left.(string) + right.(string)
		}
		panic(NewRuntimeError(expr.Operator, "operands must be two nubmers or two strings"))
	case SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}
	// unreachable
	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	}
	return true
}

func isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return reflect.ValueOf(left).Equal(reflect.ValueOf(right))
}

func checkNumberOperand(operator Token, operand interface{}) {
	_, ok := operand.(float64)
	if !ok {
		panic(NewRuntimeError(operator, "operand must be a nubmer"))
	}
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) {
	_, ok := left.(float64)
	if !ok {
		panic(NewRuntimeError(operator, "operands must be nubmers"))
	}
	_, ok = right.(float64)
	if !ok {
		panic(NewRuntimeError(operator, "operands must be nubmers"))
	}

}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	v, ok := object.(float64)
	if ok {
		return fmt.Sprintf("%f", v)
	}

	return fmt.Sprintf("%v", object)
}

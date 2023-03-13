package lox

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	hadRuntimeError bool
	environment     Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		hadRuntimeError: false,
		environment:     *NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) (err error) {
	defer func() {
		if val := recover(); val != nil {
			runtimeError := val.(*RuntimeError)
			fmt.Println(runtimeError.Error())
			err = runtimeError
			i.hadRuntimeError = true
		}
	}()
	for _, s := range statements {
		i.execute(s)
	}
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt Expression) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt Print) {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
}

func (i *Interpreter) VisitVarStmt(stmt Var) {
	var value interface{} = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitAssignExpr(expr Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
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

func (i *Interpreter) VisitVariableExpr(expr Variable) interface{} {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
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

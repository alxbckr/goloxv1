package lox

import (
	"fmt"
	"reflect"
	"time"
)

type Interpreter struct {
	hadRuntimeError bool
	globals         *Environment
	environment     *Environment
	locals          map[Expr]int
}

func NewInterpreter() *Interpreter {
	env := NewEnvironment()

	env.Define("clock", NewProtoCallable(0, func(interpreter *Interpreter, arguments []interface{}) interface{} {
		return time.Now().UnixMilli() / 1000.0
	}))

	return &Interpreter{
		hadRuntimeError: false,
		globals:         env,
		environment:     env,
		locals:          make(map[Expr]int),
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

func (i *Interpreter) VisitFunctionStmt(stmt Function) {
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
}

func (i *Interpreter) VisitPrintStmt(stmt Print) {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
}

func (i *Interpreter) VisitReturnStmt(stmt Return) {
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(NewReturnWrapper(value))
}

func (i *Interpreter) VisitVarStmt(stmt Var) {
	var value interface{} = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitWhileStmt(stmt While) {
	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

func (i *Interpreter) VisitBlockStmt(stmt Block) {
	i.executeBlock(stmt.Statements, NewEnvironmentWithEnclosing(i.environment))
}

func (i *Interpreter) VisitClassStmt(stmt Class) {
	i.environment.Define(stmt.Name.Lexeme, nil)

	methods := make(map[string]LoxFunction)
	for _, method := range stmt.Methods {
		function := NewLoxFunction(method, i.environment, (method.Name.Lexeme == "init"))
		methods[method.Name.Lexeme] = *function
	}

	class := NewLoxClass(stmt.Name.Lexeme, methods)
	i.environment.Assign(stmt.Name, class)
}

func (i *Interpreter) VisitAssignExpr(expr Assign) interface{} {
	value := i.evaluate(expr.Value)

	distance, ok := i.locals[&expr]
	if ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}
	return value
}

func (i *Interpreter) VisitIfStmt(stmt If) {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.TokenType == OR {
		if isTruthy(left) {
			return left
		}
	} else if !isTruthy(left) {
		return left
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSetExpr(expr Set) interface{} {
	object := i.evaluate(expr.Object)

	obj, ok := (object).(*LoxInstance)
	if !ok {
		panic(NewRuntimeError(expr.Name, "only instances have fields."))
	}

	value := i.evaluate(expr.Value)
	obj.Set(expr.Name, value)
	return value
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

func (i *Interpreter) VisitCallExpr(expr Call) interface{} {
	callee := i.evaluate(expr.Callee)

	var arguments []interface{}
	for _, a := range expr.Arguments {
		arguments = append(arguments, i.evaluate(a))
	}

	f, ok := callee.(Callable)
	if !ok {
		panic(NewRuntimeError(expr.Paren, "can only call functions and classes"))
	}

	if len(arguments) != f.Arity() {
		panic(NewRuntimeError(expr.Paren, fmt.Sprintf("expected %v arguments but got %v.", f.Arity(), len(arguments))))
	}

	return f.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr Get) interface{} {
	object := i.evaluate(expr.Object)
	if inst, ok := object.(*LoxInstance); ok {
		return inst.Get(expr.Name)
	}
	panic(NewRuntimeError(expr.Name, "only instances have properties"))
}

func (i *Interpreter) VisitVariableExpr(expr Variable) interface{} {
	return i.lookUpVariable(expr.Name, &expr)
}

func (i *Interpreter) VisitThisExpr(expr This) interface{} {
	return i.lookUpVariable(expr.Keyword, &expr)
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) interface{} {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.environment.Get(name)
	}
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(stmt []Stmt, environment *Environment) {
	previous := i.environment

	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	for _, s := range stmt {
		i.execute(s)
	}
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

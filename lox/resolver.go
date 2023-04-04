package lox

import (
	"fmt"

	lls "github.com/emirpasic/gods/stacks/linkedliststack"
)

type FunctionType int
type ClassType int

const (
	NONE FunctionType = iota
	FUNCTION
	METHOD
)

const (
	CLASS_NONE ClassType = iota
	CLASS_CLASS
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          lls.Stack
	currentFunction FunctionType
	currentClass    ClassType
	hadRuntimeError bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          *lls.New(),
		currentFunction: NONE,
		currentClass:    CLASS_NONE,
		hadRuntimeError: false,
	}
}

func (r *Resolver) VisitBlockStmt(stmt Block) {
	r.beginScope()
	r.ResolveStatements(stmt.Statements)
	r.endScope()
}

func (r *Resolver) VisitClassStmt(stmt Class) {
	enclosingClass := r.currentClass
	r.currentClass = CLASS_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	scope, _ := r.scopes.Peek()
	scope.(map[string]bool)["this"] = true

	for _, method := range stmt.Methods {
		declaration := METHOD
		r.resolveFunction(method, declaration)
	}

	r.currentClass = enclosingClass
	r.endScope()
}

func (r *Resolver) VisitExpressionStmt(stmt Expression) {
	r.resolveExpression(stmt.Expression)
}

func (r *Resolver) VisitVarStmt(stmt Var) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
}

func (r *Resolver) VisitFunctionStmt(stmt Function) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) VisitIfStmt(stmt If) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
}

func (r *Resolver) VisitPrintStmt(stmt Print) {
	r.resolveExpression(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt Return) {
	if r.currentFunction == NONE {
		panic(NewLoxError(stmt.Keyword, "can't return from top-level code."))
	}

	if stmt.Value != nil {
		r.resolveExpression(stmt.Value)
	}
}

func (r *Resolver) VisitWhileStmt(stmt While) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
}

func (r *Resolver) VisitAssignExpr(expr Assign) interface{} {
	r.resolveExpression(expr.Value)
	r.resolveLocal(&expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr Binary) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr Call) interface{} {
	r.resolveExpression(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpression(arg)
	}
	return nil
}

func (r *Resolver) VisitGetExpr(expr Get) interface{} {
	r.resolveExpression(expr.Object)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr Grouping) interface{} {
	r.resolveExpression(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr Logical) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitSetExpr(expr Set) interface{} {
	r.resolveExpression(expr.Value)
	r.resolveExpression(expr.Object)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr Unary) interface{} {
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr Variable) interface{} {
	if !r.scopes.Empty() {
		scope, _ := r.scopes.Peek()
		if !(scope.(map[string]bool))[expr.Name.Lexeme] {
			panic(NewLoxError(expr.Name, "can't read local variabl in its own initializer."))
		}
	}

	r.resolveLocal(&expr, expr.Name)
	return nil
}

func (r *Resolver) VisitThisExpr(expr This) interface{} {
	if r.currentClass == CLASS_NONE {
		panic(NewLoxError(expr.Keyword, "can't use 'this' outside of a class"))
	}
	r.resolveLocal(&expr, expr.Keyword)
	return nil
}

func (r *Resolver) ResolveStatements(statements []Stmt) (err error) {
	defer func() {
		if val := recover(); val != nil {
			loxError := val.(*LoxError)
			fmt.Println(loxError.Error())
			err = loxError
			r.hadRuntimeError = true
		}
	}()

	for _, s := range statements {
		r.resolveStatement(s)
	}
	return nil
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name Token) {
	if r.scopes.Empty() {
		return
	}

	scope, _ := r.scopes.Peek()

	if _, ok := scope.(map[string]bool)[name.Lexeme]; ok {
		panic(NewLoxError(name, "already a variable with this name in this scope."))
	}

	(scope.(map[string]bool))[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.scopes.Empty() {
		return
	}
	scope, _ := r.scopes.Peek()
	(scope.(map[string]bool))[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	iter := r.scopes.Iterator()
	scopeDeep := 0
	for iter.Next() {
		if _, ok := (iter.Value().(map[string]bool))[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, scopeDeep)
		}
		scopeDeep++
	}
}

func (r *Resolver) resolveFunction(function Function, typeF FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typeF

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

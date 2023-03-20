package lox

import lls "github.com/emirpasic/gods/stacks/linkedliststack"

type Resolver struct {
	interpreter Interpreter
	scopes      lls.Stack
}

func NewResolver(interpreter Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      *lls.New(),
	}
}

func (r *Resolver) VisitBlockStmt(stmt Block) {
	r.beginScope()
	r.resolveStatements(stmt.Statements)
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
	r.resolveFunction(stmt)
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

func (r *Resolver) VisitVariableExpr(expr Variable) interface{} {
	if !r.scopes.Empty() {
		scope, _ := r.scopes.Peek()
		if (scope.(map[string]bool))[expr.Name.Lexeme] == false {
			panic(NewLoxError(expr.Name, "can't read local variabl in its own initializer."))
		}
	}

	r.resolveLocal(&expr, expr.Name)
	return nil
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, s := range statements {
		r.resolveStatement(s)
	}
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
			r.interpreter.resolve(expr, scopeDeep)
		}
		scopeDeep++
	}
}

func (r *Resolver) resolveFunction(function Function) {
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.Body)
	r.endScope()
}

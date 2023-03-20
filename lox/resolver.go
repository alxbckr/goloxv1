package lox

import lls "github.com/emirpasic/gods/stacks/linkedliststack"

type Resolver struct {
	interpreter Interpreter
	scopes      lls.Stack
}

func NewResolver(interpreter Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
	}
}

func (r *Resolver) VisitBlockStmt(stmt Block) {
	r.beginScope()
	r.resolveStatements(stmt.Statements)
	r.endScope()
}

func (r *Resolver) VisitVarStmt(stmt Var) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve(stmt.Initializer)
	}
	r.define(stmt.Name)
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
	scope[name.Lexeme] = false
}

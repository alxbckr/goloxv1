package lox

type StatementVisitor interface {
	VisitPrintStmt(stmt Print)
	VisitExpressionStmt(stmt Expression)
	VisitVarStmt(stmt Var)
	VisitBlockStmt(stmt Block)
	VisitIfStmt(stmt If)
	VisitWhileStmt(stmt While)
	VisitFunctionStmt(stmt Function)
	VisitReturnStmt(stmt Return)
	VisitClassStmt(stmt Class)
}

type Stmt interface {
	Accept(visitor StatementVisitor)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type Block struct {
	Statements []Stmt
}

type Expression struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}

type Var struct {
	Name        Token
	Initializer Expr
}

type While struct {
	Condition Expr
	Body      Stmt
}

type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

type Return struct {
	Keyword Token
	Value   Expr
}

type Class struct {
	Name       Token
	Superclass *Variable
	Methods    []Function
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
	return &If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		Statements: statements,
	}
}

func NewExpression(expr Expr) *Expression {
	return &Expression{
		Expression: expr,
	}
}

func NewPrint(expr Expr) *Print {
	return &Print{
		Expression: expr,
	}
}

func NewVar(name Token, iniitializer Expr) *Var {
	return &Var{
		Name:        name,
		Initializer: iniitializer,
	}
}

func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		Condition: condition,
		Body:      body,
	}
}

func NewFunction(name Token, params []Token, body []Stmt) *Function {
	return &Function{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func NewReturn(keyword Token, value Expr) *Return {
	return &Return{
		Keyword: keyword,
		Value:   value,
	}
}

func NewClass(name Token, superclass *Variable, methods []Function) *Class {
	return &Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (i *If) Accept(visitor StatementVisitor) {
	visitor.VisitIfStmt(*i)
}

func (b *Block) Accept(visitor StatementVisitor) {
	visitor.VisitBlockStmt(*b)
}

func (s *Expression) Accept(visitor StatementVisitor) {
	visitor.VisitExpressionStmt(*s)
}

func (p *Print) Accept(visitor StatementVisitor) {
	visitor.VisitPrintStmt(*p)
}

func (v *Var) Accept(visitor StatementVisitor) {
	visitor.VisitVarStmt(*v)
}

func (w *While) Accept(visitor StatementVisitor) {
	visitor.VisitWhileStmt(*w)
}

func (f *Function) Accept(visitor StatementVisitor) {
	visitor.VisitFunctionStmt(*f)
}

func (r *Return) Accept(visitor StatementVisitor) {
	visitor.VisitReturnStmt(*r)
}

func (c *Class) Accept(visitor StatementVisitor) {
	visitor.VisitClassStmt(*c)
}

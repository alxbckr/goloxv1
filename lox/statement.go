package lox

type StatementVisitor interface {
	VisitPrintStmt(stmt Print)
	VisitExpressionStmt(stmt Expression)
	VisitVarStmt(stmt Var)
}

type Stmt interface {
	Accept(visitor StatementVisitor)
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

func (s *Expression) Accept(visitor StatementVisitor) {
	visitor.VisitExpressionStmt(*s)
}

func (p *Print) Accept(visitor StatementVisitor) {
	visitor.VisitPrintStmt(*p)
}

func (v *Var) Accept(visitor StatementVisitor) {
	visitor.VisitVarStmt(*v)
}

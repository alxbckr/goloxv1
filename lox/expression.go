package lox

type ExpressionVisitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitLiteralExpr(expr Literal) interface{}
	VisitUnaryExpr(expr Unary) interface{}
	VisitVariableExpr(expr Variable) interface{}
}

type Expr interface {
	Accept(visitor ExpressionVisitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value interface{}
}

type Unary struct {
	Operator Token
	Right    Expr
}

type Variable struct {
	Name Token
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (b *Binary) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitBinaryExpr(*b)
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}

func (g *Grouping) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitGroupingExpr(*g)
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{
		Value: value,
	}
}

func (l *Literal) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteralExpr(*l)
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (u *Unary) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitUnaryExpr(*u)
}

func NewVariable(name Token) *Variable {
	return &Variable{
		Name: name,
	}
}

func (v *Variable) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitVariableExpr(*v)
}

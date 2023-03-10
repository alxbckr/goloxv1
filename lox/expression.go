package lox

type Visitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitLiteralExpr(expr Literal) interface{}
	VisitUnaryExpr(expr Unary) interface{}
}

type Expr interface {
	Accept(visitor Visitor) interface{}
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

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (b *Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpr(*b)
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}

func (g *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpr(*g)
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{
		Value: value,
	}
}

func (l *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpr(*l)
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (u *Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(*u)
}

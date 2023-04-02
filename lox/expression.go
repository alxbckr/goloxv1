package lox

type ExpressionVisitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitCallExpr(expr Call) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitLiteralExpr(expr Literal) interface{}
	VisitLogicalExpr(expr Logical) interface{}
	VisitUnaryExpr(expr Unary) interface{}
	VisitVariableExpr(expr Variable) interface{}
	VisitAssignExpr(expr Assign) interface{}
	VisitGetExpr(expr Get) interface{}
	VisitSetExpr(expr Set) interface{}
}

type Expr interface {
	Accept(visitor ExpressionVisitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value interface{}
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Unary struct {
	Operator Token
	Right    Expr
}

type Variable struct {
	Name Token
}

type Assign struct {
	Name  Token
	Value Expr
}

type Get struct {
	Name   Token
	Object Expr
}

type Set struct {
	Name   Token
	Object Expr
	Value  Expr
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

func NewCall(callee Expr, paren Token, arguments []Expr) *Call {
	return &Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (c *Call) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitCallExpr(*c)
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

func NewLogical(left Expr, operator Token, right Expr) *Logical {
	return &Logical{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (l *Logical) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLogicalExpr(*l)
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

func NewAssign(name Token, value Expr) *Assign {
	return &Assign{
		Name:  name,
		Value: value,
	}
}

func (a *Assign) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitAssignExpr(*a)
}

func NewGet(name Token, object Expr) *Get {
	return &Get{
		Name:   name,
		Object: object,
	}
}

func (g *Get) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitGetExpr(*g)
}

func NewSet(object Expr, name Token, value Expr) *Set {
	return &Set{
		Name:   name,
		Object: object,
		Value:  value,
	}
}

func (s *Set) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitSetExpr(*s)
}

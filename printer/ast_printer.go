package printer

import (
	"fmt"
	"strings"

	"github.com/alxbckr/goloxv1/lox"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr lox.Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitBinaryExpr(expr lox.Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr lox.Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr lox.Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitUnaryExpr(expr lox.Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, parts ...interface{}) string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(name)
	for _, part := range parts {
		str.WriteString(" ")
		str.WriteString(part.(lox.Expr).Accept(a).(string))
	}
	str.WriteString(")")
	return str.String()
}

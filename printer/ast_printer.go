package printer

import (
	"fmt"
	"strings"

	"github.com/alxbckr/goloxv1/parser"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr parser.Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitBinaryExpr(expr parser.Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr parser.Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr parser.Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitUnaryExpr(expr parser.Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, parts ...interface{}) string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(name)
	for _, part := range parts {
		str.WriteString(" ")
		str.WriteString(part.(parser.Expr).Accept(a).(string))
	}
	str.WriteString(")")
	return str.String()
}

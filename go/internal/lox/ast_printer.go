package lox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(expr Expr) string {
	if expr == nil {
		return "nil"
	}
	return expr.Accept(ap).(string)
}

func (ap *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return ap.parenthesize("group", expr.Expression)
}

func (ap *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (ap *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(ap).(string))
	}

	builder.WriteString(")")
	return builder.String()
}

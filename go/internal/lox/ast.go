package lox

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func ExampleAst() string {
	ast := NewAstPrinter()
	expression := &Binary{
		Left: &Unary{
			Operator: *NewToken(MINUS, "-", nil, 1),
			Right:    &Literal{Value: 123},
		},
		Operator: *NewToken(STAR, "*", nil, 1),
		Right:    &Grouping{&Literal{Value: 45.67}},
	}

	return ast.Print(expression)
}

func (ap *AstPrinter) Print(expr Expr) string {
	return expr.Accept(ap).(string)
}

func (ap *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return ap.paranthesize(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return ap.paranthesize("group", expr.Expression)
}

func (ap *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (ap *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return ap.paranthesize(expr.Operator.lexeme, expr.Right)
}

func (ap *AstPrinter) paranthesize(name string, exprs ...Expr) string {
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

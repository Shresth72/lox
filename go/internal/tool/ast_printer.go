package tool

import (
	"fmt"
	"strings"

	"github.com/Shresth72/lox/internal/lox"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func ExampleAst() string {
	ast := NewAstPrinter()
	// expression := &lox.Binary{
	// 	Left: &lox.Unary{
	// 		Operator: *lox.NewToken(lox.MINUS, "-", nil, 1),
	// 		Right:    &lox.Literal{Value: 123},
	// 	},
	// 	Operator: *lox.NewToken(lox.STAR, "*", nil, 1),
	// 	Right:    &lox.Grouping{Expression: &lox.Literal{Value: 45.67}},
	// }

	expression := lox.NewBinary(
		lox.NewUnary(
			*lox.NewToken(lox.MINUS, "-", nil, 1),
			lox.NewLiteral(123),
		),
		*lox.NewToken(lox.STAR, "*", nil, 1),
		lox.NewGrouping(lox.NewLiteral(45.67)),
	)

	return ast.Print(expression)
}

func (ap *AstPrinter) Print(expr lox.Expr) string {
	return expr.Accept(ap).(string)
}

func (ap *AstPrinter) VisitBinaryExpr(expr *lox.Binary) interface{} {
	return ap.paranthesize(expr.Operator.GetLexeme(), expr.Left, expr.Right)
}

func (ap *AstPrinter) VisitGroupingExpr(expr *lox.Grouping) interface{} {
	return ap.paranthesize("group", expr.Expression)
}

func (ap *AstPrinter) VisitLiteralExpr(expr *lox.Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (ap *AstPrinter) VisitUnaryExpr(expr *lox.Unary) interface{} {
	return ap.paranthesize(expr.Operator.GetLexeme(), expr.Right)
}

func (ap *AstPrinter) paranthesize(name string, exprs ...lox.Expr) string {
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

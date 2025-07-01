package lox

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func (e *Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (e *Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(e)
}

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

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (e *Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		Expression: expression,
	}
}

func (e *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func NewLiteral(value any) *Literal {
	return &Literal{
		Value: value,
	}
}

func (e *Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (e *Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(e)
}

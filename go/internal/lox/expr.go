package lox

type Expr interface {
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

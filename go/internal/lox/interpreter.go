package lox

import "fmt"

type RuntimeError struct {
	Token   *Token
	Message string
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("Runtime error: %s\n[line %d]", r.Message, r.Token.Line)
}

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr Expr) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				panic(runtimeErr)
			} else {
				panic(r)
			}
		}
	}()

	value := i.evaluate(expr)
	return i.stringify(value), nil
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case MINUS:
		i.checkNumberOperands(&expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case PLUS:
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		panic(&RuntimeError{
			Token:   &expr.Operator,
			Message: "Operands must be two numbers or two strings.",
		})
	case SLASH:
		i.checkNumberOperands(&expr.Operator, left, right)
		rightNum := right.(float64)
		if rightNum == 0 {
			panic(&RuntimeError{
				Token:   &expr.Operator,
				Message: "Division by zero.",
			})
		}
		return left.(float64) / rightNum
	case STAR:
		i.checkNumberOperands(&expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		i.checkNumberOperands(&expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(&expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(&expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(&expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}
	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case MINUS:
		i.checkNumberOperand(&expr.Operator, right)
		return -right.(float64)
	case BANG:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) checkNumberOperand(operator *Token, operand interface{}) {
	if _, ok := operand.(float64); !ok {
		panic(&RuntimeError{
			Token:   operator,
			Message: "Operand must be a number.",
		})
	}
}

func (i *Interpreter) checkNumberOperands(operator *Token, left, right interface{}) {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if !leftOk || !rightOk {
		panic(&RuntimeError{
			Token:   operator,
			Message: "Operands must be numbers.",
		})
	}
}

func (i *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", value)
}

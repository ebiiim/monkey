package evaluator

import (
	"errors"
	"fmt"

	"github.com/ebiiim/monkey/ast"
	"github.com/ebiiim/monkey/object"
	"github.com/ebiiim/monkey/token"
)

// Global objects.
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Error values.
var (
	ErrTypeMismatch    = errors.New("type mismatch")
	ErrUnknownOperator = errors.New("unknown operator")
)

// Eval evaluates the program recursively.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpressions(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpressions(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt)
		// break if return or error
		switch result := obj.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return obj
}

func evalBlockStatements(stmts []ast.Statement) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt)
		if obj == nil {
			continue
		}
		// break if return or error
		objType := obj.Type()
		if objType == object.RETURN_VALUE_OBJ || objType == object.ERROR_OBJ {
			return obj
		}
	}
	return obj
}

func evalPrefixExpressions(op string, right object.Object) object.Object {
	switch op {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.MINUS:
		return evalMinusOperatorExpression(right)
	default:
		return newError(ErrUnknownOperator, "%s%s", op, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError(ErrUnknownOperator, "-%s", right.Type())
	}
	return &object.Integer{Value: -right.(*object.Integer).Value}
}

func evalInfixExpressions(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	// compare memory addresses because we have just one TRUE and FALSE
	case op == token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case op == token.NEQ:
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError(ErrTypeMismatch, "%s %s %s", left.Type(), op, right.Type())
	default:
		return newError(ErrUnknownOperator, "%s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	switch op {
	case token.PLUS:
		return &object.Integer{Value: l + r}
	case token.MINUS:
		return &object.Integer{Value: l - r}
	case token.ASTERISK:
		return &object.Integer{Value: l * r}
	case token.SLASH:
		return &object.Integer{Value: l / r}
	case token.LT:
		return nativeBoolToBooleanObject(l < r)
	case token.GT:
		return nativeBoolToBooleanObject(l > r)
	case token.EQ:
		return nativeBoolToBooleanObject(l == r)
	case token.NEQ:
		return nativeBoolToBooleanObject(l != r)
	default:
		return newError(ErrUnknownOperator, "%s %s %s", left.Type(), op, right.Type())
	}
}

func evalIfExpression(e *ast.IfExpression) object.Object {
	condition := Eval(e.Condition)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(e.Consequence)
	}
	if e.Alternative != nil {
		return Eval(e.Alternative)
	}
	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(v bool) object.Object {
	if v {
		return TRUE
	}
	return FALSE
}

func newError(errType error, format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	return &object.Error{Message: fmt.Errorf("%w: %s", errType, msg)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

package evaluator

import (
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

// Eval evaluates the program recursively.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpressions(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpressions(node.Operator, left, right)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt)
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
		return NULL
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
		return NULL
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
	default:
		return NULL
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
		return NULL
	}
}

func nativeBoolToBooleanObject(v bool) object.Object {
	if v {
		return TRUE
	}
	return FALSE
}

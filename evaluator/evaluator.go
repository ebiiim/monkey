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
	ErrTypeMismatch       = errors.New("type mismatch")
	ErrUnknownOperator    = errors.New("unknown operator")
	ErrIdentifierNotFound = errors.New("identifier not found")
	ErrIsNotFunction      = errors.New("not a function")
)

// Eval evaluates the program recursively.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	// expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpressions(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpressions(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Env: env, Parameters: params, Body: body}
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	}
	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt, env)
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

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var obj object.Object
	for _, stmt := range stmts {
		obj = Eval(stmt, env)
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	obj, ok := env.Get(node.Value)
	if !ok {
		return newError(ErrIdentifierNotFound, "%s", node.Value)
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

func evalIfExpression(e *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(e.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(e.Consequence, env)
	}
	if e.Alternative != nil {
		return Eval(e.Alternative, env)
	}
	return NULL
}

func evalExpressions(e []ast.Expression, env *object.Environment) []object.Object {
	var objs []object.Object
	for _, expr := range e {
		obj := Eval(expr, env)
		if isError(obj) {
			return []object.Object{obj}
		}
		objs = append(objs, obj)
	}
	return objs
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError(ErrIsNotFunction, "%s", fn.Type())
	}
	eEnv := extendFunctionEnv(function, args)
	ev := Eval(function.Body, eEnv)
	return unwrapReturnValue(ev)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	eEnv := object.NewEnclosedEnvironment(fn.Env)
	for i, paramName := range fn.Parameters {
		eEnv.Set(paramName.Value, args[i])
	}
	return eEnv
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue
	}
	return obj
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

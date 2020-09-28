package evaluator

import (
	"errors"

	"github.com/ebiiim/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {Fn: fnLen},
}

var (
	ErrTooManyArgs      = errors.New("too many arguments")
	ErrTooFewArgs       = errors.New("too few arguments")
	ErrTypeNotSupported = errors.New("type not supported")
)

var fnLen = func(args ...object.Object) object.Object {
	if len(args) == 0 {
		return newError(ErrTooFewArgs, "want=%d got=%d", 1, len(args))
	}
	if len(args) > 1 {
		return newError(ErrTooManyArgs, "want=%d got=%d", 1, len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError(ErrTypeNotSupported, "len(%T)", arg.Type())
	}
}

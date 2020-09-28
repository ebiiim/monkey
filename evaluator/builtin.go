package evaluator

import (
	"errors"
	"os"

	"github.com/ebiiim/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":  {Fn: fnLen},
	"exit": {Fn: fnExit},
}

var (
	ErrTooManyArgs      = errors.New("too many arguments")
	ErrTooFewArgs       = errors.New("too few arguments")
	ErrTypeNotSupported = errors.New("type not supported")
	ErrExitCode         = errors.New("exit code must be 0--125")
)

var fnLen = func(args ...object.Object) object.Object {
	if len(args) == 0 {
		return newError(ErrTooFewArgs, "want=%s got=%d", "1", len(args))
	}
	if len(args) > 1 {
		return newError(ErrTooManyArgs, "want=%s got=%d", "1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError(ErrTypeNotSupported, "len(%T)", arg.Type())
	}
}

// TODO: maybe this should return some signal instead of NULL (and do os.Exit()).
var fnExit = func(args ...object.Object) object.Object {
	if len(args) > 1 {
		return newError(ErrTooManyArgs, "want=%s got=%d", "0|1", len(args))
	}
	if len(args) == 0 {
		builtinExit(0)
		return NULL
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		if arg.Value < 0 || arg.Value > 125 {
			return newError(ErrExitCode, "exit(%d)", arg.Value)
		}
		builtinExit(int(arg.Value))
		return NULL
	default:
		return newError(ErrTypeNotSupported, "exit(%T)", arg.Type())
	}
}

const (
	BUILTIN_EXIT_RETURN_NULL = iota + 1
	BUILTIN_EXIT_OS_EXIT
)

var BUILTIN_EXIT = BUILTIN_EXIT_OS_EXIT

func builtinExit(code int) {
	if BUILTIN_EXIT == BUILTIN_EXIT_OS_EXIT {
		os.Exit(code)
	}
}

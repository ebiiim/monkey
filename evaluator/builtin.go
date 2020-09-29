package evaluator

import (
	"errors"

	"github.com/ebiiim/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: fnLen},
	"first": {Fn: fnFirst},
	"last":  {Fn: fnLast},
	"rest":  {Fn: fnRest},
	"push":  {Fn: fnPush},
	"pop":   {Fn: fnPop},
}

// Builtin function errors.
var (
	ErrTooManyArgs      = errors.New("too many arguments")
	ErrTooFewArgs       = errors.New("too few arguments")
	ErrTypeNotSupported = errors.New("type not supported")
	ErrArrayNeeded      = errors.New("argument must be Array")
	ErrFileOpenFailed   = errors.New("failed to open file")
)

var fnLen = func(args ...object.Object) object.Object {
	if errObj := hasNArgs(1, args...); errObj != nil {
		return errObj
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError(ErrTypeNotSupported, "len(%T)", arg.Type())
	}
}

var fnFirst = func(args ...object.Object) object.Object {
	if errObj := hasNArgs(1, args...); errObj != nil {
		return errObj
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError(ErrArrayNeeded, "first(%T)", args[0].Type())
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	return arr.Elements[0]
}

var fnLast = func(args ...object.Object) object.Object {
	if errObj := hasNArgs(1, args...); errObj != nil {
		return errObj
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError(ErrArrayNeeded, "last(%T)", args[0].Type())
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	return arr.Elements[len(arr.Elements)-1]
}

var fnRest = func(args ...object.Object) object.Object {
	if errObj := hasNArgs(1, args...); errObj != nil {
		return errObj
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError(ErrArrayNeeded, "rest(%T)", args[0].Type())
	}
	arr := args[0].(*object.Array)
	size := len(arr.Elements)
	if size == 0 {
		return NULL
	}
	newArr := make([]object.Object, size-1, size-1)
	copy(newArr, arr.Elements[1:size])
	return &object.Array{Elements: newArr}
}

var fnPush = func(args ...object.Object) object.Object {
	if errObj := hasNArgs(2, args...); errObj != nil {
		return errObj
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError(ErrArrayNeeded, "rest(%T)", args[0].Type())
	}
	arr := args[0].(*object.Array)
	size := len(arr.Elements)
	newArr := make([]object.Object, size+1, size+1)
	copy(newArr, arr.Elements)
	newArr[size] = args[1]
	return &object.Array{Elements: newArr}
}

var fnPop = func(args ...object.Object) object.Object {
	if errObj := hasNArgs(1, args...); errObj != nil {
		return errObj
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError(ErrArrayNeeded, "rest(%T)", args[0].Type())
	}
	arr := args[0].(*object.Array)
	size := len(arr.Elements)
	if size == 0 {
		return NULL
	}
	newArr := make([]object.Object, size-1, size-1)
	copy(newArr, arr.Elements[0:size-1])
	return &object.Array{Elements: newArr}
}

func hasNArgs(n int, args ...object.Object) object.Object {
	if len(args) == n {
		return nil
	} else if len(args) < n {
		return newError(ErrTooFewArgs, "want=%d got=%d", n, len(args))
	} else {
		return newError(ErrTooManyArgs, "want=%d got=%d", n, len(args))
	}
}

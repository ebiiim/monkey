package evaluator_test

import (
	"errors"
	"testing"

	"github.com/ebiiim/monkey/evaluator"
	"github.com/ebiiim/monkey/object"
)

func TestBuiltinFunctions(t *testing.T) {
	evaluator.BUILTIN_EXIT = evaluator.BUILTIN_EXIT_RETURN_NULL

	cases := []struct {
		input string
		want  interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, evaluator.ErrTypeNotSupported},
		{`len(1, 2)`, evaluator.ErrTooManyArgs},
		{`len()`, evaluator.ErrTooFewArgs},
		{`exit()`, evaluator.NULL},
		{`exit(0)`, evaluator.NULL},
		{`exit(1)`, evaluator.NULL},
		{`exit(125)`, evaluator.NULL},
		{`exit(-1)`, evaluator.ErrExitCode},
		{`exit(-126)`, evaluator.ErrExitCode},
		{`exit("1")`, evaluator.ErrTypeNotSupported},
		{`exit(1, 2, 3)`, evaluator.ErrTooManyArgs},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			switch want := c.want.(type) {
			case int:
				testIntegerObject(t, ev, int64(want))
			case object.Null:
				testNullObject(t, ev)
			case error:
				errObj, ok := ev.(*object.Error)
				if !ok {
					t.Errorf("no error object returned got=%T (%+v)", ev, ev)
					return
				}
				if !errors.Is(errObj.Message, want) {
					t.Errorf("wrong error type want=%+v got=%+v", want, errObj.Message)
					return
				}
			}
		})
	}
}

package evaluator_test

import (
	"errors"
	"testing"

	"github.com/ebiiim/monkey/evaluator"
	"github.com/ebiiim/monkey/object"
)

func TestBuiltinFunctions(t *testing.T) {
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
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			switch want := c.want.(type) {
			case int:
				testIntegerObject(t, ev, int64(want))
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

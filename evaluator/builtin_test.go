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
		{`len([])`, 0},
		{`len([1])`, 1},
		{`len([1, 2])`, 2},
		{`len(["hello world", 2 * 2 * 2 * 2, fn(x, y){ x * x + y * y }])`, 3},
		{`len(1)`, evaluator.ErrTypeNotSupported},
		{`len(1, 2)`, evaluator.ErrTooManyArgs},
		{`len()`, evaluator.ErrTooFewArgs},

		{`first([1])`, 1},
		{`first(["hello world", 2])`, "hello world"},
		{`first([1, 2, 3])`, 1},
		{`first([])`, evaluator.NULL},
		{`first()`, evaluator.ErrTooFewArgs},
		{`first([1], [1])`, evaluator.ErrTooManyArgs},
		{`first(1)`, evaluator.ErrArrayNeeded},

		{`last([1])`, 1},
		{`last(["hello world", 2])`, "2"},
		{`last([1, 2, 3])`, 3},
		{`last([])`, evaluator.NULL},
		{`last()`, evaluator.ErrTooFewArgs},
		{`last([1], [1])`, evaluator.ErrTooManyArgs},
		{`last(1)`, evaluator.ErrArrayNeeded},

		{`rest([1])`, "[]"},
		{`rest(["hello world", 2])`, "[2]"},
		{`rest([1, 2, 3])`, "[2, 3]"},
		{`rest([])`, evaluator.NULL},
		{`rest()`, evaluator.ErrTooFewArgs},
		{`rest([1], [1])`, evaluator.ErrTooManyArgs},
		{`rest(1)`, evaluator.ErrArrayNeeded},

		{`push([1], 2)`, "[1, 2]"},
		{`push(["hello world", 2], [3])`, `["hello world", 2, [3]]`},
		{`push([1, 2, 3], 4)`, "[1, 2, 3, 4]"},
		{`push([], 1)`, "[1]"},
		{`push()`, evaluator.ErrTooFewArgs},
		{`push([1])`, evaluator.ErrTooFewArgs},
		{`push([1], [1], [1])`, evaluator.ErrTooManyArgs},
		{`push(1, 1)`, evaluator.ErrArrayNeeded},

		{`pop([1])`, "[]"},
		{`pop(["hello world", 2])`, `["hello world"]`},
		{`pop([1, 2, 3])`, "[1, 2]"},
		{`pop([])`, evaluator.NULL},
		{`pop()`, evaluator.ErrTooFewArgs},
		{`pop([1], [1])`, evaluator.ErrTooManyArgs},
		{`pop(1)`, evaluator.ErrArrayNeeded},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			switch want := c.want.(type) {
			case int:
				testIntegerObject(t, ev, int64(want))
			case object.Null:
				testNullObject(t, ev)
			case *object.Array:
				testArrayObject(t, ev, want)
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

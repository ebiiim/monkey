package evaluator_test

import (
	"testing"

	"github.com/ebiiim/monkey/evaluator"
	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/object"
	"github.com/ebiiim/monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			testIntegerObject(t, ev, c.want)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			testBooleanObject(t, ev, c.want)
		})
	}
}

func testEval(input string) object.Object {
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	return evaluator.Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, want int64) bool {
	o, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer but %T (%+v)", obj, obj)
		return false
	}
	if o.Value != want {
		t.Errorf("object value want=%d got=%d", want, o.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, want bool) bool {
	o, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean but %T (%+v)", obj, obj)
		return false
	}
	if o.Value != want {
		t.Errorf("object value want=%v got=%v", want, o.Value)
		return false
	}
	return true

}

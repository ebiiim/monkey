package evaluator_test

import (
	"errors"
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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 5 + 5 + 5 + 5 - 10},
		{"2 * 2 * 2 * 2 * 2", 2 * 2 * 2 * 2 * 2},
		{"5 * 2 + 10", 5*2 + 10},
		{"5 + 10 * 10", 5 + 10*10},
		{"20 + 2 * -10", 20 + 2*-10},
		{"50 / 2 * 2 + 10", 50/2*2 + 10},
		{"2 * (5 + 10)", 2 * (5 + 10)},
		{"3 * 3 * 3 + 10", 3*3*3 + 10},
		{"3 * (3 * 3) + 10", 3*(3*3) + 10},
		{"(5 + 10 * 2 + 15 / 3) * 2 - 10", (5+10*2+15/3)*2 - 10},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			testBooleanObject(t, ev, c.want)
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	cases := []struct {
		input string
		want  interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10}, // truthy
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			integer, ok := c.want.(int)
			if ok {
				testIntegerObject(t, ev, int64(integer))
			} else {
				testNullObject(t, ev)
			}
		})
	}
}

func TestBangOperator(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			testBooleanObject(t, ev, c.want)
		})
	}
}

func TestReturnStatements(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; } return 1;", 10},
		{"if (10 > 1) { if (10 > 1) { return 1; 2; } return 3; }", 1},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			testIntegerObject(t, ev, c.want)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		input   string
		wantErr error
		wantMsg string
	}{
		{"5 + true;", evaluator.ErrTypeMismatch, "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", evaluator.ErrTypeMismatch, "type mismatch: INTEGER + BOOLEAN"},
		{"-true", evaluator.ErrUnknownOperator, "unknown operator: -BOOLEAN"},
		{"true + false;", evaluator.ErrUnknownOperator, "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", evaluator.ErrUnknownOperator, "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", evaluator.ErrUnknownOperator, "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false; 1; } return 1; }", evaluator.ErrUnknownOperator, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar;", evaluator.ErrIdentifierNotFound, "identifier not found: foobar"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ev := testEval(c.input)
			errObj, ok := ev.(*object.Error)
			if !ok {
				t.Errorf("no error object returned got=%T (%+v)", ev, ev)
				return
			}
			if !errors.Is(errObj.Message, c.wantErr) {
				t.Errorf("wrong error type want=%+v got=%+v", c.wantErr, errObj.Message)
				return
			}
			if errObj.Message.Error() != c.wantMsg {
				t.Errorf("wrong error message want=%s got=%s", c.wantMsg, errObj.Inspect())
				return
			}
		})
	}
}

func TestLetStatement(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			testIntegerObject(t, testEval(c.input), c.want)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	ev := testEval(input)

	fn, ok := ev.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function but %T (%+v)", fn, fn)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("len(fn.Parameters) want=1 got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("fn.Parameters[0] want=x got=%s", fn.Parameters[0])
	}
	wantBody := "(x + 2)"
	if fn.Body.String() != wantBody {
		t.Fatalf("fn.Body want=%s got=%s", fn.Body.String(), wantBody)
	}
}

func TestFunctionApplication(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x;}; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{"let newAdder = fn(x) { fn(y) { x + y }; }; let addTwo = newAdder(2); addTwo(10); ", 12},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			testIntegerObject(t, testEval(c.input), c.want)
		})
	}
}

func testEval(input string) object.Object {
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return evaluator.Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL but %T (%+v)", obj, obj)
		return false
	}
	return true
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

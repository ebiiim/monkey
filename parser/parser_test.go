package parser_test

import (
	"fmt"
	"testing"

	"github.com/ebiiim/monkey/ast"
	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/parser"
	"github.com/ebiiim/monkey/token"
)

func TestLetStatement(t *testing.T) {
	cases := []struct {
		input     string
		wantIdent string
		wantValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			p := parser.New(lexer.New(c.input))
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != 1 {
				t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
			}
			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, c.wantIdent) {
				return
			}
			val := stmt.(*ast.LetStatement).Value
			if !testLiteralExpression(t, val, c.wantValue) {
				return
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	cases := []struct {
		name            string
		input           string
		numStatements   int
		wantIdentifiers []string
	}{
		{"simple", `let x = 5;
let y = 10;
let foobar = 123456;
`,
			3, []string{
				"x",
				"y",
				"foobar",
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, ident := range c.wantIdentifiers {
				stmt := program.Statements[i]
				testLetStatement(t, stmt, ident)
			}
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, wantIdent string) bool {
	t.Helper()
	if s.TokenLiteral() != token.LET {
		t.Errorf("wrong s.TokenLiteral() want=%s got=%s", token.LET, s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement got=%T", s)
		return false
	}
	if letStmt.Name.Value != wantIdent {
		t.Errorf("wrong letStmt.Name.Value want=%s got=%s", wantIdent, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != wantIdent {
		t.Errorf("wrong letStmt.Name.TokenLiteral() want=%s got=%s", wantIdent, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func TestLetStatementErr(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		numStatements int
		numErrors     int
	}{
		{"simple", `let x = 5;
let = 10;
let foobar = 123456;
`,
			5, 2,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			if program == nil {
				t.Fatal("ParseProgram() returned nil")
			}
			if len(p.Errors()) != c.numErrors {
				t.Fatalf("p.Errors() has wrong length want=%d got=%d", c.numErrors, len(p.Errors()))
			}
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		numStatements int
	}{
		{"simple", `return 5;
return 10;
return 123456;
`,
			3,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for _, stmt := range program.Statements {
				testReturnStatement(t, stmt)
			}
		})
	}
}

func testReturnStatement(t *testing.T, s ast.Statement) {
	t.Helper()
	if s.TokenLiteral() != token.RETURN {
		t.Errorf("wrong s.TokenLiteral() want=%s got=%s", token.LET, s.TokenLiteral())
	}
	_, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("s is not *ast.ReturnStatement got=%T", s)
	}
}

func TestIdentifierExpression(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		numStatements int
		wantValues    []string
	}{
		{"#1", "foobar;", 1, []string{"foobar"}},
		{"#2", "foobar", 1, []string{"foobar"}},
		{"#3", "foobar;baz", 2, []string{"foobar", "baz"}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				if !testIdentifier(t, exprStmt.Expression, c.wantValues[i]) {
					return
				}
			}
		})
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		numStatements int
		wantValues    []int64
	}{
		{"#1", "5;", 1, []int64{5}},
		{"#2", "5", 1, []int64{5}},
		{"#3", "5;10", 2, []int64{5, 10}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				if !testIntegerLiteral(t, exprStmt.Expression, c.wantValues[i]) {
					return
				}
			}
		})
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		numStatements int
		operator      string
		integerValue  int64
	}{
		{"#1", "!5;", 1, "!", 5},
		{"#2", "-15;", 1, "-", 15},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				expr, ok := exprStmt.Expression.(*ast.PrefixExpression)
				if !ok {
					t.Fatalf("expr is not *ast.PrefixExpression but %T", exprStmt.Expression)
				}
				if expr.Operator != c.operator {
					t.Fatalf("expr.Operator want=%v got=%v", expr.Operator, c.operator)
				}
				if !testIntegerLiteral(t, expr.Right, c.integerValue) {
					return
				}
			}
		})
	}
}

func TestParsingInfixExpression(t *testing.T) {
	cases := []struct {
		name                  string
		input                 string
		numStatements         int
		operator              string
		leftValue, rightValue interface{}
	}{
		{"add", "5 + 5;", 1, "+", 5, 5},
		{"sub", "5 - 5;", 1, "-", 5, 5},
		{"mul", "5 * 5;", 1, "*", 5, 5},
		{"div", "5 / 5;", 1, "/", 5, 5},
		{"gt", "5 > 5;", 1, ">", 5, 5},
		{"lt", "5 < 5;", 1, "<", 5, 5},
		{"eq", "5 == 5;", 1, "==", 5, 5},
		{"neq", "5 != 5;", 1, "!=", 5, 5},
		{"eq_bool#1", "true == true;", 1, "==", true, true},
		{"eq_bool#2", "false == false;", 1, "==", false, false},
		{"neq_bool", "true != false;", 1, "!=", true, false},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				if !testInfixExpression(t, exprStmt.Expression, c.operator, c.leftValue, c.rightValue) {
					return
				}
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		// ident
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		// int
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		// bool
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		// grouped
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"(((((O - O)))))", "(O - O)"},
		// call
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			got := program.String()
			if got != c.want {
				t.Errorf("want=%s, got=%s", c.want, got)
			}
		})
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserError(t, p, program)
	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement but %T", program.Statements[0])
	}
	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr is not *ast.IfExpression but %T", expr)
	}
	if !testInfixExpression(t, expr.Condition, "<", "x", "y") {
		return
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Fatalf("len(expr.Consequence.Statements) want=1 got=%d", len(expr.Consequence.Statements))
	}
	conseq, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.Consequence.Statements[0] is not *ast.ExpressionStatement but %T", expr.Consequence.Statements[0])
	}
	if !testIdentifier(t, conseq.Expression, "x") {
		return
	}
	if expr.Alternative != nil {
		t.Fatalf("expr.Alternative want=nil, got=%+v", expr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserError(t, p, program)
	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement but %T", program.Statements[0])
	}
	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr is not *ast.IfExpression but %T", expr)
	}
	if !testInfixExpression(t, expr.Condition, "<", "x", "y") {
		return
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Fatalf("len(expr.Consequence.Statements) want=1 got=%d", len(expr.Consequence.Statements))
	}
	conseq, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.Consequence.Statements[0] is not *ast.ExpressionStatement but %T", expr.Consequence.Statements[0])
	}
	if !testIdentifier(t, conseq.Expression, "x") {
		return
	}
	if len(expr.Alternative.Statements) != 1 {
		t.Fatalf("len(expr.Alternative.Statements) want=1 got=%d", len(expr.Alternative.Statements))
	}
	alt, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.Alternative.Statements[0] is not *ast.ExpressionStatement but %T", expr.Alternative.Statements[0])
	}
	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserError(t, p, program)
	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement but %T", program.Statements[0])
	}
	expr, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionLiteral but %T", expr)
	}
	if len(expr.Parameters) != 2 {
		t.Fatalf("len(expr.Parameters) want=2 got=%d", len(expr.Parameters))
	}
	testOk := testLiteralExpression(t, expr.Parameters[0], "x")
	testOk = testOk && testLiteralExpression(t, expr.Parameters[1], "y")
	if len(expr.Body.Statements) != 1 {
		t.Fatalf("len(expr.Body.Statements) want=1 got=%d", len(expr.Body.Statements))
	}
	bodyStmt, ok := expr.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("bodyStmt is not *ast.ExpressionStatement but %T", expr)
	}
	testOk = testOk && testInfixExpression(t, bodyStmt.Expression, "+", "x", "y")
	if !testOk {
		t.Fatal("testOk is not true")
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	cases := []struct {
		input      string
		wantParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y) {};", []string{"x", "y"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			checkParserError(t, p, program)
			if len(program.Statements) != 1 {
				t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement but %T", program.Statements[0])
			}
			expr, ok := stmt.Expression.(*ast.FunctionLiteral)
			if !ok {
				t.Fatalf("expr is not *ast.FunctionLiteral but %T", expr)
			}
			if len(expr.Parameters) != len(c.wantParams) {
				t.Fatalf("len(expr.Parameters) want=%d got=%d", len(c.wantParams), len(expr.Parameters))
			}
			for i, p := range c.wantParams {
				testLiteralExpression(t, expr.Parameters[i], p)
			}
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserError(t, p, program)
	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement but %T", program.Statements[0])
	}
	expr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expr is not *ast.CallExpression but %T", expr)
	}
	testOk := testIdentifier(t, expr.Function, "add")
	if len(expr.Arguments) != 3 {
		t.Fatalf("len(expr.Arguments) want=3 got=%d", len(expr.Arguments))
	}
	testOk = testOk && testLiteralExpression(t, expr.Arguments[0], 1)
	testOk = testOk && testInfixExpression(t, expr.Arguments[1], "*", 2, 3)
	testOk = testOk && testInfixExpression(t, expr.Arguments[2], "+", 4, 5)
	if !testOk {
		t.Fatal("testOk is not true")
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserError(t, p, program)
	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) want=1 got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement but %T", program.Statements[0])
	}
	expr, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.StringLiteral but %T", expr)
	}
	if expr.Value != "hello world" {
		t.Errorf("expr.Value want=\"helloworld\" got=%s", expr.Value)
	}
}

// testInfixExpression tests if expr has an operator and two literals.
func testInfixExpression(t *testing.T, expr ast.Expression, op string, left, right interface{}) bool {
	t.Helper()
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not *ast.InfixExpression but %T", expr)
		return false
	}
	if opExpr.Operator != op {
		t.Errorf("opExpr.Operator want=%s got=%s", op, opExpr.Operator)
		return false
	}
	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}
	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, want interface{}) bool {
	t.Helper()
	switch v := want.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	default:
		t.Errorf("type of expr not handled: %T", expr)
		return false
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	t.Helper()
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr is not *ast.Identifier but %T", expr)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value want=%v got=%v", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() want=%v got=%v", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	t.Helper()
	intLit, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expr is not *ast.IntegerLiteral but %T", expr)
		return false
	}
	if intLit.Value != value {
		t.Errorf("intLit.Value want=%d got=%d", value, intLit.Value)
		return false
	}
	if intLit.TokenLiteral() != fmt.Sprint(value) {
		t.Errorf("intLit.TokenLiteral() want=%v got=%v", value, intLit.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	t.Helper()
	boolLit, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("expr is not *ast.BooleanLiteral but %T", expr)
		return false

	}
	if boolLit.Value != value {
		t.Errorf("boolLit.Value want=%v got=%v", value, boolLit.Value)
		return false
	}
	if boolLit.TokenLiteral() != fmt.Sprint(value) {
		t.Errorf("boolLit.TokenLiteral() want=%v got=%v", value, boolLit.TokenLiteral())
		return false
	}
	return true
}

func checkParserError(t *testing.T, parser *parser.Parser, program *ast.Program) {
	if program == nil {
		t.Fatal("program is nil")
	}
	expectNoError(t, parser.Errors())
}

func expectNoError(t *testing.T, errs []error) {
	t.Helper()
	if len(errs) != 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.Fatalf("got %d errors", len(errs))
	}
}

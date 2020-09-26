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
			if program == nil {
				t.Fatal("ParseProgram() returned nil")
			}
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
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

func testLetStatement(t *testing.T, s ast.Statement, wantIdent string) {
	t.Helper()
	if s.TokenLiteral() != token.LET {
		t.Errorf("wrong s.TokenLiteral() want=%s got=%s", token.LET, s.TokenLiteral())
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement got=%T", s)
	}
	if letStmt.Name.Value != wantIdent {
		t.Errorf("wrong letStmt.Name.Value want=%s got=%s", wantIdent, letStmt.Name.Value)
	}
	if letStmt.Name.TokenLiteral() != wantIdent {
		t.Errorf("wrong letStmt.Name.TokenLiteral() want=%s got=%s", wantIdent, letStmt.Name.TokenLiteral())
	}
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
			3, 1,
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

func TestReturnStatement(t *testing.T) {
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
			if program == nil {
				t.Fatal("ParseProgram() returned nil")
			}
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
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
			if program == nil {
				t.Fatal("ParseProgram() returned nil")
			}
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				ident, ok := exprStmt.Expression.(*ast.Identifier)
				if !ok {
					t.Fatalf("expr is not *ast.Identifier but %T", exprStmt.Expression)
				}
				if ident.Value != c.wantValues[i] {
					t.Fatalf("ident.Value want=%v got=%v", c.wantValues[i], ident.Value)
				}
				if ident.TokenLiteral() != c.wantValues[i] {
					t.Fatalf("ident.TokenLiteral() want=%v got=%v", c.wantValues[i], ident.TokenLiteral())
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
			if program == nil {
				t.Fatal("ParseProgram() returned nil")
			}
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				intLit, ok := exprStmt.Expression.(*ast.IntegerLiteral)
				if !ok {
					t.Fatalf("expr is not *ast.IntegerLiteral but %T", exprStmt.Expression)
				}
				if intLit.Value != c.wantValues[i] {
					t.Fatalf("intLit.Value want=%v got=%v", c.wantValues[i], intLit.Value)
				}
				if intLit.TokenLiteral() != fmt.Sprint(c.wantValues[i]) {
					t.Fatalf("intLit.TokenLiteral() want=%v got=%v", c.wantValues[i], intLit.TokenLiteral())
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
			if program == nil {
				t.Fatal("ParseProgram() returned nil")
			}
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
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
		leftValue, rightValue int64
	}{
		{"add", "5 + 5;", 1, "+", 5, 5},
		{"sub", "5 - 5;", 1, "-", 5, 5},
		{"mul", "5 * 5;", 1, "*", 5, 5},
		{"div", "5 / 5;", 1, "/", 5, 5},
		{"gt", "5 > 5;", 1, ">", 5, 5},
		{"lt", "5 < 5;", 1, "<", 5, 5},
		{"eq", "5 == 5;", 1, "==", 5, 5},
		{"neq", "5 != 5;", 1, "!=", 5, 5},
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
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
			for i, stmt := range program.Statements {
				exprStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("program.Statements[%d] is not *ast.ExpressionStatement but %T", i, program.Statements[0])
				}
				expr, ok := exprStmt.Expression.(*ast.InfixExpression)
				if !ok {
					t.Fatalf("expr is not *ast.InfixExpression but %T", exprStmt.Expression)
				}
				if expr.Operator != c.operator {
					t.Fatalf("expr.Operator want=%v got=%v", expr.Operator, c.operator)
				}
				if !testIntegerLiteral(t, expr.Left, c.leftValue) {
					return
				}
				if !testIntegerLiteral(t, expr.Right, c.rightValue) {
					return
				}
			}
		})
	}
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			l := lexer.New(c.input)
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("p.Errors: %v", p.Errors())
			}
			got := program.String()
			if got != c.want {
				t.Errorf("want=%s, got=%s", c.want, got)
			}
		})
	}
}

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


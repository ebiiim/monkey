package parser_test

import (
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
		t.Errorf("wrong s.TokenLiteral want=%s got=%s", token.LET, s.TokenLiteral())
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
				t.Fatalf("p.Errors() has wrong length want=%d got=%d", c.numErrors, (p.Errors()))
			}
			if len(program.Statements) != c.numStatements {
				t.Fatalf("program.Statements has wrong length want=%d got=%d", c.numStatements, len(program.Statements))
			}
		})
	}
}

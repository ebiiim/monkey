package ast_test

import (
	"testing"

	"github.com/ebiiim/monkey/ast"
	"github.com/ebiiim/monkey/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.New(token.LET, "let", 1, 1),
				Name: &ast.Identifier{
					Token: token.New(token.IDENT, "myVar", 1, 5),
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.New(token.IDENT, "anotherVar", 1, 13),
					Value: "anotherVar",
				},
			},
		},
	}
	wantCode := "let myVar = anotherVar;"

	if program.String() != wantCode {
		t.Errorf("program.String() want=%s got=%s", wantCode, program.String())
	}
}

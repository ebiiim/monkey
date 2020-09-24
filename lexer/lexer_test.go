package lexer_test

import (
	"testing"

	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/token"
)

func TestNextToken1(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		wantType         token.Type
		wantLiteral      string
		wantRow, wantCol int
	}{
		{token.ASSIGN, "=", 1, 1},
		{token.PLUS, "+", 1, 2},
		{token.LPAREN, "(", 1, 3},
		{token.RPAREN, ")", 1, 4},
		{token.LBRACE, "{", 1, 5},
		{token.RBRACE, "}", 1, 6},
		{token.COMMA, ",", 1, 7},
		{token.SEMICOLON, ";", 1, 8},
		{token.EOF, "", 1, 9},
	}

	l := lexer.New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.wantType {
			t.Fatalf("test %d: wrong token.Type want=%v got=%v", i, tt.wantType, tok.Type)
		}
		if tok.Literal != tt.wantLiteral {
			t.Fatalf("test %d: wrong token.Literal want=%v got=%v", i, tt.wantLiteral, tok.Literal)
		}
		if tok.Row != tt.wantRow {
			t.Fatalf("test %d: wrong token.Row want=%v got=%v", i, tt.wantRow, tok.Row)
		}
		if tok.Col != tt.wantCol {
			t.Fatalf("test %d: wrong token.Col want=%v got=%v", i, tt.wantCol, tok.Col)
		}
	}
}

func TestNextToken2(t *testing.T) {
	input := `let five = 5;
let ten = 10;
	
let add = fn(x, y) {
	x + y;
};
	
let result = add(five, ten);`

	tests := []struct {
		wantType         token.Type
		wantLiteral      string
		wantRow, wantCol int
	}{
		{token.LET, "let", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		{token.SEMICOLON, ";", 1, 13},

		{token.LET, "let", 2, 1},
		{token.IDENT, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "10", 2, 11},
		{token.SEMICOLON, ";", 2, 13},

		{token.LET, "let", 4, 1},
		{token.IDENT, "add", 4, 5},
		{token.ASSIGN, "=", 4, 9},
		{token.FUNCTION, "fn", 4, 11},
		{token.LPAREN, "(", 4, 13},
		{token.IDENT, "x", 4, 14},
		{token.COMMA, ",", 4, 15},
		{token.IDENT, "y", 4, 17},
		{token.RPAREN, ")", 4, 18},
		{token.LBRACE, "{", 4, 20},
		{token.IDENT, "x", 5, 5},
		{token.PLUS, "+", 5, 7},
		{token.IDENT, "y", 5, 9},
		{token.SEMICOLON, ";", 5, 10},
		{token.RBRACE, "}", 6, 1},
		{token.SEMICOLON, ";", 6, 2},

		{token.LET, "let", 8, 1},
		{token.IDENT, "result", 8, 5},
		{token.ASSIGN, "=", 8, 12},
		{token.IDENT, "add", 8, 14},
		{token.LPAREN, "(", 8, 17},
		{token.IDENT, "five", 8, 18},
		{token.COMMA, ",", 8, 22},
		{token.IDENT, "ten", 8, 24},
		{token.RPAREN, ")", 8, 27},
		{token.SEMICOLON, ";", 8, 28},

		{token.EOF, "", 8, 29},
	}

	l := lexer.New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.wantType {
			t.Fatalf("test %d: wrong token.Type want=%v got=%v", i, tt.wantType, tok.Type)
		}
		if tok.Literal != tt.wantLiteral {
			t.Fatalf("test %d: wrong token.Literal want=%v got=%v", i, tt.wantLiteral, tok.Literal)
		}
		if tok.Row != tt.wantRow {
			t.Fatalf("test %d: wrong token.Row want=%v got=%v", i, tt.wantRow, tok.Row)
		}
		if tok.Col != tt.wantCol {
			t.Fatalf("test %d: wrong token.Col want=%v got=%v", i, tt.wantCol, tok.Col)
		}
	}
}

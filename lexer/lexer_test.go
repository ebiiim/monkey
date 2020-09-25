package lexer_test

import (
	"testing"

	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/token"
)

func TestNextToken(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		wantTokens []token.Token
	}{
		{"1", `=+(){},;`, []token.Token{
			token.NewS(token.ASSIGN, "=", 1, 1),
			token.NewS(token.PLUS, "+", 1, 2),
			token.NewS(token.LPAREN, "(", 1, 3),
			token.NewS(token.RPAREN, ")", 1, 4),
			token.NewS(token.LBRACE, "{", 1, 5),
			token.NewS(token.RBRACE, "}", 1, 6),
			token.NewS(token.COMMA, ",", 1, 7),
			token.NewS(token.SEMICOLON, ";", 1, 8),
			token.NewS(token.EOF, "", 1, 9),
		}},
		{"2", `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);`, []token.Token{
			token.NewS(token.LET, "let", 1, 1),
			token.NewS(token.IDENT, "five", 1, 5),
			token.NewS(token.ASSIGN, "=", 1, 10),
			token.NewS(token.INT, "5", 1, 12),
			token.NewS(token.SEMICOLON, ";", 1, 13),

			token.NewS(token.LET, "let", 2, 1),
			token.NewS(token.IDENT, "ten", 2, 5),
			token.NewS(token.ASSIGN, "=", 2, 9),
			token.NewS(token.INT, "10", 2, 11),
			token.NewS(token.SEMICOLON, ";", 2, 13),

			token.NewS(token.LET, "let", 4, 1),
			token.NewS(token.IDENT, "add", 4, 5),
			token.NewS(token.ASSIGN, "=", 4, 9),
			token.NewS(token.FUNCTION, "fn", 4, 11),
			token.NewS(token.LPAREN, "(", 4, 13),
			token.NewS(token.IDENT, "x", 4, 14),
			token.NewS(token.COMMA, ",", 4, 15),
			token.NewS(token.IDENT, "y", 4, 17),
			token.NewS(token.RPAREN, ")", 4, 18),
			token.NewS(token.LBRACE, "{", 4, 20),
			token.NewS(token.IDENT, "x", 5, 5),
			token.NewS(token.PLUS, "+", 5, 7),
			token.NewS(token.IDENT, "y", 5, 9),
			token.NewS(token.SEMICOLON, ";", 5, 10),
			token.NewS(token.RBRACE, "}", 6, 1),
			token.NewS(token.SEMICOLON, ";", 6, 2),

			token.NewS(token.LET, "let", 8, 1),
			token.NewS(token.IDENT, "result", 8, 5),
			token.NewS(token.ASSIGN, "=", 8, 12),
			token.NewS(token.IDENT, "add", 8, 14),
			token.NewS(token.LPAREN, "(", 8, 17),
			token.NewS(token.IDENT, "five", 8, 18),
			token.NewS(token.COMMA, ",", 8, 22),
			token.NewS(token.IDENT, "ten", 8, 24),
			token.NewS(token.RPAREN, ")", 8, 27),
			token.NewS(token.SEMICOLON, ";", 8, 28),

			token.NewS(token.EOF, "", 8, 29),
		}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := lexer.New(c.input)
			for i, wt := range c.wantTokens {
				tok := l.NextToken()
				if tok.Type != wt.Type {
					t.Fatalf("token#%d: wrong token.Type want=%v got=%v", i, wt.Type, tok.Type)
				}
				if tok.Literal != wt.Literal {
					t.Fatalf("token#%d: wrong token.Literal want=%v got=%v", i, wt.Literal, tok.Literal)
				}
				if tok.Row != wt.Row {
					t.Fatalf("token#%d: wrong token.Row want=%v got=%v", i, wt.Row, tok.Row)
				}
				if tok.Col != wt.Col {
					t.Fatalf("token#%d: wrong token.Col want=%v got=%v", i, wt.Col, tok.Col)
				}
			}
		})
	}
}

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
		{"illegal", `@`, []token.Token{
			token.New(token.ILLEGAL, "@", 1, 1),
			token.New(token.EOF, "", 1, 2),
		}},
		{"simple", `=+(){},;`, []token.Token{
			token.New(token.ASSIGN, "=", 1, 1),
			token.New(token.PLUS, "+", 1, 2),
			token.New(token.LPAREN, "(", 1, 3),
			token.New(token.RPAREN, ")", 1, 4),
			token.New(token.LBRACE, "{", 1, 5),
			token.New(token.RBRACE, "}", 1, 6),
			token.New(token.COMMA, ",", 1, 7),
			token.New(token.SEMICOLON, ";", 1, 8),
			token.New(token.EOF, "", 1, 9),
		}},
		{"subset", `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);`, []token.Token{
			token.New(token.LET, "let", 1, 1),
			token.New(token.IDENT, "five", 1, 5),
			token.New(token.ASSIGN, "=", 1, 10),
			token.New(token.INT, "5", 1, 12),
			token.New(token.SEMICOLON, ";", 1, 13),

			token.New(token.LET, "let", 2, 1),
			token.New(token.IDENT, "ten", 2, 5),
			token.New(token.ASSIGN, "=", 2, 9),
			token.New(token.INT, "10", 2, 11),
			token.New(token.SEMICOLON, ";", 2, 13),

			token.New(token.LET, "let", 4, 1),
			token.New(token.IDENT, "add", 4, 5),
			token.New(token.ASSIGN, "=", 4, 9),
			token.New(token.FUNCTION, "fn", 4, 11),
			token.New(token.LPAREN, "(", 4, 13),
			token.New(token.IDENT, "x", 4, 14),
			token.New(token.COMMA, ",", 4, 15),
			token.New(token.IDENT, "y", 4, 17),
			token.New(token.RPAREN, ")", 4, 18),
			token.New(token.LBRACE, "{", 4, 20),
			token.New(token.IDENT, "x", 5, 5),
			token.New(token.PLUS, "+", 5, 7),
			token.New(token.IDENT, "y", 5, 9),
			token.New(token.SEMICOLON, ";", 5, 10),
			token.New(token.RBRACE, "}", 6, 1),
			token.New(token.SEMICOLON, ";", 6, 2),

			token.New(token.LET, "let", 8, 1),
			token.New(token.IDENT, "result", 8, 5),
			token.New(token.ASSIGN, "=", 8, 12),
			token.New(token.IDENT, "add", 8, 14),
			token.New(token.LPAREN, "(", 8, 17),
			token.New(token.IDENT, "five", 8, 18),
			token.New(token.COMMA, ",", 8, 22),
			token.New(token.IDENT, "ten", 8, 24),
			token.New(token.RPAREN, ")", 8, 27),
			token.New(token.SEMICOLON, ";", 8, 28),

			token.New(token.EOF, "", 8, 29),
		}},
		{"section1.4#1", `!-/*5;
5 < 10 > 5;`, []token.Token{
			token.New(token.BANG, "!", 1, 1),
			token.New(token.MINUS, "-", 1, 2),
			token.New(token.SLASH, "/", 1, 3),
			token.New(token.ASTERISK, "*", 1, 4),
			token.New(token.INT, "5", 1, 5),
			token.New(token.SEMICOLON, ";", 1, 6),
			token.New(token.INT, "5", 2, 1),
			token.New(token.LT, "<", 2, 3),
			token.New(token.INT, "10", 2, 5),
			token.New(token.GT, ">", 2, 8),
			token.New(token.INT, "5", 2, 10),
			token.New(token.SEMICOLON, ";", 2, 11),
			token.New(token.EOF, "", 2, 12),
		}},
		{"section1.4#2", `if (5 < 10) {
	return true;
} else {
	return false;
}`, []token.Token{
			token.New(token.IF, "if", 1, 1),
			token.New(token.LPAREN, "(", 1, 4),
			token.New(token.INT, "5", 1, 5),
			token.New(token.LT, "<", 1, 7),
			token.New(token.INT, "10", 1, 9),
			token.New(token.RPAREN, ")", 1, 11),
			token.New(token.LBRACE, "{", 1, 13),

			token.New(token.RETURN, "return", 2, 5),
			token.New(token.TRUE, "true", 2, 12),
			token.New(token.SEMICOLON, ";", 2, 16),

			token.New(token.RBRACE, "}", 3, 1),
			token.New(token.ELSE, "else", 3, 3),
			token.New(token.LBRACE, "{", 3, 8),

			token.New(token.RETURN, "return", 4, 5),
			token.New(token.FALSE, "false", 4, 12),
			token.New(token.SEMICOLON, ";", 4, 17),

			token.New(token.RBRACE, "}", 5, 1),
			token.New(token.EOF, "", 5, 2),
		}},
		{"section1.4#3", `10 == 10;
10 != 9;
!`, []token.Token{
			token.New(token.INT, "10", 1, 1),
			token.New(token.EQ, "==", 1, 4),
			token.New(token.INT, "10", 1, 7),
			token.New(token.SEMICOLON, ";", 1, 9),
			token.New(token.INT, "10", 2, 1),
			token.New(token.NEQ, "!=", 2, 4),
			token.New(token.INT, "9", 2, 7),
			token.New(token.SEMICOLON, ";", 2, 8),
			token.New(token.BANG, "!", 3, 1),
			token.New(token.EOF, "", 3, 2),
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

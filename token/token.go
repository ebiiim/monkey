package token

// Type represents token types.
type Type string

// Token contains a token.
type Token struct {
	Type     Type
	Literal  string
	Row, Col int
}

// Token types.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 123456

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ  = "=="
	NEQ = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "fn"
	LET      = "let"
	TRUE     = "true"
	FALSE    = "false"
	IF       = "if"
	ELSE     = "else"
	RETURN   = "return"
)

// NewC initializes a Token with a character.
func NewC(t Type, ch byte, row, col int) Token {
	return Token{t, string(ch), row, col}
}

// NewS initializes a Token with a string.
func NewS(t Type, s string, row, col int) Token {
	return Token{t, s, row, col}
}

var keywords = map[string]Type{
	FUNCTION: FUNCTION,
	LET:      LET,
	TRUE:     TRUE,
	FALSE:    FALSE,
	IF:       IF,
	ELSE:     ELSE,
	RETURN:   RETURN,
}

// LookupIdent finds type of an identifier.
func LookupIdent(s string) Type {
	if tok, ok := keywords[s]; ok {
		return tok
	}
	return IDENT
}

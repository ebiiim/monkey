package token

// Type represents token types.
type Type string

// Token contains a token.
type Token struct {
	Type    Type
	Literal string
}

// Token types.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 123456

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "fn"
	LET      = "let"
)

// NewC initializes a Token with a character.
func NewC(t Type, ch byte) Token {
	return Token{t, string(ch)}
}

// NewS initializes a Token with a string.
func NewS(t Type, s string) Token {
	return Token{t, s}
}

var keywords = map[string]Type{
	FUNCTION: FUNCTION,
	LET:      LET,
}

// LookupIdent finds type of an identifier.
func LookupIdent(s string) Type {
	if tok, ok := keywords[s]; ok {
		return tok
	}
	return IDENT
}

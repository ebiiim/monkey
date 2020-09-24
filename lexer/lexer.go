package lexer

import (
	"github.com/ebiiim/monkey/token"
)

var defaultTabSize = 4

// Lexer represents a lexer.
type Lexer struct {
	input                  string
	position, readPosition int
	ch                     byte
	row, col               int
	tabSize                int
}

// New initializes a lexer.
func New(input string) *Lexer {
	l := &Lexer{
		input:   input,
		tabSize: defaultTabSize,
		row:     1,
		col:     0,
	}
	l.readChar() // read the first charactor
	return l
}

// NextToken reads the next token.
func (l *Lexer) NextToken() token.Token {
	// skip non-token characters
	skip := true
	for skip {
		switch l.ch {
		case ' ', '\r':
			l.readChar()
		case '\t':
			l.col += l.tabSize - 1
			l.readChar()
		case '\n':
			l.col = 0
			l.row++
			l.readChar()
		default:
			skip = false
		}
	}

	var tok token.Token
	switch l.ch {
	case '=':
		tok = token.NewC(token.ASSIGN, l.ch, l.row, l.col)
	case '+':
		tok = token.NewC(token.PLUS, l.ch, l.row, l.col)
	case ',':
		tok = token.NewC(token.COMMA, l.ch, l.row, l.col)
	case ';':
		tok = token.NewC(token.SEMICOLON, l.ch, l.row, l.col)
	case '(':
		tok = token.NewC(token.LPAREN, l.ch, l.row, l.col)
	case ')':
		tok = token.NewC(token.RPAREN, l.ch, l.row, l.col)
	case '{':
		tok = token.NewC(token.LBRACE, l.ch, l.row, l.col)
	case '}':
		tok = token.NewC(token.RBRACE, l.ch, l.row, l.col)
	case 0:
		tok = tokenNewS(token.EOF, "", l.row, l.col)
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			t := token.LookupIdent(lit)
			return tokenNewS(t, lit, l.row, l.col)
		}
		if isDigit(l.ch) {
			return tokenNewS(token.INT, l.readNumber(), l.row, l.col)
		}
		tok = token.NewC(token.ILLEGAL, l.ch, l.row, l.col)
	}
	l.readChar()
	return tok
}

func tokenNewS(t token.Type, s string, row, col int) token.Token {
	return token.NewS(t, s, row, col-len(s))
}

func (l *Lexer) readChar() {
	l.col++
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

package ast

import (
	"bytes"
	"fmt"

	"github.com/ebiiim/monkey/token"
)

type Node interface {
	TokenLiteral() string
	fmt.Stringer
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		fmt.Fprint(&out, s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *LetStatement) statementNode() {}

func (s *LetStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%s %s = ", s.TokenLiteral(), s.Name.String())
	if s.Value != nil {
		fmt.Fprintf(&out, "%s", s.Value.String())
	}
	fmt.Fprintf(&out, ";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (s *ReturnStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%s ", s.TokenLiteral())
	if s.ReturnValue != nil {
		fmt.Fprint(&out, s.ReturnValue.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ExpressionStatement) statementNode() {}

func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (l *IntegerLiteral) TokenLiteral() string { return l.Token.Literal }

func (l *IntegerLiteral) expressionNode() {}

func (l *IntegerLiteral) String() string {
	return l.Token.Literal
}

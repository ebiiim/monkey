package ast

import (
	"bytes"
	"fmt"

	"github.com/ebiiim/monkey/token"
)

type Node interface {
	// TokenLiteral returns token.Literal
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

func (s *LetStatement) statementNode() {}

func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }

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

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }

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

func (s *ExpressionStatement) statementNode() {}

func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

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

func (e *Identifier) expressionNode() {}

func (e *Identifier) TokenLiteral() string { return e.Token.Literal }

func (e *Identifier) String() string {
	return e.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (e *IntegerLiteral) expressionNode() {}

func (e *IntegerLiteral) TokenLiteral() string { return e.Token.Literal }

func (e *IntegerLiteral) String() string { return e.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (e *BooleanLiteral) expressionNode() {}

func (e *BooleanLiteral) TokenLiteral() string { return e.Token.Literal }

func (e *BooleanLiteral) String() string { return e.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // - or !
	Operator string
	Right    Expression
}

func (e *PrefixExpression) expressionNode() {}

func (e *PrefixExpression) TokenLiteral() string { return e.Token.Literal }

func (e *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right.String())
}

type InfixExpression struct {
	Token       token.Token
	Operator    string
	Left, Right Expression
}

func (e *InfixExpression) expressionNode() {}

func (e *InfixExpression) TokenLiteral() string { return e.Token.Literal }

func (e *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator, e.Right.String())
}

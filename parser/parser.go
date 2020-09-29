package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ebiiim/monkey/ast"
	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/token"
)

// Errors
var (
	ErrTokenType      = errors.New("ErrTokenType")
	ErrInvalidLiteral = errors.New("ErrInvalidLiteral")
	ErrNoParseFunc    = errors.New("ErrNoParseFunc")
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Precedences
const (
	LOWEST     = iota + 1
	EQUALS     // ==
	LESSGRATER // < or >
	SUM        // +
	PRODUCT    // *
	PREFIX     // -X or !X
	CALL       // fn(X)
	INDEX      // array[index]
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGRATER,
	token.GT:       LESSGRATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l              *lexer.Lexer
	errs           []error
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.RPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// set the first token
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []error {
	return p.errs
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		msg := fmt.Errorf("%d:%d no prefix parse function for %s found (%w)", p.curToken.Row, p.curToken.Col, p.curToken.Type, ErrNoParseFunc)
		p.errs = append(p.errs, msg)
		return nil
	}

	leftExpr := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExpr
		}
		p.nextToken()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		err := fmt.Errorf("%d:%d could not parse \"%s\" as integer (%w)", p.curToken.Row, p.curToken.Col, p.curToken.Literal, ErrInvalidLiteral)
		p.errs = append(p.errs, err)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement()
	if !p.peekTokenIs(token.ELSE) {
		return expr
	}
	p.nextToken()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Alternative = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bs := &ast.BlockStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}
	return bs
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fn.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var idents []*ast.Identifier
	for !p.peekTokenIs(token.RPAREN) {
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		idents = append(idents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken() // skip RPAREN
	return idents
}

func (p *Parser) parseInfixExpression(leftExpr ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     leftExpr,
	}
	prec := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(prec)
	return expr
}

func (p *Parser) parseCallExpression(leftExpr ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:    p.curToken,
		Function: leftExpr,
	}
	expr.Arguments = p.parseExpressionList(token.RPAREN)
	return expr
}

func (p *Parser) parseIndexExpression(leftExpr ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  leftExpr,
	}
	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return expr
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken}
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	return arr
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var args []ast.Expression
	for !p.peekTokenIs(end) {
		p.nextToken()
		expr := p.parseExpression(LOWEST)
		if expr == nil {
			return nil
		}
		args = append(args, expr)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken() // skip end token
	return args
}

func (p *Parser) curTokenIs(t token.Type) bool {
	if p.curToken.Type == t {
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	if p.peekToken.Type == t {
		return true
	}
	return false
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	err := fmt.Errorf("%d:%d expected \"%s\" but got \"%s\" instead (%w)", p.peekToken.Row, p.peekToken.Col, t, p.peekToken.Type, ErrTokenType)
	p.errs = append(p.errs, err)
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

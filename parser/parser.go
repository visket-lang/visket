package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[token.TokenType]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	Errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}

	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, p.peekToken.Type)
	p.Errors = append(p.Errors, msg)

	return false
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Code = p.parseExpr(LOWEST)

	return program
}

func (p *Parser) parseExpr(precedence int) ast.Node {
	left := p.parsePrefixExpression()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		p.nextToken()
		left = p.parseInfixExpression(left)
	}

	return left
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

func (p *Parser) parsePrefixExpression() ast.Node {
	var left ast.Node

	switch p.curToken.Type {
	case token.MINUS:
		left = p.parseMinusPrefix()
	case token.INT:
		left = p.parseIntegerLiteral()
	case token.LPAREN:
		left = p.parseGroupedExpression()
	default:
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	return left
}

func (p *Parser) parseMinusPrefix() ast.Node {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     &ast.IntegerLiteral{Token: token.New(token.INT, "0"), Value: 0},
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpr(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpr(precedence)

	return expr
}

func (p *Parser) parseIntegerLiteral() ast.Node {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	n, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %s as integer", p.curToken.Literal)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	lit.Value = n
	return lit
}

func (p *Parser) parseGroupedExpression() ast.Node {
	p.nextToken()

	exp := p.parseExpr(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

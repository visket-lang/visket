package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/token"
	"strconv"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	left := p.parsePrefixExpression()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		p.nextToken()
		left = p.parseInfixExpression(left)
	}

	return left
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	switch p.curToken.Type {
	case token.SUB:
		return p.parseMinusPrefix()
	case token.INT:
		return p.parseIntegerLiteral()
	case token.LPAREN:
		return p.parseGroupedExpression()
	case token.IDENT:
		return p.parseIdentifier()
	}

	p.error(fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type))
	return nil
}

func (p *Parser) parseMinusPrefix() *ast.InfixExpression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     &ast.IntegerLiteral{Token: token.New(token.INT, "0"), Value: 0},
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	n, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.error(fmt.Sprintf("Could not parse %s as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = n
	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.curToken}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	op := p.curToken.Literal

	switch op {
	case "(":
		return p.parseCallExpression(left)
	}

	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: op,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseCallExpression(left ast.Expression) *ast.CallExpression {
	function, ok := left.(*ast.Identifier)
	if !ok {
		return nil
	}

	expr := &ast.CallExpression{Token: p.curToken}
	expr.Function = function
	expr.Parameters = p.parseCallParameters()

	return expr
}

func (p *Parser) parseCallParameters() []ast.Expression {
	var params []ast.Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	params = append(params, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

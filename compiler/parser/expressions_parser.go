package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/token"
	"strconv"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	left := p.parsePrefixExpression()

	// TODO rewrite
	for !p.peekTokenIs(token.SEMICOLON) && (isAssign(p.peekToken) || precedence < p.peekPrecedence()) {
		p.nextToken()
		left = p.parseInfixExpression(left)
	}

	return left
}

func isAssign(tok token.Token) bool {
	switch tok.Type {
	case
		token.ASSIGN, token.ADD_ASSIGN, token.SUB_ASSIGN,
		token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN,
		token.SHL_ASSIGN, token.SHR_ASSIGN:
		return true
	}

	return false
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

	p.error(fmt.Sprintf("%s | no prefix parse function for %s found", p.curToken.Pos, p.curToken.Type))
	return nil
}

func (p *Parser) parseMinusPrefix() *ast.InfixExpression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     &ast.IntegerLiteral{Token: token.New(token.INT, "0", p.curToken.Pos), Value: 0},
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
		p.error(fmt.Sprintf("%s | Could not parse %s as integer", p.curToken.Pos, p.curToken.Literal))
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
	case "[":
		return p.parseIndexExpression(left)
	case
		token.ASSIGN, token.ADD_ASSIGN, token.SUB_ASSIGN,
		token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN,
		token.SHL_ASSIGN, token.SHR_ASSIGN:
		return p.parseAssignExpression(left)
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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
		Index: nil,
	}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseAssignExpression(left ast.Expression) *ast.AssignExpression {
	stmt := &ast.AssignExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

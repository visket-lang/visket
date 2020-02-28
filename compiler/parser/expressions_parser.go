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
	case token.FLOAT:
		return p.parseFloatLiteral()
	case token.STRING:
		return p.parseStringLiteral()
	case token.LPAREN:
		return p.parseGroupedExpression()
	case token.IDENT:
		return p.parseIdentifier()
	case token.NEW:
		return p.parseNewExpression()
	}

	p.error(fmt.Sprintf("%s | no prefix parse function for %s found", p.curToken.Pos, p.curToken.Type))
	return nil
}

func (p *Parser) parseMinusPrefix() *ast.InfixExpression {
	expr := &ast.InfixExpression{
		Left:  &ast.IntegerLiteral{Value: 0},
		OpPos: p.curPos,
		Op:    p.curLiteral,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	lit := &ast.IntegerLiteral{Pos: p.curPos}

	n, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.error(fmt.Sprintf("%s | Could not parse %s as integer", p.curToken.Pos, p.curToken.Literal))
		return nil
	}

	lit.Value = n
	return lit
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	lit := &ast.FloatLiteral{Pos: p.curPos}

	n, err := strconv.ParseFloat(p.curToken.Literal, 32)
	if err != nil {
		p.error(fmt.Sprintf("%s | Could not parse %s as float", p.curToken.Pos, p.curToken.Literal))
		return nil
	}

	lit.Value = n
	return lit
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	lit := &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

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
	return &ast.Identifier{
		Pos:  p.curPos,
		Name: p.curLiteral,
	}
}

func (p *Parser) parseNewExpression() *ast.NewExpression {
	expr := &ast.NewExpression{New: p.curPos}
	p.nextToken()

	expr.Ident = p.parseIdentifier()

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	op := p.curToken.Literal

	switch op {
	case "(":
		return p.parseCallExpression(left)
	case "[":
		return p.parseIndexExpression(left)
	case ".":
		return p.parseLoadMemberExpression(left)
	case
		token.ASSIGN, token.ADD_ASSIGN, token.SUB_ASSIGN,
		token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN,
		token.SHL_ASSIGN, token.SHR_ASSIGN:
		return p.parseAssignExpression(left)
	}

	expr := &ast.InfixExpression{
		Left:  left,
		OpPos: p.curPos,
		Op:    p.curLiteral,
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

	expr := &ast.CallExpression{LParen: p.peekToken.Pos}
	expr.Function = function
	expr.Args = p.parseCallArguments()
	expr.RParen = p.curPos

	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
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
		Left:   left,
		LBrack: p.curPos,
	}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	exp.RBrack = p.curPos

	return exp
}

func (p *Parser) parseLoadMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.LoadMemberExpression{
		Left:   left,
		Period: p.curPos,
	}
	p.nextToken()
	exp.MemberIdent = p.parseIdentifier()

	return exp
}

func (p *Parser) parseAssignExpression(left ast.Expression) *ast.AssignExpression {
	stmt := &ast.AssignExpression{
		Left:  left,
		OpPos: p.curPos,
		Op:    p.curLiteral,
	}

	p.nextToken()
	right := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if stmt.Op == token.ASSIGN {
		stmt.Value = right
		return stmt
	}

	// "=" 以外は糖衣構文として実装する
	value := &ast.InfixExpression{
		Left:  left,
		OpPos: stmt.OpPos,
		Right: right,
	}

	switch stmt.Op {
	case token.ADD_ASSIGN:
		value.Op = token.ADD
	case token.SUB_ASSIGN:
		value.Op = token.SUB
	case token.MUL_ASSIGN:
		value.Op = token.MUL
	case token.QUO_ASSIGN:
		value.Op = token.QUO
	case token.REM_ASSIGN:
		value.Op = token.REM
	case token.SHL_ASSIGN:
		value.Op = token.SHL
	case token.SHR_ASSIGN:
		value.Op = token.SHR
	}

	stmt.Op = "="
	stmt.Value = value

	return stmt
}

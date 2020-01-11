package parser

import (
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/token"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.peekToken.Type == token.ASSIGN {
		return p.parseReassignStatement()
	}

	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Ident = p.parseIdentifier()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReassignStatement() *ast.ReassignStatement {
	stmt := &ast.ReassignStatement{Token: p.curToken}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	stmt.Ident = p.parseIdentifier()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	p.nextToken()
	stmt.Ident = p.parseIdentifier()

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	// 以下コード生成を簡単にするため

	if len(stmt.Body.Statements) == 0 {
		return stmt
	}

	// 関数の末尾は return を強制させる
	lastBodyStatement := stmt.Body.Statements[len(stmt.Body.Statements)-1]
	_, ok := lastBodyStatement.(*ast.ReturnStatement)
	if !ok {
		p.error("missing return at end of function")
		return nil
	}

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var params []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	params = append(params, p.parseIdentifier())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, p.parseIdentifier())
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if !p.peekTokenIs(token.ELSE) {
		return stmt
	}

	p.nextToken()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.Alternative = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}
	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Init = p.parseStatement()
	}
	p.expect(token.SEMICOLON)

	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
		p.nextToken()
	}
	p.expect(token.SEMICOLON)

	if !p.curTokenIs(token.LBRACE) {
		stmt.Post = p.parseStatement()
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

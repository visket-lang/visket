package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/token"
	"github.com/arata-nvm/Solitude/compiler/types"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.peekToken.Type {
	case
		token.ASSIGN, token.ADD_ASSIGN, token.SUB_ASSIGN,
		token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN,
		token.SHL_ASSIGN, token.SHR_ASSIGN:
		return p.parseAssignStatement()
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

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		pt := p.parseType()
		stmt.Type = pt
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return stmt
	}

	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: p.peekToken}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	stmt.Ident = p.parseIdentifier()

	p.nextToken()
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

	if p.curTokenIs(token.SEMICOLON) {
		return stmt
	}

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

	params, paramTypes := p.parseFunctionParameters()
	stmt.Parameters = params

	var retType types.ParserType = types.VOID

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		retType = p.parseType()
	}

	stmt.Type = types.NewFuncType(retType, paramTypes)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	// 以下コード生成を簡単にするため

	if stmt.Type.RetType == types.VOID {
		return stmt
	}

	// 関数の末尾は return を強制させる
	lastBodyStatement := stmt.Body.Statements[len(stmt.Body.Statements)-1]
	_, ok := lastBodyStatement.(*ast.ReturnStatement)
	if !ok {
		p.error(fmt.Sprintf("%s | missing return at end of function", p.curToken.Pos))
		return nil
	}

	return stmt
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, []types.ParserType) {
	var params []*ast.Identifier
	var paramTypes []types.ParserType

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params, paramTypes
	}

	p.nextToken()

	params = append(params, p.parseIdentifier())
	if !p.expectPeek(token.COLON) {
		return nil, nil
	}
	p.nextToken()
	paramTypes = append(paramTypes, p.parseType())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, p.parseIdentifier())
		p.expectPeek(token.COLON)
		p.nextToken()
		paramTypes = append(paramTypes, p.parseType())
	}

	if !p.expectPeek(token.RPAREN) {
		return nil, nil
	}

	return params, paramTypes
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

	if p.peekTokenIs(token.IN) {
		return p.parseForRangeStatement(stmt.Token)
	}

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

func (p *Parser) parseForRangeStatement(tok token.Token) *ast.ForStatement {
	stmt := &ast.ForStatement{Token: tok}

	ident := p.parseIdentifier()

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	start := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RANGE) {
		return nil
	}
	p.nextToken()

	end := p.parseExpression(LOWEST)

	// TODO すでに変数が宣言されている時の処理
	stmt.Init = &ast.VarStatement{
		Token: token.Token{Type: token.VAR, Literal: "var"},
		Ident: ident,
		Value: start,
	}

	stmt.Condition = &ast.InfixExpression{
		Token:    token.Token{Type: token.LTE, Literal: "<="},
		Left:     ident,
		Operator: "<=",
		Right:    end,
	}

	stmt.Post = &ast.AssignStatement{
		Token: token.Token{
			Type:    token.ADD_ASSIGN,
			Literal: "+=",
		},
		Ident: ident,
		Value: &ast.IntegerLiteral{
			Token: token.Token{
				Type:    token.INT,
				Literal: "1",
			},
			Value: 1,
		},
	}

	p.nextToken()
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

func (p *Parser) parseType() types.ParserType {
	if typ := types.ParseType(p.curToken.Literal); typ != nil {
		return typ
	}
	p.error(fmt.Sprintf("%s | unknown type %s", p.curToken.Pos, p.curToken.Literal))
	return nil
}

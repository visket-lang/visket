package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/token"
)

func (p *Parser) parseTopLevelStatement() ast.Statement {
	switch p.curToken.Type {
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.STRUCT:
		return p.parseStructStatement()
	}

	p.error(fmt.Sprintf("%s | func expected, got '%s'", p.curToken.Pos, p.curToken.Literal))
	return nil
}

func (p *Parser) parseStatement() ast.Statement {
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

func (p *Parser) parseStructStatement() *ast.StructStatement {
	stmt := &ast.StructStatement{
		Struct: p.curPos,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Ident = p.parseIdentifier()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.LBrace = p.curPos

	for p.peekTokenIs(token.IDENT) {
		m := &ast.MemberDecl{}

		if !p.expectPeek(token.IDENT) {
			return nil
		}
		m.Ident = p.parseIdentifier()

		if !p.expectPeek(token.IDENT) {
			return nil
		}
		m.Type = p.parseType()

		stmt.Members = append(stmt.Members, m)
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	stmt.RBrace = p.curPos

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Var: p.curPos}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Ident = p.parseIdentifier()

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		stmt.Type = p.parseType()
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return stmt
	}
	stmt.Assign = p.curPos

	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Return: p.curPos}

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
	stmt := &ast.FunctionStatement{
		Func: p.curPos,
		Sig:  &ast.FunctionSignature{},
	}

	p.nextToken()
	stmt.Ident = p.parseIdentifier()

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Sig.Params = p.parseFunctionParameters()

	retType := &ast.Type{
		Name: "void",
	}

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		retType = p.parseType()
	}

	stmt.Sig.RetType = retType

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Param {
	var params []*ast.Param

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	ident := p.parseIdentifier()
	if !p.expectPeek(token.COLON) {
		return nil
	}
	p.nextToken()
	typ := p.parseType()

	params = append(params, &ast.Param{
		Ident: ident,
		Type:  typ,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident = p.parseIdentifier()

		p.expectPeek(token.COLON)
		p.nextToken()
		typ = p.parseType()

		params = append(params, &ast.Param{
			Ident: ident,
			Type:  typ,
		})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{If: p.curPos}

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
	stmt := &ast.WhileStatement{While: p.curPos}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{For: p.curPos}
	p.nextToken()

	if p.peekTokenIs(token.IN) {
		return p.parseForRangeStatement(stmt.For)
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

func (p *Parser) parseForRangeStatement(tok token.Position) *ast.ForStatement {
	stmt := &ast.ForStatement{For: tok}

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
		Ident: ident,
		Type:  nil,
		Value: start,
	}

	stmt.Condition = &ast.InfixExpression{
		Left:  ident,
		Op:    "<=",
		Right: end,
	}

	stmt.Post = &ast.ExpressionStatement{
		Expression: &ast.AssignExpression{
			Left: ident,
			Op:   "=",
			Value: &ast.InfixExpression{
				Left: ident,
				Op:   "+",
				Right: &ast.IntegerLiteral{
					Value: 1,
				},
			},
		},
	}

	p.nextToken()
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{LBrace: p.curPos}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	block.RBrace = p.curPos

	return block
}

func (p *Parser) parseType() *ast.Type {
	typ := &ast.Type{}

	if p.curTokenIs(token.LBRACKET) {
		// 配列
		typ.IsArray = true
		p.nextToken()
		length := p.parseIntegerLiteral().Value
		typ.Len = uint64(length)
		p.expectPeek(token.RBRACKET)
		p.nextToken()
	}

	typ.Name = p.curLiteral
	typ.NamePos = p.curPos
	return typ
}

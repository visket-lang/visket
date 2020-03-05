package parser

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler/ast"
	"github.com/arata-nvm/visket/compiler/token"
)

func (p *Parser) parseTopLevelStatement() ast.Statement {
	switch p.curToken.Type {
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.STRUCT:
		return p.parseStructStatement()
	case token.VAR:
		return p.parseVarStatement()
	case token.IMPORT:
		return p.parseImportStatement()
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
		return p.parseFor()
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

		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		m.Type = p.parseType()

		stmt.Members = append(stmt.Members, m)
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	stmt.RBrace = p.curPos

	return stmt
}

// TODO rewrite
func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Import: p.curPos}

	if !p.peekTokenIs(token.STRING) {
		p.error(fmt.Sprintf("%s | expected next token to be %s, got %s instead", p.curToken.Pos, token.STRING, p.peekToken.Type))
		return nil
	}

	stmt.File = &ast.Identifier{
		Pos:  p.peekToken.Pos,
		Name: p.peekToken.Literal,
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Var: p.curPos}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Ident = p.parseIdentifier()

	if !p.peekTokenIs(token.COLON) && !p.peekTokenIs(token.ASSIGN) {
		p.error(fmt.Sprintf("%s | expected next token to be : or =, got %s instead", p.curPos, p.peekToken.Literal))
		return nil
	}

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

	isReference := p.peekTokenIs(token.REF)
	if isReference {
		p.nextToken()
	}
	p.nextToken()

	ident := p.parseIdentifier()
	if !p.expectPeek(token.COLON) {
		return nil
	}
	p.nextToken()
	typ := p.parseType()

	params = append(params, &ast.Param{
		Ident:       ident,
		Type:        typ,
		IsReference: isReference,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		isReference = p.peekTokenIs(token.REF)
		if isReference {
			p.nextToken()
		}
		p.nextToken()
		ident = p.parseIdentifier()

		p.expectPeek(token.COLON)
		p.nextToken()
		typ = p.parseType()

		params = append(params, &ast.Param{
			Ident:       ident,
			Type:        typ,
			IsReference: isReference,
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
	if p.peekTokenIs(token.IF) {
		p.nextToken()
		stmt.Alternative = &ast.BlockStatement{
			Statements: []ast.Statement{p.parseIfStatement()},
		}
		return stmt
	}

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

func (p *Parser) parseFor() ast.Statement {
	pos := p.curPos
	p.nextToken()

	var stmt ast.Statement
	if p.peekTokenIs(token.IN) {
		stmt = p.parseForRangeStatement(pos)
	} else {
		stmt = p.parseForStatement(pos)
	}

	return stmt
}

func (p *Parser) parseForStatement(pos token.Position) *ast.ForStatement {
	stmt := &ast.ForStatement{For: pos}

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

func (p *Parser) parseForRangeStatement(pos token.Position) *ast.ForRangeStatement {
	stmt := &ast.ForRangeStatement{For: pos}

	stmt.VarName = p.parseIdentifier()

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()
	stmt.In = p.curPos

	stmt.From = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RANGE) {
		return nil
	}
	p.nextToken()

	stmt.To = p.parseExpression(LOWEST)

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

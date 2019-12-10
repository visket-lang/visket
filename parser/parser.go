package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/token"
	"strconv"
)

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

func (p *Parser) ParseProgram() ast.Program {
	program := ast.Program{}

	program.Code = p.parseExpr()

	return program
}

func (p *Parser) parseExpr() ast.Node {
	left := p.parseIntegerLiteral()

	for !p.peekTokenIs(token.EOF) {
		p.nextToken()

		left = p.parseInfixExpression(left)
	}

	return left
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

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpr()

	return expr
}

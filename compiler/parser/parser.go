package parser

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/arata-nvm/Solitude/compiler/lexer"
	"github.com/arata-nvm/Solitude/compiler/token"
)

const (
	_ int = iota
	LOWEST
	RELATIONAL
	SHIFT
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       RELATIONAL,
	token.NEQ:      RELATIONAL,
	token.LT:       RELATIONAL,
	token.LTE:      RELATIONAL,
	token.GT:       RELATIONAL,
	token.GTE:      RELATIONAL,
	token.SHL:      SHIFT,
	token.SHR:      SHIFT,
	token.ADD:      SUM,
	token.SUB:      SUM,
	token.MUL:      PRODUCT,
	token.QUO:      PRODUCT,
	token.REM:      PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	Errors errors.ErrorList
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseTopLevelStatement()
		switch stmt := stmt.(type) {
		case *ast.FunctionStatement:
			program.Functions = append(program.Functions, stmt)
		default:
			p.error(fmt.Sprintf("unexpected statement: %s", stmt.Inspect()))
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// コメントはASTに含めない
	if p.curTokenIs(token.COMMENT) {
		p.nextToken()
	}
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) expect(tokenType token.TokenType) bool {
	if p.curTokenIs(tokenType) {
		p.nextToken()
		return true
	}

	p.error(fmt.Sprintf("%s | expected current token is %s, got %s instead", p.curToken.Pos, tokenType, p.peekToken.Type))
	return false
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}

	p.error(fmt.Sprintf("%s | expected next token to be %s, got %s instead", p.curToken.Pos, tokenType, p.peekToken.Type))
	return false
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

func (p *Parser) error(msg string) {
	p.Errors = append(p.Errors, msg)
}

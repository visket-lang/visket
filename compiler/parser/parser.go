package parser

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler/ast"
	"github.com/arata-nvm/visket/compiler/errors"
	"github.com/arata-nvm/visket/compiler/lexer"
	"github.com/arata-nvm/visket/compiler/token"
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
	token.PERIOD:   INDEX,
	token.MODSEP:   INDEX,
}

type Parser struct {
	l          []*lexer.Lexer
	curToken   token.Token
	curPos     token.Position
	curLiteral string
	peekToken  token.Token

	Errors errors.ErrorList
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: []*lexer.Lexer{l},
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
		case *ast.ModuleStatement:
			program.Modules = append(program.Modules, stmt)
		case *ast.FunctionStatement:
			program.Functions = append(program.Functions, stmt)
		case *ast.StructStatement:
			program.Structs = append(program.Structs, stmt)
		case *ast.VarStatement:
			program.Globals = append(program.Globals, stmt)
		case *ast.ImportStatement:
			if ok := p.importFile(stmt.File.Name); !ok {
				errors.ErrorExit(fmt.Sprintf("%s | cannot import '%s'", stmt.Import, stmt.File.Name))
			}
		default:
			p.error(fmt.Sprintf("unexpected statement: %s", ast.Show(stmt)))
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.curPos = p.curToken.Pos
	p.curLiteral = p.curToken.Literal
	p.peekToken = p.l[len(p.l)-1].NextToken()

	// コメントはASTに含めない
	if p.curTokenIs(token.COMMENT) {
		p.nextToken()
	}

	if p.curTokenIs(token.EOF) && len(p.l) > 1 {
		p.l = p.l[:len(p.l)-1]
		p.nextToken()
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

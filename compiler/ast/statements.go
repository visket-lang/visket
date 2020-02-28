package ast

import (
	"github.com/arata-nvm/Solitude/compiler/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

type FunctionStatement struct {
	Token token.Token
	Sig   *FunctionSignature
	Body  *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}

type FunctionSignature struct {
	Ident   *Identifier
	Params  []*Param
	RetType *Type
}

type Param struct {
	Ident *Identifier
	Type  *Type
}

type VarStatement struct {
	Token token.Token
	Ident *Identifier
	Value Expression
	Type  *Type
}

func (vs *VarStatement) statementNode() {}

type Type struct {
	Token token.Token

	IsArray bool
	Len     uint64
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode() {}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}

type ForStatement struct {
	Token     token.Token
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode() {}

type StructStatement struct {
	Token   token.Token
	Ident   *Identifier
	Members []*MemberDecl
}

func (ss *StructStatement) statementNode() {}

type MemberDecl struct {
	Ident *Identifier
	Type  *Type
}

package ast

import (
	"github.com/arata-nvm/visket/compiler/token"
)

type BlockStatement struct {
	LBrace     token.Position
	Statements []Statement
	RBrace     token.Position
}

func (bs *BlockStatement) statementNode() {}

type ModuleStatement struct {
	Module    token.Position
	Ident     *Identifier
	LBrace    token.Position
	Functions []*FunctionStatement
	RBrace    token.Position
}

func (ms *ModuleStatement) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

type FunctionStatement struct {
	Func  token.Position
	Ident *Identifier
	Sig   *FunctionSignature
	Body  *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}

type FunctionSignature struct {
	Params  []*Param
	RetType *Type
}

type Param struct {
	Ident       *Identifier
	Type        *Type
	IsReference bool
}

type VarStatement struct {
	Var    token.Position
	Ident  *Identifier
	Type   *Type
	Assign token.Position
	Value  Expression

	IsConstant bool
}

func (vs *VarStatement) statementNode() {}

type Type struct {
	NamePos token.Position
	Name    string

	IsArray bool
	Len     uint64
}

type ReturnStatement struct {
	Return token.Position
	Value  Expression
}

func (rs *ReturnStatement) statementNode() {}

type IfStatement struct {
	If          token.Position
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode() {}

type WhileStatement struct {
	While     token.Position
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}

type ForStatement struct {
	For       token.Position
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode() {}

type ForRangeStatement struct {
	For     token.Position
	VarName *Identifier
	In      token.Position
	From    Expression
	To      Expression
	Body    *BlockStatement
}

func (fs *ForRangeStatement) statementNode() {}

type StructStatement struct {
	Struct  token.Position
	Ident   *Identifier
	LBrace  token.Position
	Members []*MemberDecl
	RBrace  token.Position
}

func (ss *StructStatement) statementNode() {}

type MemberDecl struct {
	Ident *Identifier
	Type  *Type
}

type ImportStatement struct {
	Import token.Position
	File   *Identifier
}

func (is *ImportStatement) statementNode() {}

type IncludeStatement struct {
	Include token.Position
	File    *Identifier
}

func (is *IncludeStatement) statementNode() {}

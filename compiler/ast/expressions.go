package ast

import (
	"github.com/arata-nvm/Solitude/compiler/token"
)

type Identifier struct {
	Pos  token.Position
	Name string
}

func (i *Identifier) expressionNode() {}

type IntegerLiteral struct {
	Pos   token.Position
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

type FloatLiteral struct {
	Pos   token.Position
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

type PrefixExpression struct {
	OpPos token.Position
	Op    string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Left  Expression
	OpPos token.Position
	Op    string
	Right Expression
}

func (ie *InfixExpression) expressionNode() {}

type AssignExpression struct {
	Left  Expression
	OpPos token.Position
	Op    string
	Value Expression
}

func (rs *AssignExpression) expressionNode() {}

type CallExpression struct {
	Function *Identifier
	LParen   token.Position
	Args     []Expression
	RParen   token.Position
}

func (ce *CallExpression) expressionNode() {}

type IndexExpression struct {
	Left   Expression
	LBrack token.Position
	Index  Expression
	RBrack token.Position
}

func (ie *IndexExpression) expressionNode() {}

type NewExpression struct {
	New   token.Position
	Ident *Identifier
}

func (ne *NewExpression) expressionNode() {}

type LoadMemberExpression struct {
	Left        Expression
	Period      token.Position
	MemberIdent *Identifier
}

func (lme *LoadMemberExpression) expressionNode() {}

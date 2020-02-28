package ast

import (
	"github.com/arata-nvm/Solitude/compiler/token"
)

type Identifier struct {
	Token token.Token
}

func (i *Identifier) expressionNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

type AssignExpression struct {
	Token token.Token
	Left  Expression
	Value Expression
}

func (rs *AssignExpression) expressionNode() {}

type CallExpression struct {
	Token      token.Token
	Function   *Identifier
	Parameters []Expression
}

func (ce *CallExpression) expressionNode() {}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

type NewExpression struct {
	Token token.Token
	Ident *Identifier
}

func (ne *NewExpression) expressionNode() {}

type LoadMemberExpression struct {
	Token       token.Token
	Left        Expression
	MemberIdent *Identifier
}

func (lme *LoadMemberExpression) expressionNode() {}

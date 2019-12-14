package ast

import (
	"fmt"
	"github.com/arata-nvm/Solitude/token"
)

type Node interface {
	String() string
	Inspect() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Code Statement
}

func (p *Program) Inspect() string {
	return p.Code.Inspect()
}

func (p *Program) String() string {
	return p.Code.String()
}

type Identifier struct {
	Token token.Token
}

func (i *Identifier) Inspect() string {
	return fmt.Sprintf("Ident(%s)", i.Token.Literal)
}

func (i *Identifier) String() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) Inspect() string {
	return es.Expression.Inspect()
}

func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

func (es *ExpressionStatement) statementNode() {}

type VarStatement struct {
	Token token.Token
	Ident *Identifier
	Value Expression
}

func (vs *VarStatement) Inspect() string {
	return fmt.Sprintf("var %s = %s", vs.Ident.Inspect(), vs.Value.Inspect())
}

func (vs *VarStatement) String() string {
	return fmt.Sprintf("var %s = %s", vs.Ident.String(), vs.Value.String())
}

func (vs *VarStatement) statementNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) Inspect() string {
	return fmt.Sprintf("Int(%d)", il.Value)
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) expressionNode() {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Inspect() string {
	return fmt.Sprintf("Prefix(%s %s)", pe.Operator, pe.Right.Inspect())
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("%s %s", pe.Operator, pe.Right.String())
}

func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) Inspect() string {
	return fmt.Sprintf("Infix(%s %s %s)", ie.Left.Inspect(), ie.Operator, ie.Right.Inspect())
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("%s %s %s", ie.Left.String(), ie.Operator, ie.Right.String())
}

func (ie *InfixExpression) expressionNode() {}

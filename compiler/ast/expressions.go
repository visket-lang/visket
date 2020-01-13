package ast

import (
	"bytes"
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/token"
	"strings"
)

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

type CallExpression struct {
	Token      token.Token
	Function   *Identifier
	Parameters []Expression
}

func (ce *CallExpression) Inspect() string {
	var buf bytes.Buffer

	var p []string
	for _, param := range ce.Parameters {
		p = append(p, param.Inspect())
	}

	buf.WriteString(fmt.Sprintf("Call(%s(", ce.Function.Inspect()))
	buf.WriteString(fmt.Sprintf("%s", strings.Join(p, ",")))
	buf.WriteString("))")
	return buf.String()
}

func (ce *CallExpression) String() string {
	var buf bytes.Buffer

	var p []string
	for _, param := range ce.Parameters {
		p = append(p, param.String())
	}

	buf.WriteString(fmt.Sprintf("Call(%s(", ce.Function.String()))
	buf.WriteString(fmt.Sprintf("%s", strings.Join(p, ",")))
	buf.WriteString("))")
	return buf.String()
}

func (ce *CallExpression) expressionNode() {}
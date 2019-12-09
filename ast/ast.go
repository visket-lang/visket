package ast

import (
	"fmt"
	"github.com/arata-nvm/Solitude/token"
)

type Node interface {
	String() string
}

type Program struct {
	Code Node
}

func (p Program) String() string {
	return p.Code.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il IntegerLiteral) String() string {
	return il.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Left     Node
	Operator string
	Right    Node
}

func (ie InfixExpression) String() string {
	return fmt.Sprintf("%s %s %s", ie.Left.String(), ie.Operator, ie.Right.String())
}

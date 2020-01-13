package ast

import (
	"bytes"
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
	Statements []Statement
}

func (p *Program) Inspect() string {
	var buf bytes.Buffer
	for _, s := range p.Statements {
		buf.WriteString(s.Inspect())
	}
	return buf.String()
}

func (p *Program) String() string {
	var buf bytes.Buffer
	for _, s := range p.Statements {
		buf.WriteString(s.String())
	}
	return buf.String()
}
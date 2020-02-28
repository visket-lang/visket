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
	Functions []*FunctionStatement
	Structs   []*StructStatement
}

func (p *Program) Inspect() string {
	var buf bytes.Buffer
	for _, s := range p.Functions {
		buf.WriteString(s.Inspect())
	}

	for _, s := range p.Structs {
		buf.WriteString(s.Inspect())
	}
	return buf.String()
}

func (p *Program) String() string {
	var buf bytes.Buffer
	for _, s := range p.Functions {
		buf.WriteString(s.String())
	}

	for _, s := range p.Structs {
		buf.WriteString(s.String())
	}
	return buf.String()
}

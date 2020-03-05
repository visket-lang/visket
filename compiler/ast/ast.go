package ast

type Node interface {
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
	Globals   []*VarStatement
}

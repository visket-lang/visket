package ast

import (
	"bytes"
	"fmt"
	"github.com/arata-nvm/Solitude/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) Inspect() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for _, s := range bs.Statements {
		buf.WriteString(s.Inspect())
	}
	buf.WriteString("}")
	return buf.String()
}

func (bs *BlockStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for _, s := range bs.Statements {
		buf.WriteString("  ")
		buf.WriteString(s.String())
		buf.WriteString("\n")
	}
	buf.WriteString("}\n\n")
	return buf.String()
}

func (bs *BlockStatement) statementNode() {}

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

type FunctionStatement struct {
	Token     token.Token
	Ident     *Identifier
	Parameter *Identifier
	Body      *BlockStatement
}

func (fs *FunctionStatement) Inspect() string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	buf.WriteString(fs.Ident.Inspect())
	buf.WriteString("(")
	if fs.Parameter != nil {
		buf.WriteString(fs.Parameter.Inspect())
	}
	buf.WriteString(") ")
	buf.WriteString(fs.Body.Inspect())
	return buf.String()
}

func (fs *FunctionStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	buf.WriteString(fs.Ident.String())
	buf.WriteString("(")
	if fs.Parameter != nil {
		buf.WriteString(fs.Parameter.String())
	}
	buf.WriteString(")")
	buf.WriteString(fs.Body.String())
	return buf.String()
}

func (fs *FunctionStatement) statementNode() {}

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

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) Inspect() string {
	return fmt.Sprintf("return %s", rs.Value.Inspect())
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", rs.Value.String())
}

func (rs *ReturnStatement) statementNode() {}

package ast

import (
	"bytes"
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/token"
	"github.com/arata-nvm/Solitude/compiler/types"
	"strings"
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
	Token      token.Token
	Ident      *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	Type       *types.SlFunction
}

func (fs *FunctionStatement) Inspect() string {
	var buf bytes.Buffer

	var p []string
	for _, param := range fs.Parameters {
		p = append(p, param.Inspect())
	}

	buf.WriteString("func ")
	buf.WriteString(fs.Ident.Inspect())
	buf.WriteString("(")
	buf.WriteString(strings.Join(p, ","))
	buf.WriteString(") ")
	buf.WriteString(fs.Body.Inspect())
	return buf.String()
}

func (fs *FunctionStatement) String() string {
	var buf bytes.Buffer

	var p []string
	for _, param := range fs.Parameters {
		p = append(p, param.String())
	}

	buf.WriteString("func ")
	buf.WriteString(fs.Ident.String())
	buf.WriteString("(")
	buf.WriteString(strings.Join(p, ","))
	buf.WriteString(")")
	buf.WriteString(fs.Body.String())
	return buf.String()
}

func (fs *FunctionStatement) statementNode() {}

type VarStatement struct {
	Token token.Token
	Ident *Identifier
	Value Expression
	Type  types.SlType
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

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) Inspect() string {
	var buf bytes.Buffer
	buf.WriteString("if ")
	buf.WriteString(is.Condition.Inspect())
	buf.WriteString(" ")
	buf.WriteString(is.Consequence.Inspect())
	if is.Alternative != nil {
		buf.WriteString(" else ")
		buf.WriteString(is.Alternative.Inspect())
	}
	return buf.String()
}

func (is *IfStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("if ")
	buf.WriteString(is.Condition.String())
	buf.WriteString(" ")
	buf.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		buf.WriteString(" else ")
		buf.WriteString(is.Alternative.String())
	}
	return buf.String()
}

func (is *IfStatement) statementNode() {}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) Inspect() string {
	var buf bytes.Buffer
	buf.WriteString("while ")
	buf.WriteString(ws.Condition.Inspect())
	buf.WriteString(" ")
	buf.WriteString(ws.Body.Inspect())
	return buf.String()
}

func (ws *WhileStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("while ")
	buf.WriteString(ws.Condition.String())
	buf.WriteString(" ")
	buf.WriteString(ws.Body.String())
	return buf.String()
}

func (ws *WhileStatement) statementNode() {}

type ForStatement struct {
	Token     token.Token
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fs *ForStatement) Inspect() string {
	var buf bytes.Buffer
	buf.WriteString("for ")
	if fs.Init != nil {
		buf.WriteString(fs.Init.Inspect())
	}
	buf.WriteString("; ")
	if fs.Condition != nil {
		buf.WriteString(fs.Condition.Inspect())
	}
	buf.WriteString("; ")
	if fs.Post != nil {
		buf.WriteString(fs.Post.Inspect())
	}
	buf.WriteString(" ")
	buf.WriteString(fs.Body.Inspect())
	return buf.String()
}

func (fs *ForStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("for ")
	if fs.Init != nil {
		buf.WriteString(fs.Init.String())
	}
	buf.WriteString("; ")
	if fs.Condition != nil {
		buf.WriteString(fs.Condition.String())
	}
	buf.WriteString("; ")
	if fs.Post != nil {
		buf.WriteString(fs.Post.String())
	}
	buf.WriteString(" ")
	buf.WriteString(fs.Body.String())
	return buf.String()
}

func (fs *ForStatement) statementNode() {}

package ast

import (
	"bytes"
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/token"
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
	Token token.Token
	Sig   *FunctionSignature
	Body  *BlockStatement
}

func (fs *FunctionStatement) Inspect() string {
	var buf bytes.Buffer

	buf.WriteString(fs.Sig.Inspect())
	buf.WriteString(" ")
	buf.WriteString(fs.Body.Inspect())
	return buf.String()
}

func (fs *FunctionStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(fs.Sig.String())
	buf.WriteString(" ")
	buf.WriteString(fs.Body.String())
	return buf.String()
}

func (fs *FunctionStatement) statementNode() {}

type FunctionSignature struct {
	Ident   *Identifier
	Params  []*Param
	RetType *Type
}

func (fs *FunctionSignature) Inspect() string {
	var buf bytes.Buffer

	for i, param := range fs.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.Inspect())
	}

	return fmt.Sprintf("func %s(%s): %s", fs.Ident.Inspect(), buf.String(), fs.RetType.Inspect())
}

func (fs *FunctionSignature) String() string {
	var buf bytes.Buffer

	for i, param := range fs.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}

	return fmt.Sprintf("func %s(%s): %s", fs.Ident.String(), buf.String(), fs.RetType.String())
}

type Param struct {
	Ident *Identifier
	Type  *Type
}

func (p *Param) Inspect() string {
	return fmt.Sprintf("%s: %s", p.Ident.Inspect(), p.Type.Inspect())
}

func (p *Param) String() string {
	return fmt.Sprintf("%s: %s", p.Ident.String(), p.Type.String())
}

type VarStatement struct {
	Token token.Token
	Ident *Identifier
	Value Expression
	Type  *Type
}

func (vs *VarStatement) Inspect() string {
	return fmt.Sprintf("var %s = %s", vs.Ident.Inspect(), vs.Value.Inspect())
}

func (vs *VarStatement) String() string {
	return fmt.Sprintf("var %s = %s", vs.Ident.String(), vs.Value.String())
}

func (vs *VarStatement) statementNode() {}

type Type struct {
	Token token.Token

	IsArray bool
	Len     uint64
}

func (t *Type) Inspect() string {
	return fmt.Sprintf("Type(%s)", t)
}

func (t *Type) String() string {
	var buf bytes.Buffer

	if t.IsArray {
		buf.WriteString(fmt.Sprintf("[%d]", t.Len))
	}

	buf.WriteString(t.Token.Literal)
	return buf.String()
}

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

type StructStatement struct {
	Token   token.Token
	Ident   *Identifier
	Members []*MemberDecl
}

func (ss *StructStatement) Inspect() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("struct %s { ", ss.Ident.Inspect()))
	for _, m := range ss.Members {
		buf.WriteString(fmt.Sprintf("%s %s ", m.Ident.Inspect(), m.Type.String()))
	}
	buf.WriteString("}")

	return buf.String()
}

func (ss *StructStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("struct %s { ", ss.Ident.String()))
	for _, m := range ss.Members {
		buf.WriteString(fmt.Sprintf("%s %s ", m.Ident.String(), m.Type.String()))
	}
	buf.WriteString("}")

	return buf.String()
}

func (ss *StructStatement) statementNode() {}

type MemberDecl struct {
	Ident *Identifier
	Type  *Type
}

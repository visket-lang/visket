package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/codegen/types"
)

type Value interface {
	Ident() string
	Operand() string
}

type Var struct {
	Type types.Types
	Op   int
}

func (v Var) Ident() string {
	return fmt.Sprintf("%%%d", v.Op)
}

func (v Var) Operand() string {
	return fmt.Sprintf("%s %%%d", v.Type, v.Op)
}

type Label string

type Variable struct {
	Ident *ast.Identifier
	Num   int
}

func (v *Variable) Next() {
	v.Num++
}

func (v *Variable) peekNext() *Variable {
	return &Variable{
		Ident: v.Ident,
		Num:   v.Num + 1,
	}
}

func (v *Variable) Operand() string {
	return fmt.Sprintf("%%%s.%d", v.Ident, v.Num)
}

type Context struct {
	variables map[string]*Variable
	parent    *Context
}

func newContext(parent *Context) *Context {
	return &Context{
		variables: make(map[string]*Variable),
		parent:    parent,
	}
}

func (c *Context) newVariable(ident *ast.Identifier) *Variable {
	v := &Variable{
		Ident: ident,
		Num:   -1,
	}

	c.variables[ident.String()] = v
	return v
}

func (c *Context) findVariable(ident *ast.Identifier) (*Variable, bool) {
	v, ok := c.variables[ident.String()]

	if !ok && c.parent != nil {
		return c.parent.findVariable(ident)
	}

	return v, ok
}

func (c *CodeGen) resetIndex() {
	c.index = -1
	c.labelIndex = -1
}

func (c *CodeGen) nextVar(types types.Types) Var {
	c.index++
	return Var{types, c.index}
}

func (c *CodeGen) nextLabel(name string) Label {
	c.labelIndex++
	return Label(fmt.Sprintf("%s.%d", name, c.labelIndex))
}

func (c *CodeGen) into() {
	c.context = newContext(c.context)
}

func (c *CodeGen) outOf() {
	if c.context.parent != nil {
		c.context = c.context.parent
	}
}

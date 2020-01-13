package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen/types"
)

type Value interface {
	TypeName() types.Types
	RegName() string
	Operand() string
}

type Register struct {
	Type types.Types
	Op   int
}

func (v Register) TypeName() types.Types {
	return v.Type
}

func (v Register) RegName() string {
	return fmt.Sprintf("%%%d", v.Op)
}

func (v Register) Operand() string {
	return fmt.Sprintf("%s %%%d", v.Type, v.Op)
}

type Label string

type Named struct {
	Type  types.Types
	Ident *ast.Identifier
	Num   int
}

func (v *Named) Next() {
	v.Num++
}

func (v *Named) peekNext() *Named {
	return &Named{
		Type:  v.Type,
		Ident: v.Ident,
		Num:   v.Num + 1,
	}
}

func (v *Named) TypeName() types.Types {
	return v.Type
}

func (v *Named) RegName() string {
	return fmt.Sprintf("%%%s.%d", v.Ident, v.Num)
}

func (v *Named) Operand() string {
	return fmt.Sprintf("%s %%%s.%d", v.Type, v.Ident, v.Num)
}

type Context struct {
	variables map[string]*Named
	parent    *Context
}

func newContext(parent *Context) *Context {
	return &Context{
		variables: make(map[string]*Named),
		parent:    parent,
	}
}

func (c *Context) newNamed(types types.Types, ident *ast.Identifier) *Named {
	v := &Named{
		Type:  types,
		Ident: ident,
		Num:   0,
	}

	c.variables[ident.String()] = v
	return v
}

func (c *Context) findVariable(ident *ast.Identifier) (*Named, bool) {
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

func (c *CodeGen) nextReg(types types.Types) Register {
	c.index++
	return Register{types, c.index}
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

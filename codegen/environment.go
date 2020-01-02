package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
)

type Value interface {
	Operand() string
}

type Pointer int

func (p Pointer) Operand() string {
	return fmt.Sprintf("%%%d", p)
}

type Object int

func (o Object) Operand() string {
	return fmt.Sprintf("%%%d", o)
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

func (c *CodeGen) resetIndex() {
	c.index = -1
	c.labelIndex = -1
}

func (c *CodeGen) nextPointer() Pointer {
	c.index++
	return Pointer(c.index)
}

func (c *CodeGen) nextValue() Object {
	c.index++
	return Object(c.index)
}

func (c *CodeGen) newVariable(ident *ast.Identifier) *Variable {
	v := &Variable{
		Ident: ident,
		Num:   -1,
	}

	c.variables[ident.String()] = v
	return v
}

func (c *CodeGen) findVariable(ident *ast.Identifier) (*Variable, bool) {
	v, ok := c.variables[ident.String()]
	return v, ok
}

func (c *CodeGen) nextLabel(name string) Label {
	c.labelIndex++
	return Label(fmt.Sprintf("%s.%d", name, c.labelIndex))
}

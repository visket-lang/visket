package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
)

type Pointer int
type Value int
type Label string

type Variable struct {
	Ident *ast.Identifier
	Num   int
}

func (v *Variable) String() string {
	return fmt.Sprintf("%s.%d", v.Ident, v.Num)
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

func (c *CodeGen) resetIndex() {
	c.index = -1
	c.labelIndex = -1
}

func (c *CodeGen) nextPointer() Pointer {
	c.index++
	return Pointer(c.index)
}

func (c *CodeGen) nextValue() Value {
	c.index++
	return Value(c.index)
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


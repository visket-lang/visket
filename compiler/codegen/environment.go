package codegen

import (
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen/internal"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type Context struct {
	variables map[string]Value
	functions map[string]*ir.Func
	parent    *Context
}

type Value struct {
	Value      value.Value
	IsVariable bool
}

func (v Value) Load(block *ir.Block) value.Value {
	if v.IsVariable {
		return block.NewLoad(internal.PtrElmType(v.Value), v.Value)
	}
	return v.Value
}

func newContext(parent *Context) *Context {
	return &Context{
		variables: make(map[string]Value),
		functions: make(map[string]*ir.Func),
		parent:    parent,
	}
}

func (c *Context) addVariable(ident *ast.Identifier, v Value) {
	c.variables[ident.String()] = v
}

func (c *Context) addVariableByName(name string, v Value) {
	c.variables[name] = v
}

func (c *Context) findVariable(ident *ast.Identifier) (Value, bool) {
	v, ok := c.variables[ident.String()]

	if !ok && c.parent != nil {
		return c.parent.findVariable(ident)
	}

	return v, ok
}

func (c *Context) addFunction(ident *ast.Identifier, f *ir.Func) {
	c.functions[ident.String()] = f
}

func (c *Context) addFunctionByName(name string, f *ir.Func) {
	c.functions[name] = f
}

func (c *Context) findFunction(ident *ast.Identifier) (*ir.Func, bool) {
	f, ok := c.functions[ident.String()]

	if !ok && c.parent != nil {
		return c.parent.findFunction(ident)
	}

	return f, ok
}

func (c *CodeGen) into() {
	c.context = newContext(c.context)
}

func (c *CodeGen) outOf() {
	if c.context.parent != nil {
		c.context = c.context.parent
	}
}

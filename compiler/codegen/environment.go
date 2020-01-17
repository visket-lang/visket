package codegen

import (
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type Context struct {
	variables map[string]value.Value
	functions map[string]*ir.Func
	parent    *Context
}

func newContext(parent *Context) *Context {
	return &Context{
		variables: make(map[string]value.Value),
		functions: make(map[string]*ir.Func),
		parent:    parent,
	}
}

func (c *Context) addVariable(ident *ast.Identifier, v value.Value) {
	c.variables[ident.String()] = v
}

func (c *Context) addVariableByName(name string, v value.Value) {
	c.variables[name] = v
}

func (c *Context) findVariable(ident *ast.Identifier) (value.Value, bool) {
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

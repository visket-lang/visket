package codegen

import (
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen/internal"
	"github.com/llir/llvm/ir"
	llvmType "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Context struct {
	variables map[string]Value
	functions map[string]*ir.Func
	types     map[string]llvmType.Type
	structs   map[string]*Struct
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
	c := &Context{
		variables: make(map[string]Value),
		functions: make(map[string]*ir.Func),
		types:     make(map[string]llvmType.Type),
		structs:   make(map[string]*Struct),
		parent:    parent,
	}

	c.initType()
	return c
}

func (c *Context) initType() {
	c.addType("void", llvmType.Void)
	c.addType("bool", llvmType.I1)
	c.addType("int", llvmType.I32)
	c.addType("float", llvmType.Float)
}

func (c *Context) addVariable(ident *ast.Identifier, v Value) {
	c.variables[ident.Token.Literal] = v
}

func (c *Context) addVariableByName(name string, v Value) {
	c.variables[name] = v
}

func (c *Context) findVariable(ident *ast.Identifier) (Value, bool) {
	v, ok := c.variables[ident.Token.Literal]

	if !ok && c.parent != nil {
		return c.parent.findVariable(ident)
	}

	return v, ok
}

func (c *Context) addFunction(ident *ast.Identifier, f *ir.Func) {
	c.functions[ident.Token.Literal] = f
}

func (c *Context) addFunctionByName(name string, f *ir.Func) {
	c.functions[name] = f
}

func (c *Context) findFunction(ident *ast.Identifier) (*ir.Func, bool) {
	f, ok := c.functions[ident.Token.Literal]

	if !ok && c.parent != nil {
		return c.parent.findFunction(ident)
	}

	return f, ok
}

func (c *Context) addType(name string, t llvmType.Type) {
	c.types[name] = t
}

func (c *Context) findType(name string) (llvmType.Type, bool) {
	t, ok := c.types[name]

	if !ok && c.parent != nil {
		return c.parent.findType(name)
	}

	return t, ok
}

func (c *Context) addStruct(name string, s *Struct) {
	c.structs[name] = s
	c.addType(name, s.Type)
}

func (c *Context) findStruct(name string) (*Struct, bool) {
	s, ok := c.structs[name]

	if !ok && c.parent != nil {
		return c.parent.findStruct(name)
	}

	return s, ok
}

func (c *CodeGen) into() {
	c.context = newContext(c.context)
}

func (c *CodeGen) outOf() {
	if c.context.parent != nil {
		c.context = c.context.parent
	}
}

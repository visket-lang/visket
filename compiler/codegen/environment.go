package codegen

import (
	"github.com/arata-nvm/visket/compiler/codegen/internal"
	"github.com/llir/llvm/ir"
	llvmType "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Context struct {
	variables map[string]Value
	functions map[string]*Func
	types     map[string]llvmType.Type
	structs   map[string]*Struct
	parent    *Context
}

type Value struct {
	Value       value.Value
	IsVariable  bool
	IsReference bool
	IsConstant  bool
}

func (v Value) Load(block *ir.Block) value.Value {
	if v.IsVariable {
		return block.NewLoad(internal.PtrElmType(v.Value), v.Value)
	}
	return v.Value
}

func (v Value) Dereference(block *ir.Block) Value {
	if v.IsReference {
		return Value{
			Value:       block.NewLoad(internal.PtrElmType(v.Value), v.Value),
			IsVariable:  true,
			IsReference: v.IsReference,
		}
	}
	return v
}

// TODO rewrite
type Func struct {
	Func        *ir.Func
	IsReference []bool
}

func newContext(parent *Context) *Context {
	c := &Context{
		variables: make(map[string]Value),
		functions: make(map[string]*Func),
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
	c.addType("int8", llvmType.I8)
	c.addType("int16", llvmType.I16)
	c.addType("int32", llvmType.I32)
	c.addType("int64", llvmType.I64)

	c.addType("float", llvmType.Float)
	c.addType("float32", llvmType.Float)
	c.addType("float64", llvmType.Double)

}

func (c *Context) addVariable(name string, v Value) {
	c.variables[name] = v
}

func (c *Context) findVariable(name string) (Value, bool) {
	v, ok := c.variables[name]

	if !ok && c.parent != nil {
		return c.parent.findVariable(name)
	}

	return v, ok
}

func (c *Context) findVariableCurrent(name string) (Value, bool) {
	v, ok := c.variables[name]
	return v, ok
}

func (c *Context) addFunction(name string, f *Func) {
	c.functions[name] = f
}

func (c *Context) findFunction(name string) (*Func, bool) {
	f, ok := c.functions[name]

	if !ok && c.parent != nil {
		return c.parent.findFunction(name)
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

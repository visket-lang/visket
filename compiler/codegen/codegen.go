package codegen

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler/ast"
	"github.com/arata-nvm/visket/compiler/errors"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"io"
)

type CodeGen struct {
	program *ast.Program
	output  io.Writer
	context *Context

	module *ir.Module

	initFunc *ir.Func
	mainFunc *ir.Func

	contextModuleName string
	contextFunction   *ir.Func
	contextEntryBlock *ir.Block
	contextBlock      *ir.Block
	contextCondAfter  []*ir.Block
}

func New(program *ast.Program, w io.Writer) *CodeGen {
	c := &CodeGen{
		program: program,
		output:  w,
		context: newContext(nil),
		module:  ir.NewModule(),
	}

	c.addGlobal()

	return c
}

func (c *CodeGen) addGlobal() {
	c.initFunc = c.module.NewFunc("global-init", types.Void)
	block := c.initFunc.NewBlock("entry")
	block.NewRet(nil)

	c.mainFunc = c.module.NewFunc("main", types.I32)
	c.context.addFunction(c.mainFunc.Name(), &Func{
		Func:        c.mainFunc,
		IsReference: []bool{},
	})
	block = c.mainFunc.NewBlock("entry")
	block.NewCall(c.initFunc)
	block.NewRet(constant.NewInt(types.I32, 0))
}

func (c *CodeGen) GenerateCode() {
	c.genStdlib()

	for _, s := range c.program.Structs {
		c.genStructStatement(s)
	}

	for _, s := range c.program.Functions {
		c.genFunctionDeclaration(s)
	}

	for _, s := range c.program.Modules {
		c.genModuleDeclaration(s)
	}

	for _, s := range c.program.Globals {
		c.genGlobalVarStatement(s)
	}

	for _, s := range c.program.Functions {
		c.genFunctionBody(s)
	}

	for _, s := range c.program.Modules {
		c.genModuleBody(s)
	}

	irCode := c.module.String()
	_, err := fmt.Fprint(c.output, irCode)
	if err != nil {
		errors.ErrorExit("failed writing ir code")
	}
}

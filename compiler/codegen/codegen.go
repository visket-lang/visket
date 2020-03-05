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
	isDebug bool
	context *Context

	module *ir.Module

	mainFunc *ir.Func

	contextFunction   *ir.Func
	contextEntryBlock *ir.Block
	contextBlock      *ir.Block
	contextCondAfter  []*ir.Block
}

func New(program *ast.Program, isDebug bool, w io.Writer) *CodeGen {
	c := &CodeGen{
		program: program,
		isDebug: isDebug,
		output:  w,
		context: newContext(nil),
		module:  ir.NewModule(),
	}

	c.addGlobal()

	return c
}

func (c *CodeGen) addGlobal() {
	c.mainFunc = c.module.NewFunc("main", types.I32)
	c.context.addFunction(c.mainFunc.Name(), c.mainFunc)
	block := c.mainFunc.NewBlock("entry")
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

	for _, s := range c.program.Functions {
		c.genFunctionBody(s)
	}

	irCode := c.module.String()
	_, err := fmt.Fprint(c.output, irCode)
	if err != nil {
		errors.ErrorExit("failed writing ir code")
	}
}

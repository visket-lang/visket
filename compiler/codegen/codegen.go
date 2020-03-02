package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/llir/llvm/ir"
	"io"
)

type CodeGen struct {
	program *ast.Program
	output  io.Writer
	isDebug bool
	context *Context

	module            *ir.Module
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
	return c
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

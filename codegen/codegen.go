package codegen

import (
	"github.com/arata-nvm/Solitude/ast"
	"io"
)

type CodeGen struct {
	program      *ast.Program
	index        int
	labelIndex   int
	isDebug      bool
	isTerminated bool
	output       io.Writer
}

func New(program *ast.Program, isDebug bool, w io.Writer) *CodeGen {
	c := &CodeGen{
		program: program,
		isDebug: isDebug,
		output:  w,
	}

	c.resetIndex()
	return c
}

func (c *CodeGen) GenerateCode() {
	for _, s := range c.program.Statements {
		c.genStatement(s)
	}
}

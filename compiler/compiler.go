package compiler

import (
	"bytes"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/arata-nvm/Solitude/compiler/lexer"
	"github.com/arata-nvm/Solitude/compiler/optimizer"
	"github.com/arata-nvm/Solitude/compiler/parser"
	"log"
)

type Compiler struct {
	Program *ast.Program
	isDebug bool
}

func New(isDebug bool) *Compiler {
	return &Compiler{
		isDebug: isDebug,
	}
}

func (c *Compiler) Compile(filename string) errors.ErrorList {
	l, err := lexer.NewFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	p := parser.New(l)
	c.Program = p.ParseProgram()
	return p.Errors
}

func (c *Compiler) Optimize() {
	o := optimizer.New(c.Program)
	o.Optimize()
}

func (c *Compiler) GenIR() string {
	var b bytes.Buffer
	cg := codegen.New(c.Program, c.isDebug, &b)
	cg.GenerateCode()
	return b.String()
}

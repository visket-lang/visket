package compiler

import (
	"bytes"
	"github.com/arata-nvm/visket/compiler/ast"
	"github.com/arata-nvm/visket/compiler/codegen"
	"github.com/arata-nvm/visket/compiler/errors"
	"github.com/arata-nvm/visket/compiler/lexer"
	"github.com/arata-nvm/visket/compiler/optimizer"
	"github.com/arata-nvm/visket/compiler/parser"
	"log"
)

type Compiler struct {
	Filename string
	Program  *ast.Program
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(filename string) errors.ErrorList {
	c.Filename = filename

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
	cg := codegen.New(c.Program, &b)
	cg.GenerateCode()
	return b.String()
}

func (c *Compiler) IncludeFiles() []string {
	var filenames []string
	for _, s := range c.Program.Includes {
		filenames = append(filenames, s.File.Name)
	}
	return filenames
}

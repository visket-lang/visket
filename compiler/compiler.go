package compiler

import (
	"bytes"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen"
	"github.com/arata-nvm/Solitude/compiler/lexer"
	"github.com/arata-nvm/Solitude/compiler/optimizer"
	"github.com/arata-nvm/Solitude/compiler/parser"
	"io/ioutil"
	"log"
)

type Compiler struct {
	program *ast.Program
	isDebug bool
}

func New(isDebug bool) *Compiler {
	return &Compiler{
		isDebug: isDebug,
	}
}

func (c *Compiler) Compile(filename string) (errors []string) {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	c.program = p.ParseProgram()
	errors = p.Errors
	return
}

func (c *Compiler) Optimize() {
	o := optimizer.New(c.program)
	o.Optimize()
}

func (c *Compiler) GenIR() string {
	var b bytes.Buffer
	cg := codegen.New(c.program, c.isDebug, &b)
	cg.GenerateCode()
	return b.String()
}
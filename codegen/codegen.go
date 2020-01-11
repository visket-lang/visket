package codegen

import (
	"github.com/arata-nvm/Solitude/ast"
	"io"
)

type CodeGen struct {
	program      *ast.Program
	output       io.Writer
	index        int
	labelIndex   int
	isDebug      bool
	isTerminated bool
	context      *Context
}

func New(program *ast.Program, isDebug bool, w io.Writer) *CodeGen {
	c := &CodeGen{
		program: program,
		isDebug: isDebug,
		output:  w,
		context: NewContext(),
	}

	c.resetIndex()
	return c
}

func (c *CodeGen) GenerateCode() {
	c.genPrintFunction()
	for _, s := range c.program.Statements {
		c.genStatement(s)
	}
}

func (c *CodeGen) genPrintFunction() {
	c.gen("@.str = private unnamed_addr constant [4 x i8] c \"%%d\\0A\\00\", align 1\n")
	c.gen("declare i32 @printf(i8*, ...)\n")
	c.gen("define i32 @print(i32) nounwind {\n")
	c.gen("  call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str, i64 0, i64 0), i32 %%0)\n")
	c.gen("  ret i32 0\n")
	c.gen("}\n\n")
}

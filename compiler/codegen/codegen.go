package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
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

	module           *ir.Module
	contextFunction  *ir.Func
	contextBlock     *ir.Block
	contextCondAfter []*ir.Block
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
	c.genBool()
	c.genPrintFunction()
	c.genInputFunction()

	for _, s := range c.program.Statements {
		c.genStatement(s)
	}

	irCode := c.module.String()
	_, err := fmt.Fprint(c.output, irCode)
	if err != nil {
		errors.ErrorExit("failed writing ir code")
	}
}

func (c *CodeGen) genBool() {
	c.context.addVariableByName("false", constant.False)
	c.context.addVariableByName("true", constant.True)
}

func (c *CodeGen) genPrintFunction() {
	format := c.module.NewGlobalDef(".str.print", constant.NewCharArrayFromString("%d\x0A\x00"))
	format.Linkage = enum.LinkagePrivate
	format.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

	printf := c.module.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printf.Sig.Variadic = true

	printParam := ir.NewParam("", types.I32)
	print := c.module.NewFunc("print", types.I32, printParam)
	entryBlock := print.NewBlock("entry")

	zero := constant.NewInt(types.I64, 0)
	formatArg := constant.NewGetElementPtr(format.Typ.ElemType, format, zero, zero)
	entryBlock.NewCall(printf, formatArg, printParam)

	entryBlock.NewRet(constant.NewInt(types.I32, 0))

	c.context.addFunctionByName(print.Name(), print)
}

func (c *CodeGen) genInputFunction() {
	format := c.module.NewGlobalDef(".str.scanf", constant.NewCharArrayFromString("%d\x00"))
	format.Linkage = enum.LinkagePrivate
	format.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

	scanf := c.module.NewFunc("scanf", types.I32, ir.NewParam("", types.I8Ptr))
	scanf.Sig.Variadic = true

	input := c.module.NewFunc("input", types.I32)
	entryBlock := input.NewBlock("entry")

	scanfRet := entryBlock.NewAlloca(types.I32)

	zero := constant.NewInt(types.I64, 0)
	scanfArg := constant.NewGetElementPtr(format.Typ.ElemType, format, zero, zero)
	entryBlock.NewCall(scanf, scanfArg, scanfRet)

	result := entryBlock.NewLoad(types.I32, scanfRet)

	entryBlock.NewRet(result)

	c.context.addFunctionByName(input.Name(), input)

}

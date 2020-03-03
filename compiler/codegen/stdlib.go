package codegen

import (
	"github.com/arata-nvm/visket/compiler/codegen/builtin"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (c *CodeGen) genStdlib() {
	c.genGlibcFunc()
	c.genTypes()
	c.genPrintFunction()
	c.genInputFunction()
}

func (c CodeGen) genGlibcFunc() {
	{
		printf := c.module.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
		printf.Sig.Variadic = true
		c.context.addFunction("printf", printf)
	}

	{
		scanf := c.module.NewFunc("scanf", types.I32, ir.NewParam("", types.I8Ptr))
		scanf.Sig.Variadic = true
		c.context.addFunction("scanf", scanf)
	}
}

func (c *CodeGen) genTypes() {
	c.module.NewTypeDef("string", builtin.STRING)
	c.context.addType("string", builtin.STRING)

	c.context.addVariable("false", Value{
		Value:      constant.False,
		IsVariable: false,
	})
	c.context.addVariable("true", Value{
		Value:      constant.True,
		IsVariable: false,
	})
}

func (c *CodeGen) genPrintFunction() {
	format := c.module.NewGlobalDef(".str.print", constant.NewCharArrayFromString("%d\x0A\x00"))
	format.Linkage = enum.LinkagePrivate
	format.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

	printf, _ := c.context.findFunction("printf")

	printParam := ir.NewParam("", types.I32)
	print := c.module.NewFunc("print", types.I32, printParam)
	entryBlock := print.NewBlock("entry")

	zero := constant.NewInt(types.I64, 0)
	formatArg := constant.NewGetElementPtr(format.Typ.ElemType, format, zero, zero)
	entryBlock.NewCall(printf, formatArg, printParam)

	entryBlock.NewRet(constant.NewInt(types.I32, 0))

	c.context.addFunction(print.Name(), print)
}

func (c *CodeGen) genInputFunction() {
	format := c.module.NewGlobalDef(".str.scanf", constant.NewCharArrayFromString("%d\x00"))
	format.Linkage = enum.LinkagePrivate
	format.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

	scanf, _ := c.context.findFunction("scanf")

	input := c.module.NewFunc("input", types.I32)
	entryBlock := input.NewBlock("entry")

	scanfRet := entryBlock.NewAlloca(types.I32)

	zero := constant.NewInt(types.I64, 0)
	scanfArg := constant.NewGetElementPtr(format.Typ.ElemType, format, zero, zero)
	entryBlock.NewCall(scanf, scanfArg, scanfRet)

	result := entryBlock.NewLoad(types.I32, scanfRet)

	entryBlock.NewRet(result)

	c.context.addFunction(input.Name(), input)

}

package codegen

import (
	"github.com/arata-nvm/visket/compiler/codegen/builtin"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (c *CodeGen) genStdlib() {
	c.genGlibcFunc()
	c.genString()
}

func (c CodeGen) genGlibcFunc() {
	{
		printf := c.module.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
		printf.Sig.Variadic = true
		c.context.addFunction("printf", &Func{
			Func:        printf,
			IsReference: []bool{false},
		})
	}

	{
		scanf := c.module.NewFunc("scanf", types.I32, ir.NewParam("", types.I8Ptr))
		scanf.Sig.Variadic = true
		c.context.addFunction("scanf", &Func{
			Func:        scanf,
			IsReference: []bool{false},
		})
	}
}

func (c *CodeGen) genString() {
	c.module.NewTypeDef("string", builtin.STRING)
	c.context.addType("string", builtin.STRING)

	{
		strParam := ir.NewParam("", builtin.STRING)
		cstring := c.module.NewFunc("cstring", types.I8Ptr, strParam)
		block := cstring.NewBlock("entry")
		tmpVar := block.NewAlloca(builtin.STRING)
		block.NewStore(strParam, tmpVar)
		strVal := builtin.GetStringValue(tmpVar, block)
		block.NewRet(strVal)
		c.context.addFunction(cstring.Name(), &Func{
			Func:        cstring,
			IsReference: []bool{false},
		})
	}

	{
		strParam := ir.NewParam("", builtin.STRING)
		length := c.module.NewFunc("length", types.I32, strParam)
		block := length.NewBlock("entry")
		tmpVar := block.NewAlloca(builtin.STRING)
		block.NewStore(strParam, tmpVar)
		strLen := builtin.GetStringLength(tmpVar, block)
		block.NewRet(strLen)
		c.context.addFunction(length.Name(), &Func{
			Func:        length,
			IsReference: []bool{false},
		})
	}
}

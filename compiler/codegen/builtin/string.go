package builtin

import (
	"github.com/arata-nvm/visket/compiler/codegen/internal"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var STRING = newStringType()

func newStringType() types.Type {
	t := types.NewStruct(types.I8Ptr, types.I32)
	t.SetName("string")
	return t
}

func NewString(val string, block *ir.Block, module *ir.Module) value.Value {
	varName := internal.NextString()
	constStr := module.NewGlobalDef(varName, constant.NewCharArrayFromString(val+"\x00"))
	constStr.Linkage = enum.LinkagePrivate
	constStr.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

	zero := constant.NewInt(types.I32, 0)
	one := constant.NewInt(types.I32, 1)

	str := block.NewAlloca(STRING)
	strVal := block.NewGetElementPtr(STRING, str, zero, zero)
	block.NewStore(block.NewGetElementPtr(internal.PtrElmType(constStr), constStr, zero, zero), strVal)
	strLen := block.NewGetElementPtr(STRING, str, zero, one)
	block.NewStore(constant.NewInt(types.I32, int64(len(val))), strLen)

	return str
}

func GetStringValue(v value.Value, block *ir.Block) value.Value {
	zero := constant.NewInt(types.I32, 0)
	strValPtr := block.NewGetElementPtr(STRING, v, zero, zero)
	strVal := block.NewLoad(types.I8Ptr, strValPtr)
	return strVal
}

func GetStringLength(v value.Value, block *ir.Block) value.Value {
	zero := constant.NewInt(types.I32, 0)
	one := constant.NewInt(types.I32, 1)
	strLenPtr := block.NewGetElementPtr(STRING, v, zero, one)
	strLen := block.NewLoad(types.I32, strLenPtr)
	return strLen
}

func GetIndexedStringValue(v value.Value, index value.Value, block *ir.Block) value.Value {
	strVal := GetStringValue(v, block)
	return block.NewGetElementPtr(types.I8, strVal, index)
}

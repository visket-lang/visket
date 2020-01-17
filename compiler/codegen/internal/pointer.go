package internal

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func PtrElmType(v value.Value) types.Type {
	return v.Type().(*types.PointerType).ElemType
}

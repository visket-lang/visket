package internal

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func PtrElmType(v value.Value) types.Type {
	if ptr, ok := v.Type().(*types.PointerType); ok {
		return ptr.ElemType
	}
	return v.Type()
}

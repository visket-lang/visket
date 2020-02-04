package types

import (
	"bytes"
	"fmt"
	llvmTypes "github.com/llir/llvm/ir/types"
)

type SlType interface {
	isSlType()
	LlvmType() llvmTypes.Type
	String() string
}

type SlPointer struct {
	baseType SlType
}

func NewSlPointer(baseType SlType) *SlPointer {
	return &SlPointer{
		baseType: baseType,
	}
}

func (p *SlPointer) LlvmType() llvmTypes.Type {
	return llvmTypes.NewPointer(p.baseType.LlvmType())
}

func (p *SlPointer) String() string {
	return fmt.Sprintf("*%s", p.baseType)
}

func (p *SlPointer) isSlType() {}

type SlFunction struct {
	RetType SlType
	Params  []SlType
}

func NewSlFunction(retType SlType, params []SlType) *SlFunction {
	return &SlFunction{
		RetType: retType,
		Params:  params,
	}
}

func (f *SlFunction) LlvmType() llvmTypes.Type {
	var params []llvmTypes.Type
	for _, p := range f.Params {
		params = append(params, p.LlvmType())
	}

	return llvmTypes.NewFunc(f.RetType.LlvmType(), params...)
}

func (f *SlFunction) String() string {
	var params bytes.Buffer
	for i, p := range f.Params {
		if i != 0 {
			params.WriteString(", ")
		}
		params.WriteString(p.String())
	}
	return fmt.Sprintf("%s -> %s", params.String(), f.RetType)
}

func (f *SlFunction) isSlType() {}

type SlVoid struct {
}

func NewSlVoid() *SlVoid {
	return &SlVoid{}
}

func (v *SlVoid) LlvmType() llvmTypes.Type {
	return llvmTypes.Void
}

func (v *SlVoid) String() string {
	return "void"
}

func (v *SlVoid) isSlType() {}

type SlInt struct {
}

func NewSlInt() *SlInt {
	return &SlInt{}
}

func (i *SlInt) LlvmType() llvmTypes.Type {
	return llvmTypes.I32
}

func (i *SlInt) String() string {
	return "int"
}

func (i *SlInt) isSlType() {}

type SlFloat struct {
}

func NewSlFloat() *SlFloat {
	return &SlFloat{}
}

func (f *SlFloat) LlvmType() llvmTypes.Type {
	return llvmTypes.Float
}

func (f *SlFloat) String() string {
	return "float"
}

func (f *SlFloat) isSlType() {}

type SlBool struct {
}

func NewSlBool() *SlBool {
	return &SlBool{}
}

func (b *SlBool) LlvmType() llvmTypes.Type {
	return llvmTypes.I1
}

func (b *SlBool) String() string {
	return "bool"
}

func (b *SlBool) isSlType() {}

type SlArray struct {
	Len    uint64
	ElmTyp SlType
}

func NewSlArray(len int, elmTyp SlType) *SlArray {
	return &SlArray{
		Len:    uint64(len),
		ElmTyp: elmTyp,
	}
}

func (a *SlArray) LlvmType() llvmTypes.Type {
	return llvmTypes.NewArray(a.Len, a.ElmTyp.LlvmType())
}

func (a *SlArray) String() string {
	return fmt.Sprintf("[]%s", a.ElmTyp)
}

func (a *SlArray) isSlType() {}

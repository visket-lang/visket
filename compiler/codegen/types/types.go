package types

import "fmt"

var (
	I1  = IntType("i1")
	I32 = IntType("i32")

	I1Ptr  = PointerType{ElmType: I1}
	I32Ptr = PointerType{ElmType: I32}
)

type Types interface {
	Name() string
	String() string
}

type IntType string

func (it IntType) Name() string {
	return string(it)
}

func (it IntType) String() string {
	return it.Name()
}

type PointerType struct {
	ElmType Types
}

func NewPointer(types Types) PointerType {
	return PointerType{ElmType: types}
}

func (pt PointerType) Name() string {
	return fmt.Sprintf("%s*", pt.ElmType)
}

func (pt PointerType) String() string {
	return pt.Name()
}

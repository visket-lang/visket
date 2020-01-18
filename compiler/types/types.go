package types

import (
	"bytes"
	"fmt"
	"github.com/llir/llvm/ir/types"
)

var (
	VOID = VoidType{}
	INT  = IntType{}
)

type ParserType interface {
	ToType() types.Type
	String() string
}

type VoidType struct{}

func (v VoidType) ToType() types.Type {
	return types.Void
}

func (v VoidType) String() string {
	return "void"
}

type IntType struct{}

func (i IntType) ToType() types.Type {
	return types.I32
}

func (i IntType) String() string {
	return "int"
}

type FuncType struct {
	RetType ParserType
	Params  []ParserType
}

func NewFuncType(retType ParserType, params []ParserType) FuncType {
	return FuncType{
		RetType: retType,
		Params:  params,
	}
}

func (f FuncType) ToType() types.Type {
	var params []types.Type
	for _, p := range f.Params {
		params = append(params, p.ToType())
	}

	return types.NewFunc(f.RetType.ToType(), params...)
}

func (f FuncType) String() string {
	var params bytes.Buffer
	for i, p := range f.Params {
		if i != 0 {
			params.WriteString(", ")
		}
		params.WriteString(p.String())
	}
	return fmt.Sprintf("%s -> %s", params.String(), f.RetType)
}

var typeNameToType = map[string]ParserType{
	"void": VOID,
	"int":  INT,
}

func ParseType(name string) ParserType {
	if typ, ok := typeNameToType[name]; ok {
		return typ
	}

	return nil
}

package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"log"
	"strings"
)

type Pointer int
type Value int
type Label string

func (c *CodeGen) resetIndex() {
	c.index = -1
	c.labelIndex = -1
}

func (c *CodeGen) nextPointer() Pointer {
	c.index++
	return Pointer(c.index)
}

func (c *CodeGen) nextValue() Value {
	c.index++
	return Value(c.index)
}

func (c *CodeGen) nextLabel(name string) Label {
	c.labelIndex++
	return Label(fmt.Sprintf("%s.%d", name, c.labelIndex))
}

func (c *CodeGen) gen(format string, a ...interface{}) {
	code := fmt.Sprintf(format, a...)
	_, err := fmt.Fprintf(c.output, format, a...)
	if err != nil {
		log.Fatal(err)
	}

	c.isTerminated = strings.Contains(code, "ret") || strings.Contains(code, "br")
}

func (c *CodeGen) comment(format string, a ...interface{}) {
	if !c.isDebug {
		return
	}

	c.gen("")
	c.gen(format, a...)
}

func (c *CodeGen) genAlloca() Pointer {
	result := c.nextPointer()
	c.gen("  %%%d = alloca i32, align 4\n", result)
	return result
}

func (c *CodeGen) genNamedAlloca(ident *ast.Identifier) {
	c.gen("  %%%s = alloca i32, align 4\n", ident.Token.Literal)
}

func (c *CodeGen) genStore(value Value, ptrToStore Pointer) {
	c.gen("  store i32 %%%d, i32* %%%d\n", value, ptrToStore)
}

func (c *CodeGen) genNamedStore(ident *ast.Identifier, ptrToStore Pointer) {
	c.gen("  store i32 %%%d, i32* %%%s\n", ptrToStore, ident.Token.Literal)
}

func (c *CodeGen) genStoreImmediate(value int, ptrToStore Pointer) {
	c.gen("  store i32 %d, i32* %%%d\n", value, ptrToStore)
}

func (c *CodeGen) genLoad(ptrToLoad Pointer) Value {
	result := c.nextValue()
	c.gen("  %%%d = load i32, i32* %%%d, align 4\n", result, ptrToLoad)
	return result
}

func (c *CodeGen) genNamedLoad(ident *ast.Identifier) Value {
	result := c.nextValue()
	c.gen("  %%%d = load i32, i32* %%%s, align 4\n", result, ident.Token.Literal)
	return result
}

func (c *CodeGen) genAdd(op1 Value, op2 Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = add i32 %%%d, %%%d\n", result, op1, op2)
	return result
}

func (c *CodeGen) genSub(op1 Value, op2 Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = sub i32 %%%d, %%%d\n", result, op1, op2)
	return result
}

func (c *CodeGen) genMul(op1 Value, op2 Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = mul i32 %%%d, %%%d\n", result, op1, op2)
	return result
}

func (c *CodeGen) genIDiv(op1 Value, op2 Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = idiv i32 %%%d, %%%d\n", result, op1, op2)
	return result
}

type IcmpCond string

const (
	EQ  IcmpCond = "eq"
	NEQ          = "ne"
	LT           = "slt"
	LTE          = "sle"
	GT           = "sgt"
	GTE          = "sge"
)

func (c *CodeGen) genIcmp(cond IcmpCond, op1, op2 Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = icmp %s i32 %%%d, %%%d\n", result, cond, op1, op2)
	return result
}

func (c *CodeGen) genZext(typeFrom, typeTo string, value Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = zext %s %%%d to %s\n", result, typeFrom, value, typeTo)
	return result
}

func (c *CodeGen) genRet(value Value) {
	c.gen("  ret i32 %%%d\n", value)
}

func (c *CodeGen) genDefineFunction(ident *ast.Identifier) {
	c.gen("define i32 @%s(", ident.Token.Literal)
}

func (c *CodeGen) genFunctionParameters(params []*ast.Identifier) {
	var p []string
	for _, _ = range params {
		p = append(p, "i32")
	}

	c.gen(strings.Join(p, ","))
}

func (c *CodeGen) genBeginFunction() {
	c.gen(") nounwind {\n")
}

func (c *CodeGen) genEndFunction() {
	c.gen("}\n\n")
}

func (c *CodeGen) genCall(function *ast.Identifier, params []Value) {
	var p []string
	for _, param := range params {
		p = append(p, fmt.Sprintf("i32 %%%d", param))
	}

	c.gen("call i32 @%s(%s)\n", function.Token.Literal, strings.Join(p, ","))
}

func (c *CodeGen) genCallWithReturn(function *ast.Identifier, params []Value) Value {
	result := c.nextValue()

	var p []string
	for _, param := range params {
		p = append(p, fmt.Sprintf("i32 %%%d", param))
	}

	c.gen("  %%%d = call i32 @%s(%s)\n", result, function.Token.Literal, strings.Join(p, ","))
	return result
}

func (c *CodeGen) genLabel(name Label) {
	c.gen("%s:\n", name)
}

func (c *CodeGen) genBr(label Label) {
	c.gen("  br label %%%s\n", label)
}

func (c *CodeGen) genBrWithCond(condition Value, ifTrue Label, itFalse Label) {
	c.gen("  br i1 %%%d, label %%%s, label %%%s\n", condition, ifTrue, itFalse)
}

func (c *CodeGen) genTrunc(typeFrom, typeTo string, value Value) Value {
	result := c.nextValue()
	c.gen("  %%%d = trunc %s %%%d to %s\n", result, typeFrom, value, typeTo)
	return result
}

package codegen

import (
	"fmt"
	"log"
)

type Pointer int
type Value int

func (c *CodeGen) nextPointer() Pointer {
	c.index++
	return Pointer(c.index)
}

func (c *CodeGen) nextValue() Value {
	c.index++
	return Value(c.index)
}

func (c *CodeGen) gen(format string, a ...interface{}) {
	_, err := fmt.Fprintf(c.output, format, a...)
	if err != nil {
		log.Fatal(err)
	}
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

func (c *CodeGen) genStore(value Value, ptrToStore Pointer) {
	c.gen("  store i32 %%%d, i32* %%%d\n", value, ptrToStore)
}

func (c *CodeGen) genStoreImmediate(value int, ptrToStore Pointer) {
	c.gen("  store i32 %d, i32* %%%d\n", value, ptrToStore)
}

func (c *CodeGen) genLoad(ptrToLoad Pointer) Value {
	result := c.nextValue()
	c.gen("  %%%d = load i32, i32* %%%d, align 4\n", result, ptrToLoad)
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

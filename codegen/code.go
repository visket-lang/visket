package codegen

import (
	"fmt"
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
	fmt.Printf(format, a...)
}

func (c *CodeGen) comment(format string, a ...interface{}) {
	if !c.isDebug {
		return
	}

	fmt.Println("")
	fmt.Printf(format, a...)
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

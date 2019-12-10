package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/parser"
)

type Pointer int
type Value int

type CodeGen struct {
	p     *parser.Parser
	index int
}

func New(p *parser.Parser) *CodeGen {
	c := &CodeGen{
		p: p,
	}

	return c
}

func (c *CodeGen) GenerateCode() {
	program := c.p.ParseProgram()
	c.gen("define i32 @main() nounwind {\n")
	c.genExpr(program.Code)

	c.comment("  ; Ret\n")
	numProcLastIndex := c.index
	regIndex := c.genLoad(numProcLastIndex)
	c.gen("  ret i32 %%%d\n", regIndex)
	c.gen("}\n")
}

func (c *CodeGen) genExpr(node ast.Node) {
	switch node := node.(type) {
	case *ast.InfixExpression:
		c.genInfix(node)
	case *ast.IntegerLiteral:
		c.comment("  ; Assign\n")
		c.genAlloca()
		c.genStoreImmediate(node.Value, c.index)
	}
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) {
	c.comment("  ; Infix\n")

	c.genExpr(ie.Left)
	lhsIndex := c.index
	c.genExpr(ie.Right)
	rhsIndex := c.index

	c.comment("  ; Op\n")

	lhsRegIndex := c.genLoad(lhsIndex)
	rhsRegIndex := c.genLoad(rhsIndex)

	var resRegIndex int

	switch ie.Operator {
	case "+":
		c.comment("  ; Add\n")
		resRegIndex = c.genAdd(lhsRegIndex, rhsRegIndex)
	case "-":
		c.comment("  ; Sub\n")
		resRegIndex = c.genSub(lhsRegIndex, rhsRegIndex)
	}

	resMemIndex := c.genAlloca()
	c.genStore(resRegIndex, resMemIndex)
}

func (c *CodeGen) gen(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (c *CodeGen) comment(format string, a ...interface{}) {
	fmt.Println("")
	fmt.Printf(format, a...)
}

func (c *CodeGen) genAlloca() int {
	c.index++
	c.gen("  %%%d = alloca i32, align 4\n", c.index)
	return c.index
}

func (c *CodeGen) genStore(pointer1 int, pointer2 int) {
	c.gen("  store i32 %%%d, i32* %%%d\n", pointer1, pointer2)
}

func (c *CodeGen) genStoreImmediate(value int, pointer int) {
	c.gen("  store i32 %d, i32* %%%d\n", value, pointer)
}

func (c *CodeGen) genLoad(pointer int) int {
	c.index++
	c.gen("  %%%d = load i32, i32* %%%d, align 4\n", c.index, pointer)
	return c.index
}

func (c *CodeGen) genAdd(op1 int, op2 int) int {
	c.index++
	c.gen("  %%%d = add i32 %%%d, %%%d\n", c.index, op1, op2)
	return c.index
}

func (c *CodeGen) genSub(op1 int, op2 int) int {
	c.index++
	c.gen("  %%%d = sub i32 %%%d, %%%d\n", c.index, op1, op2)
	return c.index
}

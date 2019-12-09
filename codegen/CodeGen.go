package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/parser"
)

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
	fmt.Println("define i32 @main() nounwind {")

	c.gen(program.Code)

	numProcLastIndex := c.index
	c.index++
	regIndex := c.index

	fmt.Printf("  %%%d = load i32, i32* %%%d, align 4\n", regIndex, numProcLastIndex)
	fmt.Printf("  ret i32 %%%d\n", regIndex)
	fmt.Println("}")
}

func (c *CodeGen) gen(node ast.Node) {
	switch node := node.(type) {
	case *ast.InfixExpression:
		c.genInfix(node)
	case *ast.IntegerLiteral:
		c.index++
		fmt.Println("  ; Assign")
		fmt.Printf("  %%%d = alloca i32, align 4\n", c.index)
		fmt.Printf("  store i32 %d, i32* %%%d\n", node.Value, c.index)
	}
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) {
	c.gen(ie.Left)
	lhsIndex := c.index
	c.gen(ie.Right)
	rhsIndex := c.index

	c.index++
	lhsRegIndex := c.index
	fmt.Printf("  %%%d= load i32, i32* %%%d, align 4\n", c.index, lhsIndex)

	c.index++
	rhsRegIndex := c.index
	fmt.Printf("  %%%d = load i32, i32* %%%d, align 4\n", c.index, rhsIndex)

	c.index++
	resRegIndex := c.index

	switch ie.Operator {
	case "+":
		fmt.Printf("  %%%d = add i32 %%%d, %%%d\n", c.index, lhsRegIndex, rhsRegIndex)
	case "-":
		fmt.Printf("  %%%d = sub i32 %%%d, %%%d\n", c.index, lhsRegIndex, rhsRegIndex)
	}

	c.index++
	resMemIndex := c.index
	fmt.Printf("  %%%d = alloca i32, align 4\n", c.index)

	fmt.Printf("  store i32 %%%d, i32* %%%d, align 4\n", resRegIndex, resMemIndex)
}

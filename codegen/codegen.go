package codegen

import (
	"github.com/arata-nvm/Solitude/ast"
)

type CodeGen struct {
	program *ast.Program
	index   int
	isDebug bool
}

func New(program *ast.Program, isDebug bool) *CodeGen {
	c := &CodeGen{
		program: program,
		isDebug: isDebug,
	}

	return c
}

func (c *CodeGen) GenerateCode() {
	c.gen("define i32 @main() nounwind {\n")

	result := c.genExpr(c.program.Code)

	c.comment("  ; Ret\n")
	returnPtr := c.genLoad(result)
	c.gen("  ret i32 %%%d\n", returnPtr)
	c.gen("}\n")
}

func (c *CodeGen) genExpr(node ast.Node) Pointer {
	var result Pointer
	switch node := node.(type) {
	case *ast.InfixExpression:
		result = c.genInfix(node)
	case *ast.IntegerLiteral:
		c.comment("  ; Assign\n")
		result = c.genAlloca()
		c.genStoreImmediate(node.Value, result)
	}

	return result
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) Pointer {
	c.comment("  ; Infix\n")

	lhsPtr := c.genExpr(ie.Left)
	rhsPtr := c.genExpr(ie.Right)

	c.comment("  ; Op\n")

	lhs := c.genLoad(lhsPtr)
	rhs := c.genLoad(rhsPtr)

	var result Value

	switch ie.Operator {
	case "+":
		c.comment("  ; Add\n")
		result = c.genAdd(lhs, rhs)
	case "-":
		c.comment("  ; Sub\n")
		result = c.genSub(lhs, rhs)
	}

	resultPtr := c.genAlloca()
	c.genStore(result, resultPtr)

	return resultPtr
}

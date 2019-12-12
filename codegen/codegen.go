package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"io"
	"os"
)

type CodeGen struct {
	program *ast.Program
	index   int
	isDebug bool
	output  io.Writer
}

func New(program *ast.Program, isDebug bool, w io.Writer) *CodeGen {
	c := &CodeGen{
		program: program,
		isDebug: isDebug,
		output:  w,
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
	default:
		fmt.Printf("unexpexted node: %s\n", node.Inspect())
		os.Exit(1)
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
	case "*":
		c.comment("  ; Mul\n")
		result = c.genMul(lhs, rhs)
	case "/":
		c.comment("  ; Div\n")
		result = c.genIDiv(lhs, rhs)
	case "==":
		c.comment("  ; Equal\n")
		result = c.genIcmp(EQ, lhs, rhs)
		result = c.genZext("i1", "i32", result)
	case "!=":
		c.comment("  ; Not Equal\n")
		result = c.genIcmp(NEQ, lhs, rhs)
		result = c.genZext("i1", "i32", result)
	case "<":
		c.comment("  ; Less Than\n")
		result = c.genIcmp(LT, lhs, rhs)
		result = c.genZext("i1", "i32", result)
	case "<=":
		c.comment("  ; Less Than or Equal\n")
		result = c.genIcmp(LTE, lhs, rhs)
		result = c.genZext("i1", "i32", result)
	case ">":
		c.comment("  ; Greater Than\n")
		result = c.genIcmp(GT, lhs, rhs)
		result = c.genZext("i1", "i32", result)
	case ">=":
		c.comment("  ; Greater Than or Equal\n")
		result = c.genIcmp(GTE, lhs, rhs)
		result = c.genZext("i1", "i32", result)
	}

	resultPtr := c.genAlloca()
	c.genStore(result, resultPtr)

	return resultPtr
}

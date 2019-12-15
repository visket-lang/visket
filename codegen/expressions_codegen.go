package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"os"
)

func (c *CodeGen) genExpression(expr ast.Expression) Value {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return c.genInfix(expr)
	case *ast.CallExpression:
		return c.genCallExpression(expr)
	case *ast.IntegerLiteral:
		c.comment("  ; Int\n")
		result := c.genAlloca()
		c.genStoreImmediate(expr.Value, result)
		return c.genLoad(result)
	case *ast.Identifier:
		c.comment("  ; Ident\n")
		return c.genNamedLoad(expr)
	}

	fmt.Printf("unexpexted expression: %s\n", expr.Inspect())
	os.Exit(1)
	return -1
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) Value {
	c.comment("  ; Infix\n")

	lhs := c.genExpression(ie.Left)
	rhs := c.genExpression(ie.Right)

	c.comment("  ; Op\n")

	switch ie.Operator {
	case "+":
		c.comment("  ; Add\n")
		return c.genAdd(lhs, rhs)
	case "-":
		c.comment("  ; Sub\n")
		return c.genSub(lhs, rhs)
	case "*":
		c.comment("  ; Mul\n")
		return c.genMul(lhs, rhs)
	case "/":
		c.comment("  ; Div\n")
		return c.genIDiv(lhs, rhs)
	case "==":
		c.comment("  ; Equal\n")
		result := c.genIcmp(EQ, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case "!=":
		c.comment("  ; Not Equal\n")
		result := c.genIcmp(NEQ, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case "<":
		c.comment("  ; Less Than\n")
		result := c.genIcmp(LT, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case "<=":
		c.comment("  ; Less Than or Equal\n")
		result := c.genIcmp(LTE, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case ">":
		c.comment("  ; Greater Than\n")
		result := c.genIcmp(GT, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case ">=":
		c.comment("  ; Greater Than or Equal\n")
		result := c.genIcmp(GTE, lhs, rhs)
		return c.genZext("i1", "i32", result)
	}

	return Value(c.index)
}

func (c *CodeGen) genCallExpression(expr *ast.CallExpression) Value {
	var params []Value

	for _, param := range expr.Parameters {
		params = append(params, c.genExpression(param))
	}

	return c.genCallWithReturn(expr.Function, params)
}

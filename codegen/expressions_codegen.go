package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/codegen/constant"
	"github.com/arata-nvm/Solitude/codegen/types"
	"os"
)

func (c *CodeGen) genExpression(expr ast.Expression) Value {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return c.genInfix(expr)
	case *ast.CallExpression:
		return c.genCallExpression(expr)
	case *ast.IntegerLiteral:
		return constant.NewInt(types.I32, expr.Value)
	case *ast.Identifier:
		c.comment("  ; RegName\n")
		v, ok := c.context.findVariable(expr)
		if !ok {
			fmt.Printf("unresolved variable: %s\n", expr.String())
			os.Exit(1)
		}
		return c.genLoad(types.I32, v)
	}

	fmt.Fprintf(os.Stderr, "unexpexted expression: %s\n", expr.Inspect())
	os.Exit(1)
	return Register{}
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
		c.comment("  ; Quo\n")
		return c.genSDiv(lhs, rhs)
	case "%":
		c.comment("  ; Rem\n")
		return c.genSRem(lhs, rhs)
	case "<<":
		c.comment("  ; Shl\n")
		return c.genShl(lhs, rhs)
	case ">>":
		c.comment("  ; Shr\n")
		return c.genAShr(lhs, rhs)
	case "==":
		c.comment("  ; Equal\n")
		result := c.genIcmp(EQ, lhs, rhs)
		return c.genZext(types.I32, result)
	case "!=":
		c.comment("  ; Not Equal\n")
		result := c.genIcmp(NEQ, lhs, rhs)
		return c.genZext(types.I32, result)
	case "<":
		c.comment("  ; Less Than\n")
		result := c.genIcmp(LT, lhs, rhs)
		return c.genZext(types.I32, result)
	case "<=":
		c.comment("  ; Less Than or Equal\n")
		result := c.genIcmp(LTE, lhs, rhs)
		return c.genZext(types.I32, result)
	case ">":
		c.comment("  ; Greater Than\n")
		result := c.genIcmp(GT, lhs, rhs)
		return c.genZext(types.I32, result)
	case ">=":
		c.comment("  ; Greater Than or Equal\n")
		result := c.genIcmp(GTE, lhs, rhs)
		return c.genZext(types.I32, result)
	}

	fmt.Fprintf(os.Stdout, "unexpected operator: %s\n", ie.Operator)
	os.Exit(1)
	return Register{}
}

func (c *CodeGen) genCallExpression(expr *ast.CallExpression) Value {
	var params []Value

	for _, param := range expr.Parameters {
		params = append(params, c.genExpression(param))
	}

	return c.genCallWithReturn(expr.Function, params)
}

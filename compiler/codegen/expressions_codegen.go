package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen/internal"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (c *CodeGen) genExpression(expr ast.Expression) value.Value {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return c.genInfix(expr)
	case *ast.CallExpression:
		return c.genCallExpression(expr)
	case *ast.IntegerLiteral:
		return constant.NewInt(types.I32, int64(expr.Value))
	case *ast.Identifier:
		v, ok := c.context.findVariable(expr)
		if !ok {
			errors.ErrorExit(fmt.Sprintf("unresolved variable: %s\n", expr.String()))
		}
		_, ok = v.Type().(*types.PointerType)
		if ok {
			return c.contextBlock.NewLoad(internal.PtrElmType(v), v)
		} else {
			return v
		}
	}

	errors.ErrorExit(fmt.Sprintf("unexpexted expression: %s\n", expr.Inspect()))
	return nil
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) value.Value {
	lhs := c.genExpression(ie.Left)
	rhs := c.genExpression(ie.Right)

	switch ie.Operator {
	case "+":
		return c.contextBlock.NewAdd(lhs, rhs)
	case "-":
		return c.contextBlock.NewSub(lhs, rhs)
	case "*":
		return c.contextBlock.NewMul(lhs, rhs)
	case "/":
		return c.contextBlock.NewSDiv(lhs, rhs)
	case "%":
		return c.contextBlock.NewSRem(lhs, rhs)
	case "<<":
		return c.contextBlock.NewShl(lhs, rhs)
	case ">>":
		return c.contextBlock.NewAShr(lhs, rhs)
	case "==":
		return c.contextBlock.NewICmp(enum.IPredEQ, lhs, rhs)
	case "!=":
		return c.contextBlock.NewICmp(enum.IPredNE, lhs, rhs)
	case "<":
		return c.contextBlock.NewICmp(enum.IPredULT, lhs, rhs)
	case "<=":
		return c.contextBlock.NewICmp(enum.IPredULE, lhs, rhs)
	case ">":
		return c.contextBlock.NewICmp(enum.IPredUGT, lhs, rhs)
	case ">=":
		return c.contextBlock.NewICmp(enum.IPredUGE, lhs, rhs)
	}

	errors.ErrorExit(fmt.Sprintf("unexpected operator: %s\n", ie.Operator))
	return nil
}

func (c *CodeGen) genCallExpression(expr *ast.CallExpression) value.Value {
	var params []value.Value

	for _, param := range expr.Parameters {
		params = append(params, c.genExpression(param))
	}

	f, ok := c.context.findFunction(expr.Function)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | undefined function %s", expr.Token.Pos, expr.Function.String()))
	}

	return c.contextBlock.NewCall(f, params...)
}

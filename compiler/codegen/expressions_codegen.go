package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/codegen/internal"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/arata-nvm/Solitude/compiler/token"
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
	case *ast.AssignExpression:
		return c.genAssignExpression(expr)
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
	f, ok := c.context.findFunction(expr.Function)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | undefined function %s", expr.Token.Pos, expr.Function.String()))
	}

	if len(expr.Parameters) < len(f.Params) {
		errors.ErrorExit(fmt.Sprintf("%s | not enough arguments in call to %s", expr.Token.Pos, expr.Function.String()))
	} else if len(expr.Parameters) > len(f.Params) {
		errors.ErrorExit(fmt.Sprintf("%s | too many arguments in call to %s", expr.Token.Pos, expr.Function.String()))
	}

	var params []value.Value

	for i, param := range expr.Parameters {
		v := c.genExpression(param)
		params = append(params, v)
		if v.Type() != f.Sig.Params[i] {
			errors.ErrorExit(fmt.Sprintf("%s | type mismatched %s and %s", expr.Token.Pos, v.Type(), f.Sig.Params[i].String()))

		}
	}

	return c.contextBlock.NewCall(f, params...)
}

func (c *CodeGen) genAssignExpression(stmt *ast.AssignExpression) value.Value {
	// TODO rewrite
	ident, ok := stmt.Left.(*ast.Identifier)
	if !ok {
		errors.ErrorExit("not implemented")
	}

	lhs, ok := c.context.findVariable(ident)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("unresolved variable: %s\n", ident))
	}

	rhs := c.genExpression(stmt.Value)

	lhsTyp := internal.PtrElmType(lhs)
	rhsTyp := rhs.Type()

	if !lhsTyp.Equal(rhsTyp) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Token.Pos, lhsTyp, rhsTyp))
	}

	switch stmt.Token.Type {
	case token.ASSIGN:
		c.contextBlock.NewStore(rhs, lhs)
	case token.ADD_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewAdd(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	case token.SUB_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewSub(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	case token.MUL_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewMul(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	case token.QUO_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewSDiv(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	case token.REM_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewSRem(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	case token.SHL_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewShl(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	case token.SHR_ASSIGN:
		vValue := c.contextBlock.NewLoad(lhsTyp, lhs)
		rhs = c.contextBlock.NewAShr(vValue, rhs)
		c.contextBlock.NewStore(rhs, lhs)
	}

	return rhs
}

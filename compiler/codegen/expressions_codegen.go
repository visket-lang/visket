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

func (c *CodeGen) genExpression(expr ast.Expression) Value {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return c.genInfix(expr)
	case *ast.CallExpression:
		return c.genCallExpression(expr)
	case *ast.AssignExpression:
		return c.genAssignExpression(expr)
	case *ast.IntegerLiteral:
		return Value{
			Value:      constant.NewInt(types.I32, int64(expr.Value)),
			IsVariable: false,
		}
	case *ast.Identifier:
		v, ok := c.context.findVariable(expr)
		if !ok {
			errors.ErrorExit(fmt.Sprintf("unresolved variable: %s\n", expr.String()))
		}

		return v
	}

	errors.ErrorExit(fmt.Sprintf("unexpexted expression: %s\n", expr.Inspect()))
	return Value{} //unreachable
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) Value {
	lhs := c.genExpression(ie.Left).Load(c.contextBlock)
	rhs := c.genExpression(ie.Right).Load(c.contextBlock)

	var opResult value.Value

	switch ie.Operator {
	case "+":
		opResult = c.contextBlock.NewAdd(lhs, rhs)
	case "-":
		opResult = c.contextBlock.NewSub(lhs, rhs)
	case "*":
		opResult = c.contextBlock.NewMul(lhs, rhs)
	case "/":
		opResult = c.contextBlock.NewSDiv(lhs, rhs)
	case "%":
		opResult = c.contextBlock.NewSRem(lhs, rhs)
	case "<<":
		opResult = c.contextBlock.NewShl(lhs, rhs)
	case ">>":
		opResult = c.contextBlock.NewAShr(lhs, rhs)
	case "==":
		opResult = c.contextBlock.NewICmp(enum.IPredEQ, lhs, rhs)
	case "!=":
		opResult = c.contextBlock.NewICmp(enum.IPredNE, lhs, rhs)
	case "<":
		opResult = c.contextBlock.NewICmp(enum.IPredULT, lhs, rhs)
	case "<=":
		opResult = c.contextBlock.NewICmp(enum.IPredULE, lhs, rhs)
	case ">":
		opResult = c.contextBlock.NewICmp(enum.IPredUGT, lhs, rhs)
	case ">=":
		opResult = c.contextBlock.NewICmp(enum.IPredUGE, lhs, rhs)
	default:
		errors.ErrorExit(fmt.Sprintf("unexpected operator: %s\n", ie.Operator))
	}

	return Value{
		Value:      opResult,
		IsVariable: false,
	}
}

func (c *CodeGen) genCallExpression(expr *ast.CallExpression) Value {
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
		v := c.genExpression(param).Load(c.contextBlock)
		params = append(params, v)
		if v.Type() != f.Sig.Params[i] {
			errors.ErrorExit(fmt.Sprintf("%s | type mismatched %s and %s", expr.Token.Pos, v.Type(), f.Sig.Params[i].String()))

		}
	}

	funcRet := c.contextBlock.NewCall(f, params...)

	return Value{
		Value:      funcRet,
		IsVariable: false,
	}
}
func (c *CodeGen) genAssignExpression(stmt *ast.AssignExpression) Value {
	lhs := c.genExpression(stmt.Left).Value
	rhs := c.genExpression(stmt.Value).Load(c.contextBlock)

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

	return Value{
		Value:      rhs,
		IsVariable: true,
	}
}

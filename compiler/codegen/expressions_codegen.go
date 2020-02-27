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

func (c *CodeGen) genExpression(expr ast.Expression) Value {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return c.genInfix(expr)
	case *ast.CallExpression:
		return c.genCallExpression(expr)
	case *ast.IndexExpression:
		return c.genIndexExpression(expr)
	case *ast.AssignExpression:
		return c.genAssignExpression(expr)
	case *ast.IntegerLiteral:
		return c.genIntegerLiteral(expr)
	case *ast.FloatLiteral:
		return c.genFloatLiteral(expr)
	case *ast.Identifier:
		return c.genIdentifier(expr)
	case *ast.NewExpression:
		return c.genNewExpression(expr)
	case *ast.LoadMemberExpression:
		return c.genLoadMemberExpression(expr)
	}

	errors.ErrorExit(fmt.Sprintf("unexpexted expression: %s\n", expr.Inspect()))
	return Value{} //unreachable
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) Value {
	lhs := c.genExpression(ie.Left).Load(c.contextBlock)
	rhs := c.genExpression(ie.Right).Load(c.contextBlock)

	lhsTyp := lhs.Type()
	rhsTyp := rhs.Type()
	if lhsTyp.Equal(types.I32) && rhsTyp.Equal(types.I32) {
		return c.genInfixInteger(ie.Operator, lhs, rhs)
	} else if lhsTyp.Equal(types.Float) && rhsTyp.Equal(types.Float) {
		return c.genInfixFloat(ie.Operator, lhs, rhs)
	}

	errors.ErrorExit(fmt.Sprintf("unexpected operator: %s %s %s\n", lhsTyp, ie.Operator, rhsTyp))
	return Value{} // unreachable
}

func (c *CodeGen) genInfixInteger(op string, lhs value.Value, rhs value.Value) Value {
	var opResult value.Value

	switch op {
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
		errors.ErrorExit(fmt.Sprintf("unexpected operator: int %s int\n", op))
	}

	return Value{
		Value:      opResult,
		IsVariable: false,
	}
}

func (c *CodeGen) genInfixFloat(op string, lhs value.Value, rhs value.Value) Value {
	var opResult value.Value

	switch op {
	case "+":
		opResult = c.contextBlock.NewFAdd(lhs, rhs)
	case "-":
		opResult = c.contextBlock.NewFSub(lhs, rhs)
	case "*":
		opResult = c.contextBlock.NewFMul(lhs, rhs)
	case "/":
		opResult = c.contextBlock.NewFDiv(lhs, rhs)
	case "==":
		opResult = c.contextBlock.NewFCmp(enum.FPredOEQ, lhs, rhs)
	case "!=":
		opResult = c.contextBlock.NewFCmp(enum.FPredONE, lhs, rhs)
	case "<":
		opResult = c.contextBlock.NewFCmp(enum.FPredOLT, lhs, rhs)
	case "<=":
		opResult = c.contextBlock.NewFCmp(enum.FPredOLE, lhs, rhs)
	case ">":
		opResult = c.contextBlock.NewFCmp(enum.FPredOGT, lhs, rhs)
	case ">=":
		opResult = c.contextBlock.NewFCmp(enum.FPredOGE, lhs, rhs)
	default:
		errors.ErrorExit(fmt.Sprintf("unexpected operator: float %s float\n", op))
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
func (c *CodeGen) genAssignExpression(expr *ast.AssignExpression) Value {
	lhs := c.genExpression(expr.Left).Value
	rhs := c.genExpression(expr.Value).Load(c.contextBlock)

	lhsTyp := internal.PtrElmType(lhs)
	rhsTyp := rhs.Type()

	if !lhsTyp.Equal(rhsTyp) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", expr.Token.Pos, lhsTyp, rhsTyp))
	}

	c.contextBlock.NewStore(rhs, lhs)

	return Value{
		Value:      rhs,
		IsVariable: true,
	}
}

func (c *CodeGen) genIndexExpression(expr *ast.IndexExpression) Value {
	left := c.genExpression(expr.Left).Value
	leftTyp := internal.PtrElmType(left)

	if _, ok := leftTyp.(*types.ArrayType); !ok {
		errors.ErrorExit(fmt.Sprintf("%s | cannot index %s", expr.Token.Pos, leftTyp))

	}

	index := c.genExpression(expr.Index).Load(c.contextBlock)
	val := c.contextBlock.NewGetElementPtr(leftTyp, left, constant.NewInt(types.I64, 0), index)
	val.InBounds = true
	return Value{
		Value:      val,
		IsVariable: true,
	}
}

func (c *CodeGen) genIntegerLiteral(expr *ast.IntegerLiteral) Value {
	return Value{
		Value:      constant.NewInt(types.I32, int64(expr.Value)),
		IsVariable: false,
	}
}

func (c *CodeGen) genFloatLiteral(expr *ast.FloatLiteral) Value {
	return Value{
		Value:      constant.NewFloat(types.Float, expr.Value),
		IsVariable: false,
	}
}

func (c *CodeGen) genIdentifier(expr *ast.Identifier) Value {
	v, ok := c.context.findVariable(expr)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("unresolved variable: %s\n", expr.String()))
	}

	return v
}

func (c *CodeGen) genNewExpression(expr *ast.NewExpression) Value {
	typ, ok := c.context.findType(expr.Ident.String())
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unknown type %s", expr.Token.Pos, expr.Ident))
	}

	return Value{
		Value:      c.contextBlock.NewAlloca(typ),
		IsVariable: true,
	}
}
func (c *CodeGen) genLoadMemberExpression(expr *ast.LoadMemberExpression) Value {
	lhs := c.genExpression(expr.Left).Value
	lhsTyp := internal.PtrElmType(lhs)

	structLlvmTyp, ok := lhsTyp.(*types.StructType)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unexpected operator: %s.%s\n", expr.Token.Pos, lhsTyp, expr.MemberIdent))

	}

	structTyp, ok := c.context.findStruct(structLlvmTyp.Name())
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unexpected operator: %s.%s\n", expr.Token.Pos, lhsTyp, expr.MemberIdent))

	}

	id := structTyp.findMember(expr.MemberIdent.String())
	if id == -1 {
		errors.ErrorExit(fmt.Sprintf("%s | unresolved member: %s\n", expr.Token.Pos, expr.String()))
	}

	member := structTyp.Members[id]

	zero := constant.NewInt(types.I32, 0)
	index := constant.NewInt(types.I32, int64(member.Id))
	val := c.contextBlock.NewGetElementPtr(lhsTyp, lhs, zero, index)

	return Value{
		Value:      val,
		IsVariable: true, // TODO
	}
}

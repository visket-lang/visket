package codegen

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler/ast"
	"github.com/arata-nvm/visket/compiler/codegen/builtin"
	"github.com/arata-nvm/visket/compiler/codegen/internal"
	"github.com/arata-nvm/visket/compiler/errors"
	"github.com/arata-nvm/visket/compiler/token"
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
	case *ast.StringLiteral:
		return c.genStringLiteral(expr)
	case *ast.CharLiteral:
		return c.genCharLiteral(expr)
	case *ast.Identifier:
		return c.genIdentifier(expr)
	case *ast.NewExpression:
		return c.genNewExpression(expr)
	case *ast.LoadMemberExpression:
		return c.genLoadMemberExpression(expr)
	}

	errors.ErrorExit(fmt.Sprintf("unexpexted expression: %s\n", ast.Show(expr)))
	return Value{} //unreachable
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) Value {
	lhs := c.genExpression(ie.Left).Load(c.contextBlock)
	rhs := c.genExpression(ie.Right).Load(c.contextBlock)

	lhsTyp := lhs.Type()
	rhsTyp := rhs.Type()

	if !lhsTyp.Equal(rhsTyp) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", ie.OpPos, lhsTyp, rhsTyp))
	}

	if lhsTyp.Equal(types.Float) {
		return c.genInfixFloat(ie.Op, lhs, rhs, ie.OpPos)
	}

	// TODO make default infix expr gen
	return c.genInfixInteger(ie.Op, lhs, rhs, ie.OpPos)
}

func (c *CodeGen) genInfixInteger(op string, lhs value.Value, rhs value.Value, pos token.Position) Value {
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
		errors.ErrorExit(fmt.Sprintf("%s | unexpected operator: %s %s %s", pos, lhs.Type(), op, rhs.Type()))
	}

	return Value{
		Value:      opResult,
		IsVariable: false,
	}
}

func (c *CodeGen) genInfixFloat(op string, lhs value.Value, rhs value.Value, pos token.Position) Value {
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
		errors.ErrorExit(fmt.Sprintf("%s | unexpected operator: %s %s %s", pos, lhs.Type(), op, rhs.Type()))
	}

	return Value{
		Value:      opResult,
		IsVariable: false,
	}
}

func (c *CodeGen) genCallExpression(expr *ast.CallExpression) Value {
	f, ok := c.context.findFunction(expr.Function.Name)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | undefined function '%s'", expr.LParen, expr.Function.Name))
	}

	if len(expr.Args) < len(f.Func.Params) {
		errors.ErrorExit(fmt.Sprintf("%s | not enough arguments in call to '%s'", expr.LParen, expr.Function.Name))
	} else if !f.Func.Sig.Variadic && len(expr.Args) > len(f.Func.Params) {
		errors.ErrorExit(fmt.Sprintf("%s | too many arguments in call to '%s'", expr.LParen, expr.Function.Name))
	}

	var params []value.Value

	for i, param := range expr.Args {
		// TODO rewrite
		// isReference
		exprVal := c.genExpression(param)
		var v value.Value
		if i < len(f.IsReference) && f.IsReference[i] {
			if !exprVal.IsVariable || exprVal.IsConstant {
				errors.ErrorExit(fmt.Sprintf("%s | a ref value must be an assignable variable", expr.RParen))
			}
			v = exprVal.Value
		} else {
			v = exprVal.Load(c.contextBlock)
		}
		params = append(params, v)
		if i >= len(f.Func.Sig.Params) {
			// variadic function
			continue
		}
		if !v.Type().Equal(f.Func.Sig.Params[i]) {
			errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", expr.LParen, v.Type(), f.Func.Sig.Params[i]))
		}
	}

	funcRet := c.contextBlock.NewCall(f.Func, params...)

	return Value{
		Value:      funcRet,
		IsVariable: false,
	}
}

func (c *CodeGen) genAssignExpression(expr *ast.AssignExpression) Value {
	left := c.genExpression(expr.Left)
	if left.IsConstant {
		errors.ErrorExit(fmt.Sprintf("%s | constant '%s' cannot be reassigned", expr.OpPos, ast.Show(expr.Left)))
	}
	lhs := left.Value

	rhs := c.genExpression(expr.Value).Load(c.contextBlock)

	lhsTyp := internal.PtrElmType(lhs)
	rhsTyp := rhs.Type()

	if !lhsTyp.Equal(rhsTyp) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", expr.OpPos, lhsTyp, rhsTyp))
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

	if _, ok := leftTyp.(*types.ArrayType); ok {
		return c.genArrayIndexing(left, leftTyp, expr)
	}

	if leftTyp.Equal(builtin.STRING) {
		return c.genStringIndexing(left, expr)
	}

	errors.ErrorExit(fmt.Sprintf("%s | cannot index '%s'", expr.LBrack, leftTyp))
	return Value{} // unreachable
}

func (c *CodeGen) genArrayIndexing(left value.Value, leftTyp types.Type, expr *ast.IndexExpression) Value {
	index := c.genExpression(expr.Index).Load(c.contextBlock)
	val := c.contextBlock.NewGetElementPtr(leftTyp, left, constant.NewInt(types.I64, 0), index)
	val.InBounds = true
	return Value{
		Value:      val,
		IsVariable: true,
	}
}

func (c *CodeGen) genStringIndexing(left value.Value, expr *ast.IndexExpression) Value {
	index := c.genExpression(expr.Index).Load(c.contextBlock)
	val := builtin.GetIndexedStringValue(left, index, c.contextBlock)
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

func (c *CodeGen) genStringLiteral(expr *ast.StringLiteral) Value {
	str := builtin.NewString(expr.Value, c.contextBlock, c.module)
	return Value{
		Value:      c.contextBlock.NewLoad(builtin.STRING, str),
		IsVariable: false,
	}
}

func (c *CodeGen) genCharLiteral(expr *ast.CharLiteral) Value {
	return Value{
		Value:      constant.NewInt(types.I8, int64(expr.Value)),
		IsVariable: false,
	}
}

func (c *CodeGen) genIdentifier(expr *ast.Identifier) Value {
	v, ok := c.context.findVariable(expr.Name)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unresolved variable '%s'", expr.Pos, expr.Name))
	}

	return v.Dereference(c.contextBlock)
}

func (c *CodeGen) genNewExpression(expr *ast.NewExpression) Value {
	typ := c.llvmType(expr.Type)
	val := c.contextBlock.NewAlloca(typ)
	initVal := constant.NewZeroInitializer(typ)
	c.contextBlock.NewStore(initVal, val)

	return Value{
		Value:      val,
		IsVariable: true,
	}
}

func (c *CodeGen) genLoadMemberExpression(expr *ast.LoadMemberExpression) Value {
	lhs := c.genExpression(expr.Left).Value
	lhsTyp := internal.PtrElmType(lhs)

	structLlvmTyp, ok := lhsTyp.(*types.StructType)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unexpected operator: %s.%s", expr.Period, lhsTyp, expr.MemberIdent.Name))
	}

	structTyp, ok := c.context.findStruct(structLlvmTyp.Name())
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unexpected operator: %s.%s", expr.Period, lhsTyp, expr.MemberIdent.Name))
	}

	if structTyp.IsIncomplete {
		errors.ErrorExit(fmt.Sprintf("%s | cannot load the member of incomplete structure: %s.%s", expr.Period, lhsTyp, expr.MemberIdent.Name))
	}

	id := structTyp.findMember(expr.MemberIdent.Name)
	if id == -1 {
		errors.ErrorExit(fmt.Sprintf("%s | unresolved member '%s'", expr.Period, expr.MemberIdent.Name))
	}

	member := structTyp.Members[id]

	zero := constant.NewInt(types.I32, 0)
	index := constant.NewInt(types.I32, int64(member.Id))
	val := c.contextBlock.NewGetElementPtr(lhsTyp, lhs, zero, index)

	return Value{
		Value:      val,
		IsVariable: true,
	}
}

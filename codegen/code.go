package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/codegen/types"
	"log"
	"strings"
)

func (c *CodeGen) gen(format string, a ...interface{}) {
	code := fmt.Sprintf(format, a...)
	_, err := fmt.Fprint(c.output, code)
	if err != nil {
		log.Fatal(err)
	}

	c.isTerminated = c.isTerminatorInst(code)
}

func (c *CodeGen) isTerminatorInst(code string) bool {
	rawCode := strings.TrimPrefix(code, "  ")
	inst := strings.Split(rawCode, " ")[0]
	switch inst {
	case "ret", "br":
		return true
	}
	return false
}

func (c *CodeGen) comment(format string, a ...interface{}) {
	if !c.isDebug {
		return
	}

	c.gen("")
	c.gen(format, a...)
}

func (c *CodeGen) genAlloca() Var {
	result := c.nextVar(types.I32Ptr)
	c.gen("  %s = alloca %s, align 4\n", result.Ident(), types.I32)
	return result
}

func (c *CodeGen) genNamedAlloca(v *Variable) {
	c.gen("  %s = alloca %s, align 4\n", v.Operand(), types.I32)
}

func (c *CodeGen) genStore(object Var, ptrToStore Var) {
	c.gen("  store %s, %s\n", object.Operand(), ptrToStore.Operand())
}

func (c *CodeGen) genNamedStore(v *Variable, ptrToStore Var) {
	c.gen("  store %s, i32* %s\n", ptrToStore.Operand(), v.Operand())
}

func (c *CodeGen) genStoreImmediate(value int, ptrToStore Var) {
	c.gen("  store %s %d, %s\n", types.I32, value, ptrToStore.Operand())
}

func (c *CodeGen) genLoad(ptrToLoad Var) Var {
	result := c.nextVar(types.I32)
	c.gen("  %s = load %s, %s, align 4\n", result.Ident(), types.I32, ptrToLoad.Operand())
	return result
}

func (c *CodeGen) genNamedLoad(v *Variable) Var {
	result := c.nextVar(types.I32)
	c.gen("  %s = load %s, i32* %s, align 4\n", result.Ident(), types.I32, v.Operand())
	return result
}

func (c *CodeGen) genAdd(op1 Var, op2 Var) Var {
	result := c.nextVar(types.I32)
	c.gen("  %s = add %s, %s\n", result.Ident(), op1.Operand(), op2.Ident())
	return result
}

func (c *CodeGen) genSub(op1 Var, op2 Var) Var {
	result := c.nextVar(types.I32)
	c.gen("  %s = sub %s, %s\n", result.Ident(), op1.Operand(), op2.Ident())
	return result
}

func (c *CodeGen) genMul(op1 Var, op2 Var) Var {
	result := c.nextVar(types.I32)
	c.gen("  %s = mul %s, %s\n", result.Ident(), op1.Operand(), op2.Ident())
	return result
}

func (c *CodeGen) genIDiv(op1 Var, op2 Var) Var {
	result := c.nextVar(types.I32)
	c.gen("  %s = idiv %s, %s\n", result.Ident(), op1.Operand(), op2.Ident())
	return result
}

type IcmpCond string

const (
	EQ  IcmpCond = "eq"
	NEQ          = "ne"
	LT           = "slt"
	LTE          = "sle"
	GT           = "sgt"
	GTE          = "sge"
)

func (c *CodeGen) genIcmp(cond IcmpCond, op1, op2 Var) Var {
	result := c.nextVar(types.I1)
	c.gen("  %s = icmp %s %s, %s\n", result.Ident(), cond, op1.Operand(), op2.Ident())
	return result
}

func (c *CodeGen) genIcmpWithNum(cond IcmpCond, op1 Var, op2 int) Var {
	result := c.nextVar(types.I1)
	c.gen("  %s = icmp %s %s, %d\n", result.Ident(), cond, op1.Operand(), op2)
	return result
}

func (c *CodeGen) genZext(typeTo types.Types, object Var) Var {
	result := c.nextVar(typeTo)
	c.gen("  %s = zext %s to %s\n", result.Ident(), object.Operand(), typeTo)
	return result
}

func (c *CodeGen) genRet(object Var) {
	c.gen("  ret %s\n", object.Operand())
}

func (c *CodeGen) genDefineFunction(ident *ast.Identifier) {
	c.gen("define %s @%s(", types.I32, ident.Token.Literal)
}

func (c *CodeGen) genFunctionParameters(params []*ast.Identifier) {
	var p []string
	for _, _ = range params {
		p = append(p, types.I32)
	}

	c.gen(strings.Join(p, ","))
}

func (c *CodeGen) genBeginFunction() {
	c.gen(") nounwind {\n")
}

func (c *CodeGen) genEndFunction() {
	c.gen("}\n\n")
}

func (c *CodeGen) genCall(function *ast.Identifier, params []Var) {
	var p []string
	for _, param := range params {
		p = append(p, param.Operand())
	}

	c.gen("  call %s @%s(%s)\n", types.I32, function.Token.Literal, strings.Join(p, ","))
}

func (c *CodeGen) genCallWithReturn(function *ast.Identifier, params []Var) Var {
	result := c.nextVar(types.I32)

	var p []string
	for _, param := range params {
		p = append(p, param.Operand())
	}

	c.gen("  %s = call %s @%s(%s)\n", result.Ident(), types.I32, function.Token.Literal, strings.Join(p, ","))
	return result
}

func (c *CodeGen) genLabel(name Label) {
	c.gen("%s:\n", name)
}

func (c *CodeGen) genBr(label Label) {
	c.gen("  br label %%%s\n", label)
}

func (c *CodeGen) genBrWithCond(condition Var, ifTrue Label, itFalse Label) {
	c.gen("  br %s, label %%%s, label %%%s\n", condition.Operand(), ifTrue, itFalse)
}

func (c *CodeGen) genTrunc(typeTo types.Types, object Var) Var {
	result := c.nextVar(typeTo)
	c.gen("  %s = trunc %s to %s\n", result.Ident(), object.Operand(), typeTo)
	return result
}

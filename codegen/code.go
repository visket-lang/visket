package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
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
	result := c.nextVar()
	c.gen("  %s = alloca i32, align 4\n", result.Operand())
	return result
}

func (c *CodeGen) genNamedAlloca(v *Variable) {
	c.gen("  %s = alloca i32, align 4\n", v.Operand())
}

func (c *CodeGen) genStore(object Var, ptrToStore Var) {
	c.gen("  store i32 %s, i32* %s\n", object.Operand(), ptrToStore.Operand())
}

func (c *CodeGen) genNamedStore(v *Variable, ptrToStore Var) {
	c.gen("  store i32 %s, i32* %s\n", ptrToStore.Operand(), v.Operand())
}

func (c *CodeGen) genStoreImmediate(value int, ptrToStore Var) {
	c.gen("  store i32 %d, i32* %s\n", value, ptrToStore.Operand())
}

func (c *CodeGen) genLoad(ptrToLoad Var) Var {
	result := c.nextVar()
	c.gen("  %s = load i32, i32* %s, align 4\n", result.Operand(), ptrToLoad.Operand())
	return result
}

func (c *CodeGen) genNamedLoad(v *Variable) Var {
	result := c.nextVar()
	c.gen("  %s = load i32, i32* %s, align 4\n", result.Operand(), v.Operand())
	return result
}

func (c *CodeGen) genAdd(op1 Var, op2 Var) Var {
	result := c.nextVar()
	c.gen("  %s = add i32 %s, %s\n", result.Operand(), op1.Operand(), op2.Operand())
	return result
}

func (c *CodeGen) genSub(op1 Var, op2 Var) Var {
	result := c.nextVar()
	c.gen("  %s = sub i32 %s, %s\n", result.Operand(), op1.Operand(), op2.Operand())
	return result
}

func (c *CodeGen) genMul(op1 Var, op2 Var) Var {
	result := c.nextVar()
	c.gen("  %s = mul i32 %s, %s\n", result.Operand(), op1.Operand(), op2.Operand())
	return result
}

func (c *CodeGen) genIDiv(op1 Var, op2 Var) Var {
	result := c.nextVar()
	c.gen("  %s = idiv i32 %s, %s\n", result.Operand(), op1.Operand(), op2.Operand())
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
	result := c.nextVar()
	c.gen("  %s = icmp %s i32 %s, %s\n", result.Operand(), cond, op1.Operand(), op2.Operand())
	return result
}

func (c *CodeGen) genIcmpWithNum(cond IcmpCond, op1 Var, op2 int) Var {
	result := c.nextVar()
	c.gen("  %s = icmp %s i32 %s, %d\n", result.Operand(), cond, op1.Operand(), op2)
	return result
}

func (c *CodeGen) genZext(typeFrom, typeTo string, object Var) Var {
	result := c.nextVar()
	c.gen("  %s = zext %s %s to %s\n", result.Operand(), typeFrom, object.Operand(), typeTo)
	return result
}

func (c *CodeGen) genRet(object Var) {
	c.gen("  ret i32 %s\n", object.Operand())
}

func (c *CodeGen) genDefineFunction(ident *ast.Identifier) {
	c.gen("define i32 @%s(", ident.Token.Literal)
}

func (c *CodeGen) genFunctionParameters(params []*ast.Identifier) {
	var p []string
	for _, _ = range params {
		p = append(p, "i32")
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
		p = append(p, fmt.Sprintf("i32 %s", param.Operand()))
	}

	c.gen("  call i32 @%s(%s)\n", function.Token.Literal, strings.Join(p, ","))
}

func (c *CodeGen) genCallWithReturn(function *ast.Identifier, params []Var) Var {
	result := c.nextVar()

	var p []string
	for _, param := range params {
		p = append(p, fmt.Sprintf("i32 %s", param.Operand()))
	}

	c.gen("  %s = call i32 @%s(%s)\n", result.Operand(), function.Token.Literal, strings.Join(p, ","))
	return result
}

func (c *CodeGen) genLabel(name Label) {
	c.gen("%s:\n", name)
}

func (c *CodeGen) genBr(label Label) {
	c.gen("  br label %%%s\n", label)
}

func (c *CodeGen) genBrWithCond(condition Var, ifTrue Label, itFalse Label) {
	c.gen("  br i1 %s, label %%%s, label %%%s\n", condition.Operand(), ifTrue, itFalse)
}

func (c *CodeGen) genTrunc(typeFrom, typeTo string, object Var) Var {
	result := c.nextVar()
	c.gen("  %s = trunc %s %s to %s\n", result.Operand(), typeFrom, object.Operand(), typeTo)
	return result
}

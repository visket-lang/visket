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
	inst := strings.Split(code, " ")[0]
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

func (c *CodeGen) indent() {
	c.gen("  ")
}

func (c *CodeGen) genAlloca(t types.Types) Value {
	result := c.nextReg(types.NewPointer(t))
	c.indent()
	c.gen("%s = alloca %s, align 4\n", result.RegName(), t.Name())
	return result
}

func (c *CodeGen) genNamedAlloca(v *Named) {
	c.indent()
	c.gen("%s = alloca %s, align 4\n", v.RegName(), v.TypeName())
	v.Type = types.NewPointer(v.Type)
}

func (c *CodeGen) genStore(src Value, dst Value) {
	c.indent()
	c.gen("store %s, %s\n", src.Operand(), dst.Operand())
}

func (c *CodeGen) genLoad(t types.Types, src Value) Value {
	result := c.nextReg(t)
	c.indent()
	c.gen("%s = load %s, %s %s, align 4\n", result.RegName(), t, types.NewPointer(t), src.RegName())
	return result
}

func (c *CodeGen) genAdd(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = add %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genSub(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = sub %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genMul(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = mul %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genSDiv(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = sdiv %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genSRem(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = srem %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genShl(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = shl %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genAShr(op1 Value, op2 Value) Value {
	result := c.nextReg(types.I32)
	c.indent()
	c.gen("%s = ashr %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
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

func (c *CodeGen) genIcmp(cond IcmpCond, op1, op2 Value) Value {
	result := c.nextReg(types.I1)
	c.indent()
	c.gen("%s = icmp %s %s, %s\n", result.RegName(), cond, op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genZext(typeTo types.Types, object Value) Value {
	result := c.nextReg(typeTo)
	c.indent()
	c.gen("%s = zext %s to %s\n", result.RegName(), object.Operand(), typeTo)
	return result
}

func (c *CodeGen) genTrunc(typeTo types.Types, object Value) Value {
	result := c.nextReg(typeTo)
	c.indent()
	c.gen("%s = trunc %s to %s\n", result.RegName(), object.Operand(), typeTo)
	return result
}

func (c *CodeGen) genRet(object Value) {
	c.indent()
	c.gen("ret %s\n", object.Operand())
}

func (c *CodeGen) genDefineFunction(ident *ast.Identifier) {
	c.gen("define %s @%s(", types.I32, ident.String())
}

func (c *CodeGen) genFunctionParameters(params []*ast.Identifier) {
	var p []string
	for range params {
		p = append(p, types.I32.String())
	}

	c.gen(strings.Join(p, ","))
}

func (c *CodeGen) genBeginFunction() {
	c.gen(") nounwind {\n")
}

func (c *CodeGen) genEndFunction() {
	c.gen("}\n\n")
}

func (c *CodeGen) genCall(function *ast.Identifier, params []Register) {
	var p []string
	for _, param := range params {
		p = append(p, param.Operand())
	}

	c.indent()
	c.gen("call %s @%s(%s)\n", types.I32, function.Token.Literal, strings.Join(p, ","))
}

func (c *CodeGen) genCallWithReturn(function *ast.Identifier, params []Value) Value {
	result := c.nextReg(types.I32)

	var p []string
	for _, param := range params {
		p = append(p, param.Operand())
	}

	c.indent()
	c.gen("%s = call %s @%s(%s)\n", result.RegName(), types.I32, function.Token.Literal, strings.Join(p, ","))
	return result
}

func (c *CodeGen) genLabel(name Label) {
	c.gen("%s:\n", name)
}

func (c *CodeGen) genBr(label Label) {
	c.indent()
	c.gen("br label %%%s\n", label)
}

func (c *CodeGen) genBrWithCond(condition Value, ifTrue Label, itFalse Label) {
	c.indent()
	c.gen("br %s, label %%%s, label %%%s\n", condition.Operand(), ifTrue, itFalse)
}

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

func (c *CodeGen) genAlloca(t types.Types) Value {
	result := c.nextVar(types.NewPointer(t))
	c.gen("  %s = alloca %s, align 4\n", result.RegName(), t.Name())
	return result
}

func (c *CodeGen) genNamedAlloca(v *Variable) {
	c.gen("  %s = alloca %s, align 4\n", v.RegName(), v.TypeName())
	v.Type = types.NewPointer(v.Type)
}

func (c *CodeGen) genStore(object Value, ptrToStore Value) {
	c.gen("  store %s, %s\n", object.Operand(), ptrToStore.Operand())
}

func (c *CodeGen) genStoreImmediate(value int, ptrToStore Value) {
	c.gen("  store %s %d, %s\n", types.I32, value, ptrToStore.Operand())
}

func (c *CodeGen) genLoad(t types.Types, ptrToLoad Value) Value {
	result := c.nextVar(t)
	c.gen("  %s = load %s, %s %s, align 4\n", result.RegName(), t, types.NewPointer(t), ptrToLoad.RegName())
	return result
}

func (c *CodeGen) genAdd(op1 Value, op2 Value) Value {
	result := c.nextVar(types.I32)
	c.gen("  %s = add %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genSub(op1 Value, op2 Value) Value {
	result := c.nextVar(types.I32)
	c.gen("  %s = sub %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genMul(op1 Value, op2 Value) Value {
	result := c.nextVar(types.I32)
	c.gen("  %s = mul %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genIDiv(op1 Value, op2 Value) Value {
	result := c.nextVar(types.I32)
	c.gen("  %s = idiv %s, %s\n", result.RegName(), op1.Operand(), op2.RegName())
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
	result := c.nextVar(types.I1)
	c.gen("  %s = icmp %s %s, %s\n", result.RegName(), cond, op1.Operand(), op2.RegName())
	return result
}

func (c *CodeGen) genIcmpWithNum(cond IcmpCond, op1 Value, op2 int) Value {
	result := c.nextVar(types.I1)
	c.gen("  %s = icmp %s %s, %d\n", result.RegName(), cond, op1.Operand(), op2)
	return result
}

func (c *CodeGen) genZext(typeTo types.Types, object Value) Value {
	result := c.nextVar(typeTo)
	c.gen("  %s = zext %s to %s\n", result.RegName(), object.Operand(), typeTo)
	return result
}

func (c *CodeGen) genRet(object Value) {
	c.gen("  ret %s\n", object.Operand())
}

func (c *CodeGen) genDefineFunction(ident *ast.Identifier) {
	c.gen("define %s @%s(", types.I32, ident.Token.Literal)
}

func (c *CodeGen) genFunctionParameters(params []*ast.Identifier) {
	var p []string
	for _, _ = range params {
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

func (c *CodeGen) genCall(function *ast.Identifier, params []Var) {
	var p []string
	for _, param := range params {
		p = append(p, param.Operand())
	}

	c.gen("  call %s @%s(%s)\n", types.I32, function.Token.Literal, strings.Join(p, ","))
}

func (c *CodeGen) genCallWithReturn(function *ast.Identifier, params []Value) Value {
	result := c.nextVar(types.I32)

	var p []string
	for _, param := range params {
		p = append(p, param.Operand())
	}

	c.gen("  %s = call %s @%s(%s)\n", result.RegName(), types.I32, function.Token.Literal, strings.Join(p, ","))
	return result
}

func (c *CodeGen) genLabel(name Label) {
	c.gen("%s:\n", name)
}

func (c *CodeGen) genBr(label Label) {
	c.gen("  br label %%%s\n", label)
}

func (c *CodeGen) genBrWithCond(condition Value, ifTrue Label, itFalse Label) {
	c.gen("  br %s, label %%%s, label %%%s\n", condition.Operand(), ifTrue, itFalse)
}

func (c *CodeGen) genTrunc(typeTo types.Types, object Value) Value {
	result := c.nextVar(typeTo)
	c.gen("  %s = trunc %s to %s\n", result.RegName(), object.Operand(), typeTo)
	return result
}

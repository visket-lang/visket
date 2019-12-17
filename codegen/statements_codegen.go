package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"os"
)

func (c *CodeGen) genStatement(stmt ast.Statement) {
	switch stmt := stmt.(type) {
	case *ast.VarStatement:
		c.genVarStatement(stmt)
	case *ast.ReturnStatement:
		c.genReturnStatement(stmt)
	case *ast.FunctionStatement:
		c.genFunctionStatement(stmt)
	case *ast.ExpressionStatement:
		c.genExpression(stmt.Expression)
	case *ast.IfStatement:
		c.genIfStatement(stmt)
	default:
		fmt.Printf("unexpexted statement: %s\n", stmt.Inspect())
		os.Exit(1)
	}
}

func (c *CodeGen) genVarStatement(stmt *ast.VarStatement) Value {
	c.comment("  ; Var\n")
	c.genNamedAlloca(stmt.Ident)
	resultPtr := c.genExpression(stmt.Value)
	// TODO Pointer への変換がよくわからない
	c.genNamedStore(stmt.Ident, Pointer(resultPtr))
	return c.genNamedLoad(stmt.Ident)
}

func (c *CodeGen) genReturnStatement(stmt *ast.ReturnStatement) {
	c.comment("  ; Ret\n")
	result := c.genExpression(stmt.Value)
	c.genRet(result)
}

func (c *CodeGen) genFunctionStatement(stmt *ast.FunctionStatement) {
	c.resetIndex()
	c.genDefineFunction(stmt.Ident)
	c.genFunctionParameters(stmt.Parameters)
	c.genBeginFunction()
	c.genLabel(c.nextLabel("entry"))

	for _, param := range stmt.Parameters {
		c.nextPointer()
		c.genNamedAlloca(param)
		c.genNamedStore(param, Pointer(c.index))
	}
	c.genBlockStatement(stmt.Body)
	c.genEndFunction()
}

func (c *CodeGen) genIfStatement(stmt *ast.IfStatement) {
	if stmt.Alternative != nil {
		c.genIfStatementWithElse(stmt)
		return
	}

	c.comment("  ; If\n")
	condition := c.genExpression(stmt.Condition)
	lThen := c.nextLabel("if.then")
	lMerge := c.nextLabel("if.merge")
	conditionI1 := c.genTrunc("i32", "i1", condition)
	c.genBrWithCond(conditionI1, lThen, lMerge)

	c.genLabel(lThen)
	c.genBlockStatement(stmt.Consequence)
	terminated := c.isTerminated
	if !terminated {
		c.genBr(lMerge)
	}

	c.genLabel(lMerge)
}

func (c *CodeGen) genIfStatementWithElse(stmt *ast.IfStatement) {
	c.comment("  ; If\n")
	condition := c.genExpression(stmt.Condition)
	lThen := c.nextLabel("if.then")
	lElse := c.nextLabel("if.else")
	lMerge := c.nextLabel("if.merge")
	conditionI1 := c.genTrunc("i32", "i1", condition)
	c.genBrWithCond(conditionI1, lThen, lElse)

	c.genLabel(lThen)
	c.genBlockStatement(stmt.Consequence)
	terminated := c.isTerminated
	if !c.isTerminated {
		c.genBr(lMerge)
	}

	c.genLabel(lElse)
	c.genBlockStatement(stmt.Alternative)
	terminated = terminated && c.isTerminated
	if !c.isTerminated {
		c.genBr(lMerge)
	}

	if !terminated {
		c.genLabel(lMerge)
	}
}

func (c *CodeGen) genBlockStatement(stmt *ast.BlockStatement) {
	for _, s := range stmt.Statements {
		c.genStatement(s)
	}
}
package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/codegen/types"
	"os"
)

func (c *CodeGen) genStatement(stmt ast.Statement) {
	switch stmt := stmt.(type) {
	case *ast.VarStatement:
		c.genVarStatement(stmt)
	case *ast.ReassignStatement:
		c.genReassignStatement(stmt)
	case *ast.ReturnStatement:
		c.genReturnStatement(stmt)
	case *ast.FunctionStatement:
		c.genFunctionStatement(stmt)
	case *ast.ExpressionStatement:
		c.genExpression(stmt.Expression)
	case *ast.IfStatement:
		c.genIfStatement(stmt)
	case *ast.WhileStatement:
		c.genWhileStatement(stmt)
	case *ast.ForStatement:
		c.genForStatement(stmt)
	default:
		fmt.Printf("unexpexted statement: %s\n", stmt.Inspect())
		os.Exit(1)
	}
}

func (c *CodeGen) genVarStatement(stmt *ast.VarStatement) {
	c.comment("  ; Var\n")

	_, ok := c.context.findVariable(stmt.Ident)
	if ok {
		fmt.Printf("already declared variable: %s\n", stmt.Ident.String())
		os.Exit(1)
	}

	v := c.context.newVariable(types.I32, stmt.Ident)
	c.genNamedAlloca(v)
	resultPtr := c.genExpression(stmt.Value)
	c.genStore(resultPtr, v)
}

func (c *CodeGen) genReassignStatement(stmt *ast.ReassignStatement) {
	c.comment("  ; Reassign\n")

	v, ok := c.context.findVariable(stmt.Ident)
	if !ok {
		fmt.Printf("unresolved variable: %s\n", stmt.Ident.String())
		os.Exit(1)
	}
	resultPtr := c.genExpression(stmt.Value)
	c.genStore(resultPtr, v)
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
	c.into()
	c.genLabel(c.nextLabel("entry"))

	for _, param := range stmt.Parameters {
		v := c.context.newVariable(types.I32, param)
		vv := c.nextVar(types.I32)
		c.genNamedAlloca(v)
		c.genStore(vv, v)

	}
	c.genBlockStatement(stmt.Body)

	c.outOf()
	c.genEndFunction()
}

func (c *CodeGen) genIfStatement(stmt *ast.IfStatement) {
	c.comment("  ; If\n")

	hasAlternative := stmt.Alternative != nil

	condition := c.genExpression(stmt.Condition)
	lThen := c.nextLabel("if.then")
	lElse := c.nextLabel("if.else")
	lMerge := c.nextLabel("if.merge")
	conditionI1 := c.genTrunc(types.I1, condition)
	if hasAlternative {
		c.genBrWithCond(conditionI1, lThen, lElse)
	} else {
		c.genBrWithCond(conditionI1, lThen, lMerge)
	}

	c.genLabel(lThen)
	c.into()
	c.genBlockStatement(stmt.Consequence)
	if !c.isTerminated {
		c.genBr(lMerge)
	}
	c.outOf()

	if hasAlternative {
		c.into()
		c.genLabel(lElse)
		c.genBlockStatement(stmt.Alternative)
		if !c.isTerminated {
			c.genBr(lMerge)
		}
		c.outOf()
	}

	c.genLabel(lMerge)
}

func (c *CodeGen) genWhileStatement(stmt *ast.WhileStatement) {
	c.comment("  ; While\n")
	lLoop := c.nextLabel("while.loop")
	lExit := c.nextLabel("while.exit")

	cond := c.genExpression(stmt.Condition)
	result := c.genIcmpWithNum(NEQ, cond, 0)
	c.genBrWithCond(result, lLoop, lExit)

	c.genLabel(lLoop)
	c.into()

	c.genBlockStatement(stmt.Body)

	cond = c.genExpression(stmt.Condition)
	result = c.genIcmpWithNum(NEQ, cond, 0)
	c.genBrWithCond(result, lLoop, lExit)

	c.outOf()
	c.genLabel(lExit)
}

func (c *CodeGen) genForStatement(stmt *ast.ForStatement) {
	c.comment("  ; For\n")
	lLoop := c.nextLabel("for.loop")
	lExit := c.nextLabel("for.exit")

	if stmt.Init != nil {
		c.genStatement(stmt.Init)
	}

	if stmt.Condition != nil {
		cond := c.genExpression(stmt.Condition)
		result := c.genIcmpWithNum(NEQ, cond, 0)
		c.genBrWithCond(result, lLoop, lExit)
	} else {
		c.genBr(lLoop)
	}

	c.genLabel(lLoop)
	c.into()

	c.genBlockStatement(stmt.Body)

	if stmt.Post != nil {
		c.genStatement(stmt.Post)
	}

	if stmt.Condition != nil {
		cond := c.genExpression(stmt.Condition)
		result := c.genIcmpWithNum(NEQ, cond, 0)
		c.genBrWithCond(result, lLoop, lExit)
	} else {
		c.genBr(lLoop)
	}

	c.outOf()
	c.genLabel(lExit)
}

func (c *CodeGen) genBlockStatement(stmt *ast.BlockStatement) {
	for _, s := range stmt.Statements {
		c.genStatement(s)
	}
}

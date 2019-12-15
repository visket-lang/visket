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

	for _, param := range stmt.Parameters {
		c.genNamedAlloca(param)
		c.genNamedStore(param, Pointer(c.index))
		c.nextPointer()
	}
	c.genBlockStatement(stmt.Body)
	c.genEndFunction()
}

func (c *CodeGen) genBlockStatement(stmt *ast.BlockStatement) {
	for _, s := range stmt.Statements {
		c.genStatement(s)
	}
}

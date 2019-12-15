package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"io"
	"os"
)

type CodeGen struct {
	program *ast.Program
	index   int
	isDebug bool
	output  io.Writer
}

func New(program *ast.Program, isDebug bool, w io.Writer) *CodeGen {
	c := &CodeGen{
		program: program,
		isDebug: isDebug,
		output:  w,
	}

	return c
}

func (c *CodeGen) GenerateCode() {
	for _, s := range c.program.Statements {
		c.genStatement(s)
	}
}

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
	c.genDefineFunction(stmt.Ident)
	if stmt.Parameter != nil {
		c.genFunctionParameter(stmt.Parameter)
	}
	c.genBeginFunction()

	if stmt.Parameter != nil {
		c.genNamedAlloca(stmt.Parameter)
		c.genNamedStore(stmt.Parameter, Pointer(c.index))
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

func (c *CodeGen) genExpression(expr ast.Expression) Value {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return c.genInfix(expr)
	case *ast.IntegerLiteral:
		c.comment("  ; Int\n")
		result := c.genAlloca()
		c.genStoreImmediate(expr.Value, result)
		return c.genLoad(result)
	case *ast.Identifier:
		c.comment("  ; Ident\n")
		return c.genNamedLoad(expr)
	}

	fmt.Printf("unexpexted expression: %s\n", expr.Inspect())
	os.Exit(1)
	return -1
}

func (c *CodeGen) genInfix(ie *ast.InfixExpression) Value {
	c.comment("  ; Infix\n")

	lhs := c.genExpression(ie.Left)
	rhs := c.genExpression(ie.Right)

	c.comment("  ; Op\n")

	switch ie.Operator {
	case "+":
		c.comment("  ; Add\n")
		return c.genAdd(lhs, rhs)
	case "-":
		c.comment("  ; Sub\n")
		return c.genSub(lhs, rhs)
	case "*":
		c.comment("  ; Mul\n")
		return c.genMul(lhs, rhs)
	case "/":
		c.comment("  ; Div\n")
		return c.genIDiv(lhs, rhs)
	case "==":
		c.comment("  ; Equal\n")
		result := c.genIcmp(EQ, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case "!=":
		c.comment("  ; Not Equal\n")
		result := c.genIcmp(NEQ, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case "<":
		c.comment("  ; Less Than\n")
		result := c.genIcmp(LT, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case "<=":
		c.comment("  ; Less Than or Equal\n")
		result := c.genIcmp(LTE, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case ">":
		c.comment("  ; Greater Than\n")
		result := c.genIcmp(GT, lhs, rhs)
		return c.genZext("i1", "i32", result)
	case ">=":
		c.comment("  ; Greater Than or Equal\n")
		result := c.genIcmp(GTE, lhs, rhs)
		return c.genZext("i1", "i32", result)
	}

	return Value(c.index)
}

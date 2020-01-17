package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/arata-nvm/Solitude/compiler/token"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (c *CodeGen) genStatement(stmt ast.Statement) {
	switch stmt := stmt.(type) {
	case *ast.VarStatement:
		c.genVarStatement(stmt)
	case *ast.AssignStatement:
		c.genAssignStatement(stmt)
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
		errors.ErrorExit(fmt.Sprintf("unexpexted statement: %s\n", stmt.Inspect()))
	}
}

func (c *CodeGen) genVarStatement(stmt *ast.VarStatement) {
	_, ok := c.context.findVariable(stmt.Ident)
	if ok {
		errors.ErrorExit(fmt.Sprintf("already declared variable: %s\n", stmt.Ident.String()))
	}

	value := c.genExpression(stmt.Value)
	named := c.contextBlock.NewAlloca(value.Type())
	named.SetName(stmt.Ident.String())
	c.context.addVariable(stmt.Ident, named)
	c.contextBlock.NewStore(value, named)
}

func (c *CodeGen) genAssignStatement(stmt *ast.AssignStatement) {
	v, ok := c.context.findVariable(stmt.Ident)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("unresolved variable: %s\n", stmt.Ident.String()))
	}

	rhs := c.genExpression(stmt.Value)

	switch stmt.Token.Type {
	case token.ASSIGN:
		c.contextBlock.NewStore(rhs, v)
	case token.ADD_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewAdd(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	case token.SUB_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewSub(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	case token.MUL_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewMul(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	case token.QUO_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewSDiv(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	case token.REM_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewSRem(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	case token.SHL_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewShl(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	case token.SHR_ASSIGN:
		vValue := c.contextBlock.NewLoad(types.I32, v)
		rhs = c.contextBlock.NewAShr(vValue, rhs)
		c.contextBlock.NewStore(rhs, v)
	}
}

func (c *CodeGen) genReturnStatement(stmt *ast.ReturnStatement) {
	result := c.genExpression(stmt.Value)
	c.contextBlock.NewRet(result)
}

func (c *CodeGen) genFunctionStatement(stmt *ast.FunctionStatement) {
	var params []*ir.Param

	for _, p := range stmt.Parameters {
		param := ir.NewParam(p.String(), types.I32)
		params = append(params, param)
		c.context.addVariable(p, param)
	}
	c.contextFunction = c.module.NewFunc(stmt.Ident.String(), types.I32, params...)
	c.context.addFunction(stmt.Ident, c.contextFunction)

	c.into()
	c.contextBlock = c.contextFunction.NewBlock("entry")

	c.genBlockStatement(stmt.Body)

	c.contextBlock = nil
	c.outOf()

	c.contextFunction = nil
}

func addLineNum(blockName string, tok token.Token) string {
	return fmt.Sprintf("%s.%d", blockName, tok.Pos.Line)
}

func (c *CodeGen) genIfStatement(stmt *ast.IfStatement) {
	hasAlternative := stmt.Alternative != nil

	condition := c.genExpression(stmt.Condition)
	blockThen := c.contextFunction.NewBlock(addLineNum("if.then", stmt.Token))
	var blockElse *ir.Block
	blockMerge := c.contextFunction.NewBlock(addLineNum("if.merge", stmt.Token))
	c.contextCondAfter = append(c.contextCondAfter, blockMerge)

	if hasAlternative {
		blockElse = c.contextFunction.NewBlock(addLineNum("if.else", stmt.Token))
		c.contextBlock.NewCondBr(condition, blockThen, blockElse)
	} else {
		c.contextBlock.NewCondBr(condition, blockThen, blockMerge)
	}

	c.into()
	c.contextBlock = blockThen
	c.contextBlock.NewBr(blockMerge)
	c.genBlockStatement(stmt.Consequence)
	c.outOf()

	if hasAlternative {
		c.into()
		c.contextBlock = blockElse
		c.contextBlock.NewBr(blockMerge)
		c.genBlockStatement(stmt.Alternative)
		c.outOf()
	}

	c.contextBlock = blockMerge
	c.contextCondAfter = c.contextCondAfter[:len(c.contextCondAfter)-1]

	if len(c.contextCondAfter) > 0 {
		c.contextBlock.NewBr(c.contextCondAfter[len(c.contextCondAfter)-1])
	}
}

func (c *CodeGen) genWhileStatement(stmt *ast.WhileStatement) {
	blockLoop := c.contextFunction.NewBlock(addLineNum("while.loop", stmt.Token))
	blockExit := c.contextFunction.NewBlock(addLineNum("while.exit", stmt.Token))

	cond := c.genExpression(stmt.Condition)
	result := c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
	c.contextBlock.NewCondBr(result, blockLoop, blockExit)

	c.into()
	c.contextBlock = blockLoop

	c.genBlockStatement(stmt.Body)

	cond = c.genExpression(stmt.Condition)
	result = c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
	c.contextBlock.NewCondBr(result, blockLoop, blockExit)
	c.outOf()

	c.contextBlock = blockExit
}

func (c *CodeGen) genForStatement(stmt *ast.ForStatement) {
	blockLoop := c.contextFunction.NewBlock(addLineNum("for.loop", stmt.Token))
	blockExit := c.contextFunction.NewBlock(addLineNum("for.exit", stmt.Token))

	if stmt.Init != nil {
		c.genStatement(stmt.Init)
	}

	if stmt.Condition != nil {
		cond := c.genExpression(stmt.Condition)
		result := c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
		c.contextBlock.NewCondBr(result, blockLoop, blockExit)
	} else {
		c.contextBlock.NewBr(blockLoop)
	}

	c.into()
	c.contextBlock = blockLoop

	c.genBlockStatement(stmt.Body)

	if stmt.Post != nil {
		c.genStatement(stmt.Post)
	}

	if stmt.Condition != nil {
		cond := c.genExpression(stmt.Condition)
		result := c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
		c.contextBlock.NewCondBr(result, blockLoop, blockExit)
	} else {
		c.contextBlock.NewBr(blockLoop)
	}

	c.outOf()

	c.contextBlock = blockExit
}

func (c *CodeGen) genBlockStatement(stmt *ast.BlockStatement) {
	for _, s := range stmt.Statements {
		c.genStatement(s)
	}
}

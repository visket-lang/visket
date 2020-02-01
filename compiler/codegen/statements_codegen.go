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

	// TODO 書き直す
	if stmt.Type != nil && stmt.Value == nil {
		typ := stmt.Type.LlvmType()
		named := c.contextBlock.NewAlloca(typ)
		named.SetName(stmt.Ident.String())
		c.context.addVariable(stmt.Ident, named)
	}

	if stmt.Type == nil && stmt.Value != nil {
		value := c.genExpression(stmt.Value)
		named := c.contextBlock.NewAlloca(value.Type())
		named.SetName(stmt.Ident.String())
		c.context.addVariable(stmt.Ident, named)
		c.contextBlock.NewStore(value, named)
	}

	if stmt.Type != nil && stmt.Value != nil {
		typ := stmt.Type.LlvmType()
		value := c.genExpression(stmt.Value)
		if typ != value.Type() {
			errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Token.Pos, typ, value.Type()))
		}
		named := c.contextBlock.NewAlloca(value.Type())
		named.SetName(stmt.Ident.String())
		c.context.addVariable(stmt.Ident, named)
		c.contextBlock.NewStore(value, named)
	}
}

func (c *CodeGen) genReturnStatement(stmt *ast.ReturnStatement) {
	retType := c.contextFunction.Sig.RetType

	if stmt.Value == nil {
		if retType != types.Void {
			errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Token.Pos, retType, types.Void))
		}
		c.contextBlock.NewRet(nil)
		return
	}

	result := c.genExpression(stmt.Value)

	if retType != result.Type() {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Token.Pos, retType, result.Type()))
	}

	c.contextBlock.NewRet(result)
}

func (c *CodeGen) genFunctionStatement(stmt *ast.FunctionStatement) {
	var params []*ir.Param

	for i, _ := range stmt.Parameters {
		typ := stmt.Type.Params[i].LlvmType()
		param := ir.NewParam("", typ)
		params = append(params, param)
	}

	returnTyp := stmt.Type.RetType.LlvmType()

	c.contextFunction = c.module.NewFunc(stmt.Ident.String(), returnTyp, params...)
	c.context.addFunction(stmt.Ident, c.contextFunction)

	c.into()
	c.contextBlock = c.contextFunction.NewBlock("entry")

	// 引数の再代入のために必要
	for i, p := range stmt.Parameters {
		typ := stmt.Type.Params[i].LlvmType()
		param := ir.NewParam("", typ)
		param.LocalID = int64(i)
		pp := c.contextBlock.NewAlloca(typ)
		c.contextBlock.NewStore(param, pp)
		c.context.addVariable(p, pp)
	}

	c.genBlockStatement(stmt.Body)

	if returnTyp == types.Void {
		c.contextBlock.NewRet(nil)
	}

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

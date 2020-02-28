package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	. "github.com/arata-nvm/Solitude/compiler/codegen/internal"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (c *CodeGen) genStatement(stmt ast.Statement) {
	switch stmt := stmt.(type) {
	case *ast.VarStatement:
		c.genVarStatement(stmt)
	case *ast.ReturnStatement:
		c.genReturnStatement(stmt)
	case *ast.ExpressionStatement:
		c.genExpression(stmt.Expression)
	case *ast.IfStatement:
		c.genIfStatement(stmt)
	case *ast.WhileStatement:
		c.genWhileStatement(stmt)
	case *ast.ForStatement:
		c.genForStatement(stmt)
	default:
		errors.ErrorExit(fmt.Sprintf("unexpexted statement: %s\n", ast.Show(stmt)))
	}
}

func (c *CodeGen) genVarStatement(stmt *ast.VarStatement) {
	_, ok := c.context.findVariable(stmt.Ident)
	if ok {
		errors.ErrorExit(fmt.Sprintf("already declared variable: %s\n", stmt.Ident.Token.Literal))
	}

	var typ types.Type
	var val value.Value
	if stmt.Value != nil {
		val = c.genExpression(stmt.Value).Load(c.contextBlock)
	} else {
		typ = c.llvmType(stmt.Type)
		val = constant.NewZeroInitializer(typ)
	}

	if stmt.Type == nil {
		typ = val.Type()
	}

	if !typ.Equal(val.Type()) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Token.Pos, typ, val.Type()))
	}

	named := c.contextEntryBlock.NewAlloca(val.Type())
	named.SetName(stmt.Ident.Token.Literal)
	c.context.addVariable(stmt.Ident, Value{
		Value:      named,
		IsVariable: true,
	})
	c.contextBlock.NewStore(val, named)
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

	result := c.genExpression(stmt.Value).Load(c.contextBlock)

	if retType != result.Type() {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Token.Pos, retType, result.Type()))
	}

	c.contextBlock.NewRet(result)
}

func (c *CodeGen) genFunctionDeclaration(stmt *ast.FunctionStatement) {
	_, ok := c.context.findFunction(stmt.Sig.Ident)
	if ok {
		errors.ErrorExit(fmt.Sprintf("%s | already declared function %s", stmt.Token.Pos, stmt.Sig.Ident.Token.Literal))
	}

	var params []*ir.Param

	for _, p := range stmt.Sig.Params {
		typ := c.llvmType(p.Type)
		param := ir.NewParam(p.Ident.Token.Literal, typ)
		params = append(params, param)
	}

	returnTyp := c.llvmType(stmt.Sig.RetType)

	function := c.module.NewFunc(stmt.Sig.Ident.Token.Literal, returnTyp, params...)
	c.context.addFunction(stmt.Sig.Ident, function)
}

func (c *CodeGen) genFunctionBody(stmt *ast.FunctionStatement) {
	f, ok := c.context.findFunction(stmt.Sig.Ident)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | undeclared function %s", stmt.Token.Pos, stmt.Sig.Ident.Token.Literal))
	}

	c.contextFunction = f

	c.into()
	c.contextBlock = c.contextFunction.NewBlock("entry")
	c.contextEntryBlock = c.contextBlock

	for i, p := range stmt.Sig.Params {
		c.context.addVariable(p.Ident, Value{
			Value:      f.Params[i],
			IsVariable: false,
		})
	}

	c.genBlockStatement(stmt.Body)

	if f.Sig.RetType == types.Void {
		c.contextBlock.NewRet(nil)
	}

	if c.contextBlock.Term == nil {
		errors.ErrorExit(fmt.Sprintf("%s | missing return at end of function", stmt.Token.Pos))
	}

	c.contextEntryBlock = nil
	c.contextBlock = nil
	c.outOf()

	c.contextFunction = nil
}

func (c *CodeGen) genIfStatement(stmt *ast.IfStatement) {
	hasAlternative := stmt.Alternative != nil

	condition := c.genExpression(stmt.Condition).Load(c.contextBlock)
	blockThen := c.contextFunction.NewBlock(NextLabel("if.then"))
	var blockElse *ir.Block
	blockMerge := c.contextFunction.NewBlock(NextLabel("if.merge"))
	c.contextCondAfter = append(c.contextCondAfter, blockMerge)

	if hasAlternative {
		blockElse = c.contextFunction.NewBlock(NextLabel("if.else"))
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
	blockLoop := c.contextFunction.NewBlock(NextLabel("while.loop"))
	blockExit := c.contextFunction.NewBlock(NextLabel("while.exit"))

	cond := c.genExpression(stmt.Condition).Load(c.contextBlock)
	result := c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
	c.contextBlock.NewCondBr(result, blockLoop, blockExit)

	c.into()
	c.contextBlock = blockLoop

	c.genBlockStatement(stmt.Body)

	cond = c.genExpression(stmt.Condition).Load(c.contextBlock)
	result = c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
	c.contextBlock.NewCondBr(result, blockLoop, blockExit)
	c.outOf()

	c.contextBlock = blockExit
}

func (c *CodeGen) genForStatement(stmt *ast.ForStatement) {
	blockLoop := c.contextFunction.NewBlock(NextLabel("for.loop"))
	blockExit := c.contextFunction.NewBlock(NextLabel("for.exit"))

	if stmt.Init != nil {
		c.genStatement(stmt.Init)
	}

	if stmt.Condition != nil {
		cond := c.genExpression(stmt.Condition).Load(c.contextBlock)
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
		cond := c.genExpression(stmt.Condition).Load(c.contextBlock)
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

func (c *CodeGen) genStructStatement(stmt *ast.StructStatement) {
	s := &Struct{
		Name: stmt.Ident.Token.Literal,
	}

	var llvmMembers []types.Type
	for i, m := range stmt.Members {
		typ := c.llvmType(m.Type)
		s.Members = append(s.Members, &Member{
			Name: m.Ident.Token.Literal,
			Id:   i,
			Type: typ,
		})

		llvmMembers = append(llvmMembers, typ)
	}

	s.Type = types.NewStruct(llvmMembers...)

	c.module.NewTypeDef(s.Name, s.Type)
	c.context.addStruct(s.Name, s)
}

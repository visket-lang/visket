package codegen

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler/ast"
	. "github.com/arata-nvm/visket/compiler/codegen/internal"
	"github.com/arata-nvm/visket/compiler/errors"
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
	case *ast.ForRangeStatement:
		c.genForRangeStatement(stmt)
	default:
		errors.ErrorExit(fmt.Sprintf("unexpexted statement: %s\n", ast.Show(stmt)))
	}
}

func (c *CodeGen) genVarStatement(stmt *ast.VarStatement) {
	_, ok := c.context.findVariable(stmt.Ident.Name)
	if ok {
		errors.ErrorExit(fmt.Sprintf("%s | already declared variable '%s'", stmt.Var, stmt.Ident.Name))
	}

	var typ types.Type
	var val value.Value
	if stmt.Value != nil {
		val = c.genExpression(stmt.Value).Load(c.contextBlock)
		typ = val.Type()
	} else {
		typ = c.llvmType(stmt.Type)
		val = constant.NewZeroInitializer(typ)
	}

	if stmt.Type != nil {
		typ = c.llvmType(stmt.Type)
	}

	if !typ.Equal(val.Type()) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Var, typ, val.Type()))
	}

	named := c.contextEntryBlock.NewAlloca(val.Type())
	named.SetName(stmt.Ident.Name)
	c.context.addVariable(stmt.Ident.Name, Value{
		Value:      named,
		IsVariable: true,
	})
	c.contextBlock.NewStore(val, named)
}

func (c *CodeGen) genReturnStatement(stmt *ast.ReturnStatement) {
	retType := c.contextFunction.Sig.RetType

	if stmt.Value == nil {
		if c.contextFunction.Name() == "main" {
			c.contextBlock.NewRet(constant.NewInt(types.I32, 0))
			return
		}

		if retType != types.Void {
			errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Return, retType, types.Void))
		}
		c.contextBlock.NewRet(nil)
		return
	}

	result := c.genExpression(stmt.Value).Load(c.contextBlock)

	if retType != result.Type() {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.Return, retType, result.Type()))
	}

	c.contextBlock.NewRet(result)
}

func (c *CodeGen) genFunctionDeclaration(stmt *ast.FunctionStatement) {
	if stmt.Ident.Name == "main" {
		return
	}

	_, ok := c.context.findFunction(stmt.Ident.Name)
	if ok {
		errors.ErrorExit(fmt.Sprintf("%s | already declared function '%s'", stmt.Func, stmt.Ident.Name))
	}

	var params []*ir.Param

	for _, p := range stmt.Sig.Params {
		typ := c.llvmType(p.Type)
		if p.IsReference {
			typ = types.NewPointer(typ)
		}
		param := ir.NewParam("", typ)
		params = append(params, param)
	}

	returnTyp := c.llvmType(stmt.Sig.RetType)

	function := c.module.NewFunc(stmt.Ident.Name, returnTyp, params...)
	c.context.addFunction(stmt.Ident.Name, function)
}

func (c *CodeGen) genFunctionBody(stmt *ast.FunctionStatement) {
	f, ok := c.context.findFunction(stmt.Ident.Name)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | undeclared function '%s'", stmt.Func, stmt.Ident.Name))
	}

	c.contextFunction = f

	c.into()
	if stmt.Ident.Name == "main" {
		retTyp := stmt.Sig.RetType.Name
		if retTyp != "void" {
			errors.ErrorExit(fmt.Sprintf("%s | main func cannot have a return type", stmt.Func))
		}

		if len(stmt.Sig.Params) != 0 {
			errors.ErrorExit(fmt.Sprintf("%s | main func cannot have parameters", stmt.Func))
		}

		c.contextBlock = c.contextFunction.Blocks[0]
	} else {
		c.contextBlock = c.contextFunction.NewBlock("entry")
	}
	c.contextEntryBlock = c.contextBlock

	for i, p := range stmt.Sig.Params {
		typ := f.Params[i].Typ
		val := c.contextBlock.NewAlloca(typ)
		val.SetName(p.Ident.Name)
		c.contextBlock.NewStore(f.Params[i], val)
		c.context.addVariable(p.Ident.Name, Value{
			Value:       val,
			IsVariable:  true,
			IsReference: p.IsReference,
		})
	}

	c.genBlockStatement(stmt.Body)

	if f.Sig.RetType == types.Void {
		c.contextBlock.NewRet(nil)
	}

	if stmt.Ident.Name == "main" {
		c.contextBlock.NewRet(constant.NewInt(types.I32, 0))
	}

	if c.contextBlock.Term == nil {
		errors.ErrorExit(fmt.Sprintf("%s | missing return at end of function", stmt.Body.RBrace))
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

// TODO rewrite
func (c *CodeGen) genForRangeStatement(stmt *ast.ForRangeStatement) {
	blockLoop := c.contextFunction.NewBlock(NextLabel("for.loop"))
	blockExit := c.contextFunction.NewBlock(NextLabel("for.exit"))

	c.into()
	from := c.genExpression(stmt.From).Load(c.contextBlock)
	to := c.genExpression(stmt.To).Load(c.contextBlock)

	if !from.Type().Equal(to.Type()) {
		errors.ErrorExit(fmt.Sprintf("%s | type mismatch '%s' and '%s'", stmt.For, from.Type(), to.Type()))
	}

	typ := from.Type()
	namedVar := c.contextEntryBlock.NewAlloca(typ)
	namedVar.SetName(NextForNum(stmt.VarName.Name))
	c.contextBlock.NewStore(from, namedVar)
	c.context.addVariable("i", Value{
		Value:      namedVar,
		IsVariable: true,
	})

	val := c.contextBlock.NewLoad(typ, namedVar)
	cond := c.contextBlock.NewICmp(enum.IPredULE, val, to)
	result := c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
	c.contextBlock.NewCondBr(result, blockLoop, blockExit)

	c.contextBlock = blockLoop
	c.genBlockStatement(stmt.Body)

	val = c.contextBlock.NewLoad(typ, namedVar)
	nextVal := c.contextBlock.NewAdd(val, constant.NewInt(types.I32, 1))
	c.contextBlock.NewStore(nextVal, namedVar)

	val = c.contextBlock.NewLoad(typ, namedVar)
	cond = c.contextBlock.NewICmp(enum.IPredULE, val, to)
	result = c.contextBlock.NewICmp(enum.IPredNE, cond, constant.False)
	c.contextBlock.NewCondBr(result, blockLoop, blockExit)

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
		Name: stmt.Ident.Name,
	}

	var llvmMembers []types.Type
	for i, m := range stmt.Members {
		typ := c.llvmType(m.Type)
		s.Members = append(s.Members, &Member{
			Name: m.Ident.Name,
			Id:   i,
			Type: typ,
		})

		llvmMembers = append(llvmMembers, typ)
	}

	s.Type = types.NewStruct(llvmMembers...)

	c.module.NewTypeDef(s.Name, s.Type)
	c.context.addStruct(s.Name, s)
}

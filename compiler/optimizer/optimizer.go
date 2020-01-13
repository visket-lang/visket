package optimizer

import "github.com/arata-nvm/Solitude/compiler/ast"

type Optimizer struct {
	Program *ast.Program
}

func New(program *ast.Program) *Optimizer {
	return &Optimizer{Program: program}
}

func (o *Optimizer) Optimize() {
	var s []ast.Statement
	for _, stmt := range o.Program.Statements {
		s = append(s, o.optStatement(stmt))
	}

	o.Program.Statements = s
}

func (o *Optimizer) optBlockStatement(stmt *ast.BlockStatement) *ast.BlockStatement {
	if stmt == nil {
		return stmt
	}

	for i := range stmt.Statements {
		stmt.Statements[i] = o.optStatement(stmt.Statements[i])
	}
	return stmt
}

func (o *Optimizer) optStatement(stmt ast.Statement) ast.Statement {
	switch stmt := stmt.(type) {
	case *ast.BlockStatement:
		return o.optBlockStatement(stmt)
	case *ast.ExpressionStatement:
		stmt.Expression = o.optExpression(stmt.Expression)
	case *ast.VarStatement:
		stmt.Value = o.optExpression(stmt.Value)
	case *ast.AssignStatement:
		stmt.Value = o.optExpression(stmt.Value)
	case *ast.ReturnStatement:
		stmt.Value = o.optExpression(stmt.Value)
	case *ast.FunctionStatement:
		stmt.Body = o.optBlockStatement(stmt.Body)
	case *ast.IfStatement:
		stmt.Condition = o.optExpression(stmt.Condition)
		stmt.Consequence = o.optBlockStatement(stmt.Consequence)
		stmt.Alternative = o.optBlockStatement(stmt.Alternative)
	case *ast.WhileStatement:
		stmt.Condition = o.optExpression(stmt.Condition)
		stmt.Body = o.optBlockStatement(stmt.Body)
	case *ast.ForStatement:
		stmt.Init = o.optStatement(stmt.Init)
		stmt.Condition = o.optExpression(stmt.Condition)
		stmt.Post = o.optStatement(stmt.Post)
		stmt.Body = o.optBlockStatement(stmt.Body)
	}

	return stmt
}

func (o *Optimizer) optExpression(expr ast.Expression) ast.Expression {
	switch expr := expr.(type) {
	case *ast.InfixExpression:
		return o.optInfixExpression(expr)
	}

	return expr
}

func (o *Optimizer) optInfixExpression(expr *ast.InfixExpression) ast.Expression {
	expr.Left = o.optExpression(expr.Left)
	expr.Right = o.optExpression(expr.Right)

	lil, ok := expr.Left.(*ast.IntegerLiteral)
	if !ok {
		return expr
	}

	ril, ok := expr.Right.(*ast.IntegerLiteral)
	if !ok {
		return expr
	}

	switch expr.Operator {
	case "+":
		return &ast.IntegerLiteral{Value: lil.Value + ril.Value}
	case "-":
		return &ast.IntegerLiteral{Value: lil.Value - ril.Value}
	case "*":
		return &ast.IntegerLiteral{Value: lil.Value * ril.Value}
	case "/":
		return &ast.IntegerLiteral{Value: lil.Value / ril.Value}
	}

	return expr
}

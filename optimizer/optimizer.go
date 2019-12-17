package optimizer

import "github.com/arata-nvm/Solitude/ast"

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

func (o *Optimizer) optStatement(stmt ast.Statement) ast.Statement {
	switch stmt := stmt.(type) {
	case *ast.ExpressionStatement:
		stmt.Expression = o.optExpression(stmt.Expression)
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

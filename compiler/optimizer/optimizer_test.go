package optimizer

import (
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/token"
	"testing"
)

func TestOptimize(t *testing.T) {
	program := &ast.Program{
		Functions: []*ast.FunctionStatement{{
			Sig: &ast.FunctionSignature{
				Ident:   &ast.Identifier{Token: token.Token{Literal: "main"}},
				Params:  make([]*ast.Param, 0),
				RetType: &ast.Type{Token: token.Token{Literal: "void"}},
			},
			Body: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left:     &ast.IntegerLiteral{Value: 2},
								Operator: "*",
								Right:    &ast.IntegerLiteral{Value: 3},
							},
							Operator: "*",
							Right:    &ast.Identifier{Token: token.Token{Literal: "x"}},
						},
					},
				},
			},
		}},
	}

	expected := "(def-func main(): void ((6 * x)))"

	o := New(program)
	o.Optimize()
	if ast.Show(program) != expected {
		t.Fatalf("expected=%q, got=%q", expected, ast.Show(program))
	}
}

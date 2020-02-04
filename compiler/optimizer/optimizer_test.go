package optimizer

import (
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/token"
	"testing"
)

func TestOptimize(t *testing.T) {
	program := &ast.Program{
		Functions: []*ast.FunctionStatement{{
			Ident: &ast.Identifier{Token: token.Token{Literal: "main"}},
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

	expected := "func Ident(main)() {Infix(Int(6) * Ident(x))}"

	o := New(program)
	o.Optimize()
	if program.Inspect() != expected {
		t.Fatalf("expected=%q, got=%q", expected, program.Inspect())
	}
}

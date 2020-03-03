package optimizer

import (
	"github.com/arata-nvm/visket/compiler/ast"
	"testing"
)

func TestOptimize(t *testing.T) {
	program := &ast.Program{
		Functions: []*ast.FunctionStatement{{
			Ident: &ast.Identifier{
				Name: "main",
			},
			Sig: &ast.FunctionSignature{
				Params: make([]*ast.Param, 0),
				RetType: &ast.Type{
					Name: "void",
				},
			},
			Body: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.InfixExpression{
							Left: &ast.InfixExpression{
								Left: &ast.IntegerLiteral{
									Value: 2,
								},
								Op: "*",
								Right: &ast.IntegerLiteral{
									Value: 3,
								},
							},
							Op: "*",
							Right: &ast.Identifier{
								Name: "x",
							},
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

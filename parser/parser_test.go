package parser

import (
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/lexer"
	"testing"
)

func TestIntegerLiteral(t *testing.T) {
	input := `42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	es, ok := program.Code.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("code is not ast.ExpressionStatement. got=%T", program.Code)
	}

	il, ok := es.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("code is not ast.IntegerLiteral. got=%T", program.Code)
	}

	if il.Value != 42 {
		t.Errorf("il.Value is not %d. got=%d", 42, il.Value)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"4 + 4", "4", "+", "4"},
		{"4 - 4", "4", "-", "4"},
		{"4 * 4", "4", "*", "4"},
		{"4 / 4", "4", "/", "4"},
		{"4 == 4", "4", "==", "4"},
		{"4 != 4", "4", "!=", "4"},
		{"4 < 4", "4", "<", "4"},
		{"4 <= 4", "4", "<=", "4"},
		{"4 > 4", "4", ">", "4"},
		{"4 >= 4", "4", ">=", "4"},
	}

	for i, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		es, ok := program.Code.(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("tests[%d] - code is not ast.ExpressionStatement. got=%T", i, program.Code)
		}

		ie, ok := es.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("tests[%d] - code is not ast.InfixExpression. got=%T", i, program.Code)
		}

		if ie.Left.String() != test.leftValue {
			t.Fatalf("tests[%d] - ie.Left wrong. expected=%q, got=%q", i, test.leftValue, ie.Left)
		}

		if ie.Operator != test.operator {
			t.Fatalf("tests[%d] - ie.Operator wrong. expected=%q, got=%q", i, test.operator, ie.Operator)
		}

		if ie.Right.String() != test.rightValue {
			t.Fatalf("tests[%d] - ie.Right wrong. expected=%q, got=%q", i, test.rightValue, ie.Right)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.Errors))

	for _, msg := range p.Errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

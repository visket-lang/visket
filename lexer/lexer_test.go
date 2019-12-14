package lexer

import (
	"github.com/arata-nvm/Solitude/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
10 + 2
8 - 4
42 * 89
32 / 4
((10 + 10) * 2)
57 == 72
43 != 83
32 < 33 <= 33
59 > 58 >= 58

var a = 1
return a
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.PLUS, "+"},
		{token.INT, "2"},

		{token.INT, "8"},
		{token.MINUS, "-"},
		{token.INT, "4"},

		{token.INT, "42"},
		{token.ASTERISK, "*"},
		{token.INT, "89"},

		{token.INT, "32"},
		{token.SLASH, "/"},
		{token.INT, "4"},

		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.INT, "10"},
		{token.PLUS, "+"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.ASTERISK, "*"},
		{token.INT, "2"},
		{token.RPAREN, ")"},

		{token.INT, "57"},
		{token.EQ, "=="},
		{token.INT, "72"},

		{token.INT, "43"},
		{token.NEQ, "!="},
		{token.INT, "83"},

		{token.INT, "32"},
		{token.LT, "<"},
		{token.INT, "33"},
		{token.LTE, "<="},
		{token.INT, "33"},

		{token.INT, "59"},
		{token.GT, ">"},
		{token.INT, "58"},
		{token.GTE, ">="},
		{token.INT, "58"},

		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},

		{token.EOF, ""},
	}

	lexer := New(input)

	for i, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, test.expectedLiteral, token.Literal)
		}
	}
}

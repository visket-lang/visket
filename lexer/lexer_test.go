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

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
54 % 3
((10 + 10) * 2)
57 == 72
43 != 83
32 < 33 <= 33
59 > 58 >= 58

a += 1
b -= 2
c *= 3
d /= 4
e %= 5

var a = 1
return a
{ 1 }
func f(arg) {}
if a { return 1 }
else { return 0 }
while 1 { 1 }
for var i = 0; i < 10; i=i+1 { 1 }
// while 1 { 1 }
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.ADD, "+"},
		{token.INT, "2"},

		{token.INT, "8"},
		{token.SUB, "-"},
		{token.INT, "4"},

		{token.INT, "42"},
		{token.MUL, "*"},
		{token.INT, "89"},

		{token.INT, "32"},
		{token.QUO, "/"},
		{token.INT, "4"},

		{token.INT, "54"},
		{token.REM, "%"},
		{token.INT, "3"},

		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.INT, "10"},
		{token.ADD, "+"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.MUL, "*"},
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

		{token.IDENT, "a"},
		{token.ADD_ASSIGN, "+="},
		{token.INT, "1"},

		{token.IDENT, "b"},
		{token.SUB_ASSIGN, "-="},
		{token.INT, "2"},

		{token.IDENT, "c"},
		{token.MUL_ASSIGN, "*="},
		{token.INT, "3"},

		{token.IDENT, "d"},
		{token.QUO_ASSIGN, "/="},
		{token.INT, "4"},

		{token.IDENT, "e"},
		{token.REM_ASSIGN, "%="},
		{token.INT, "5"},

		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},

		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.RBRACE, "}"},

		{token.FUNCTION, "func"},
		{token.IDENT, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "arg"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},

		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "1"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "0"},
		{token.RBRACE, "}"},

		{token.WHILE, "while"},
		{token.INT, "1"},
		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.RBRACE, "}"},

		{token.FOR, "for"},
		{token.VAR, "var"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.IDENT, "i"},
		{token.ADD, "+"},
		{token.INT, "1"},
		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.RBRACE, "}"},

		{token.COMMENT, "// while 1 { 1 }"},

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

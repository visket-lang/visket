package lexer

import (
	"github.com/arata-nvm/visket/compiler/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
10
10.945
"hoge"
10 + 2
8 - 4
42 * 89
32 / 4
54 % 3
2 << 3
16 >> 3
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
a <<= 2
b >>= 3

var a = 1
return a
{ 1 }
func f(arg) {}
if a { return 1 }
else { return 0 }
while 1 { 1 }
for var i = 0; i < 10; i=i+1 { 1 }
for i in 0..10 { 1 }
// while 1 { 1 }
[1, 2, 3]
array[1]
struct Foo {
  X: int
  Y: float
}
new Foo
bar.X
"\a\b\f\n\r\t\v\"\\"
import "std"
1.upto(10)
func f(ref i: int) {}
f(ref i)
val i = 10
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.FLOAT, "10.945"},
		{token.STRING, "hoge"},

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

		{token.INT, "2"},
		{token.SHL, "<<"},
		{token.INT, "3"},

		{token.INT, "16"},
		{token.SHR, ">>"},
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

		{token.IDENT, "a"},
		{token.SHL_ASSIGN, "<<="},
		{token.INT, "2"},

		{token.IDENT, "b"},
		{token.SHR_ASSIGN, ">>="},
		{token.INT, "3"},

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

		{token.FOR, "for"},
		{token.IDENT, "i"},
		{token.IN, "in"},
		{token.INT, "0"},
		{token.RANGE, ".."},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.RBRACE, "}"},

		{token.COMMENT, "// while 1 { 1 }"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},

		{token.IDENT, "array"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.RBRACKET, "]"},

		{token.STRUCT, "struct"},
		{token.IDENT, "Foo"},
		{token.LBRACE, "{"},
		{token.IDENT, "X"},
		{token.COLON, ":"},
		{token.IDENT, "int"},
		{token.IDENT, "Y"},
		{token.COLON, ":"},
		{token.IDENT, "float"},
		{token.RBRACE, "}"},

		{token.NEW, "new"},
		{token.IDENT, "Foo"},

		{token.IDENT, "bar"},
		{token.PERIOD, "."},
		{token.IDENT, "X"},

		{token.STRING, "\a\b\f\n\r\t\v\"\\"},

		{token.IMPORT, "import"},
		{token.STRING, "std"},

		{token.INT, "1"},
		{token.PERIOD, "."},
		{token.IDENT, "upto"},
		{token.LPAREN, "("},
		{token.INT, "10"},
		{token.RPAREN, ")"},

		{token.FUNCTION, "func"},
		{token.IDENT, "f"},
		{token.LPAREN, "("},
		{token.REF, "ref"},
		{token.IDENT, "i"},
		{token.COLON, ":"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.IDENT, "f"},
		{token.LPAREN, "("},
		{token.REF, "ref"},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},

		{token.VAL, "val"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "10"},

		{token.EOF, ""},
	}

	lexer := NewFromString(input)

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

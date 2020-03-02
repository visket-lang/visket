package parser

import (
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/lexer"
	"testing"
)

func TestParseProgram(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"func f(a: int): int {return 1}", "(def-func f(a: int): int ((return 1)))"},
		{"func hoge(fuga: int): int {return fuga}", "(def-func hoge(fuga: int): int ((return fuga)))"},

		{"func num(): int {return 2} func main(): int {return num()}", "(def-func num(): int ((return 2)))(def-func main(): int ((return (func-call num()))))"},
		{"func add(n: int): int {return n + 2} func main(): int {return num(1)}", "(def-func add(n: int): int ((return (n + 2))))(def-func main(): int ((return (func-call num(1)))))"},
		{"func add(a: int, b: int): int {return a + b} func main(): int {return num(1, 2)}", "(def-func add(a: int, b: int): int ((return (a + b))))(def-func main(): int ((return (func-call num(1, 2)))))"},
		{"struct Foo { X: int Y: float }", "(struct Foo(X: int, Y: float))"},
	}

	for i, test := range tests {
		l := lexer.NewFromString(test.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := ast.Show(program)

		if actual != test.expected {
			t.Fatalf("tests[%d] - expected=%q, got=%q", i, test.expected, actual)
		}
	}
}

func TestParseStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"42", "42"},
		{"32.0", "32.000000"},
		{"\"hoge\"", "\"hoge\""},

		{"4 + 4", "(4 + 4)"},
		{"4 - 4", "(4 - 4)"},
		{"4 * 4", "(4 * 4)"},
		{"4 % 4", "(4 % 4)"},
		{"4 / 4", "(4 / 4)"},
		{"4 << 4", "(4 << 4)"},
		{"4 >> 4", "(4 >> 4)"},
		{"4 == 4", "(4 == 4)"},
		{"4 != 4", "(4 != 4)"},
		{"4 < 4", "(4 < 4)"},
		{"4 <= 4", "(4 <= 4)"},
		{"4 > 4", "(4 > 4)"},
		{"4 >= 4", "(4 >= 4)"},

		{"4 + 4 * 4", "(4 + (4 * 4))"},
		{"4 * 4 + 4", "((4 * 4) + 4)"},

		{"a += 1", "(a = (a + 1))"},
		{"b -= 2", "(b = (b - 2))"},
		{"c *= 3", "(c = (c * 3))"},
		{"d /= 4", "(d = (d / 4))"},
		{"e %= 5", "(e = (e % 5))"},
		{"f <<= 6", "(f = (f << 6))"},
		{"g >>= 7", "(g = (g >> 7))"},

		{"a += 1 + 2", "(a = (a + (1 + 2)))"},

		{"var hoge = 1", "(var hoge = 1)"},
		{"var fuga = hoge * 2 + 2", "(var fuga = ((hoge * 2) + 2))"},

		{"return 0", "(return 0)"},
		{"return hoge", "(return hoge)"},

		{"func f(a: int): int {return 1}", "(def-func f(a: int): int ((return 1)))"},
		{"func hoge(fuga: int): int {return fuga}", "(def-func hoge(fuga: int): int ((return fuga)))"},

		{"if 1 { 1 } else { 0 }", "(if 1(1)(0))"},
		{"if 1 { 1 } else if 0 { 2 } else { 3 }", "(if 1(1)((if 0(2)(3))))"},

		{"while 1 { 1 }", "(while 1(1))"},

		{"a = a + 1", "(a = (a + 1))"},

		{"for i in 0..10 {1}", "(for i in 0..10(1))"},
		{"for var i = 0; i < 10; i = i + 1 {1}", "(for (var i = 0); (i < 10); (i = (i + 1))(1))"},

		{"array[1]", "(array[1])"},
		{"array[a * 10 + 1]", "(array[((a * 10) + 1)])"},

		{"new Foo", "(new Foo)"},
		{"foo.X", "(foo.X)"},
	}

	for i, test := range tests {
		l := lexer.NewFromString(test.input)
		p := New(l)
		program := p.parseStatement()
		checkParserErrors(t, p)
		actual := ast.Show(program)

		if actual != test.expected {
			t.Fatalf("tests[%d] - expected=%q, got=%q", i, test.expected, actual)
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

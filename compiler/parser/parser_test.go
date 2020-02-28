package parser

import (
	"github.com/arata-nvm/Solitude/compiler/lexer"
	"testing"
)

func TestParseProgram(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"func f(a: int) {return 1}", "func Ident(f)(Ident(a): Type(int)): Type(void) {return Int(1)}"},
		{"func hoge(fuga: int): int {return fuga}", "func Ident(hoge)(Ident(fuga): Type(int)): Type(int) {return Ident(fuga)}"},

		{"func num(): int {return 2} func main() {return num()}", "func Ident(num)(): Type(int) {return Int(2)}func Ident(main)(): Type(void) {return Call(Ident(num)())}"},
		{"func add(n: int): int {return n + 2} func main() {return num(1)}", "func Ident(add)(Ident(n): Type(int)): Type(int) {return Infix(Ident(n) + Int(2))}func Ident(main)(): Type(void) {return Call(Ident(num)(Int(1)))}"},
		{"func add(a: int, b: int): int {return a + b} func main() {return num(1, 2)}", "func Ident(add)(Ident(a): Type(int), Ident(b): Type(int)): Type(int) {return Infix(Ident(a) + Ident(b))}func Ident(main)(): Type(void) {return Call(Ident(num)(Int(1),Int(2)))}"},
		{"struct Foo { X int Y float }", "struct Ident(Foo) { Ident(X) int Ident(Y) float }"},
	}

	for i, test := range tests {
		l := lexer.NewFromString(test.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.Inspect()

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
		{"0", "Int(0)"},
		{"42", "Int(42)"},
		{"32.0", "Float(32.000000)"},

		{"4 + 4", "Infix(Int(4) + Int(4))"},
		{"4 - 4", "Infix(Int(4) - Int(4))"},
		{"4 * 4", "Infix(Int(4) * Int(4))"},
		{"4 / 4", "Infix(Int(4) / Int(4))"},
		{"4 % 4", "Infix(Int(4) % Int(4))"},
		{"4 << 4", "Infix(Int(4) << Int(4))"},
		{"4 >> 4", "Infix(Int(4) >> Int(4))"},
		{"4 == 4", "Infix(Int(4) == Int(4))"},
		{"4 != 4", "Infix(Int(4) != Int(4))"},
		{"4 < 4", "Infix(Int(4) < Int(4))"},
		{"4 <= 4", "Infix(Int(4) <= Int(4))"},
		{"4 > 4", "Infix(Int(4) > Int(4))"},
		{"4 >= 4", "Infix(Int(4) >= Int(4))"},

		{"4 + 4 * 4", "Infix(Int(4) + Infix(Int(4) * Int(4)))"},
		{"4 * 4 + 4", "Infix(Infix(Int(4) * Int(4)) + Int(4))"},

		{"a += 1", "Ident(a) = Infix(Ident(a) + Int(1))"},
		{"b -= 2", "Ident(b) = Infix(Ident(b) - Int(2))"},
		{"c *= 3", "Ident(c) = Infix(Ident(c) * Int(3))"},
		{"d /= 4", "Ident(d) = Infix(Ident(d) / Int(4))"},
		{"e %= 5", "Ident(e) = Infix(Ident(e) % Int(5))"},
		{"f <<= 6", "Ident(f) = Infix(Ident(f) << Int(6))"},
		{"g >>= 7", "Ident(g) = Infix(Ident(g) >> Int(7))"},

		{"a += 1 + 2", "Ident(a) = Infix(Ident(a) + Infix(Int(1) + Int(2)))"},

		{"var hoge = 1", "var Ident(hoge) = Int(1)"},
		{"var fuga = hoge * 2 + 2", "var Ident(fuga) = Infix(Infix(Ident(hoge) * Int(2)) + Int(2))"},

		{"return 0", "return Int(0)"},
		{"return hoge", "return Ident(hoge)"},

		{"func f(a: int) {return 1}", "func Ident(f)(Ident(a): Type(int)): Type(void) {return Int(1)}"},
		{"func hoge(fuga: int) {return fuga}", "func Ident(hoge)(Ident(fuga): Type(int)): Type(void) {return Ident(fuga)}"},

		{"if 1 { 1 } else { 0 }", "if Int(1) {Int(1)} else {Int(0)}"},

		{"while 1 { 1 }", "while Int(1) {Int(1)}"},

		{"a = a + 1", "Ident(a) = Infix(Ident(a) + Int(1))"},

		{"for i in 0..10 {1}", "for var Ident(i) = Int(0); Infix(Ident(i) <= Int(10)); Ident(i) = Infix(Ident(i) + Int(1)) {Int(1)}"},
		{"for var i = 0; i < 10; i = i + 1 {1}", "for var Ident(i) = Int(0); Infix(Ident(i) < Int(10)); Ident(i) = Infix(Ident(i) + Int(1)) {Int(1)}"},

		{"array[1]", "Ident(array)[Int(1)]"},
		{"array[a * 10 + 1]", "Ident(array)[Infix(Infix(Ident(a) * Int(10)) + Int(1))]"},

		{"new Foo", "new Ident(Foo)"},
		{"foo.X", "Ident(foo).Ident(X)"},
	}

	for i, test := range tests {
		l := lexer.NewFromString(test.input)
		p := New(l)
		program := p.parseStatement()
		checkParserErrors(t, p)
		actual := program.Inspect()

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

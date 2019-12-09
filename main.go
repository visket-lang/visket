package main

import (
	"bufio"
	"fmt"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/token"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimRight(input, "\n")

	l := lexer.New(input)

	fmt.Println("define i32 @main() nounwind {")

	tok := expectRead(l, token.INT)

	fmt.Println("  %1 = alloca i32, align 4")
	fmt.Printf("  store i32 %d, i32* %%1\n", tok.Val)

	op := l.NextToken()

	if op.Type == token.EOF {
		fmt.Println("  %2 = load i32, i32* %1, align 4")
		fmt.Println("  ret i32 %2")
		fmt.Println("}")
		os.Exit(1)
	}

	tok = expectRead(l, token.INT)

	fmt.Println("  %2 = alloca i32, align 4")
	fmt.Printf("  store i32 %d, i32* %%2", tok.Val)

	fmt.Println("  %3 = load i32, i32* %1, align 4")
	fmt.Println("  %4 = load i32, i32* %2, align 4")

	switch op.Type {
	case token.PLUS:
		fmt.Println("  %5 = add i32 %3, %4")
	case token.MINUS:
		fmt.Println("  %5 = sub i32 %3, %4")
	default:
		fmt.Printf("Unexpected token: %s\n", op)
		os.Exit(1)
	}

	fmt.Println("  %6 = alloca i32, align 4")
	fmt.Println("  store i32 %5, i32* %6, align 4")

	fmt.Println("  ret i32 %5")
	fmt.Println("}")
}

func expectRead(l *lexer.Lexer, tokenType token.TokenType) token.Token {
	tok := l.NextToken()
	if tok.Type != tokenType {
		fmt.Printf("Unexpected token: %s\n", tok)
		os.Exit(1)
	}

	return tok
}

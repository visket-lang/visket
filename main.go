package main

import (
	"bufio"
	"fmt"
	"github.com/arata-nvm/Solitude/ast"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/parser"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimRight(input, "\n")

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	genCode(program)
}

func genCode(program ast.Program) {
	fmt.Println("define i32 @main() nounwind {")

	index := 0

	index = gen(program.Code, index)

	numProcLastIndex := index
	index++
	regIndex := index

	fmt.Printf("  %%%d = load i32, i32* %%%d, align 4\n", regIndex, numProcLastIndex)
	fmt.Printf("  ret i32 %%%d\n", regIndex)
	fmt.Println("}")
}

func gen(node ast.Node, index int) int {
	switch node := node.(type) {
	case *ast.InfixExpression:
		return genInfix(node, index)
	case *ast.IntegerLiteral:
		index++
		fmt.Println("  ; Assign")
		fmt.Printf("  %%%d = alloca i32, align 4\n", index)
		fmt.Printf("  store i32 %d, i32* %%%d\n", node.Value, index)
	}
	return index
}

func genInfix(ie *ast.InfixExpression, index int) int {
	index = gen(ie.Left, index)
	lhsIndex := index
	index = gen(ie.Right, index)
	rhsIndex := index

	index++
	lhsRegIndex := index
	fmt.Printf("  %%%d= load i32, i32* %%%d, align 4\n", index, lhsIndex)

	index++
	rhsRegIndex := index
	fmt.Printf("  %%%d = load i32, i32* %%%d, align 4\n", index, rhsIndex)

	index++
	resRegIndex := index

	switch ie.Operator {
	case "+":
		fmt.Printf("  %%%d = add i32 %%%d, %%%d\n", index, lhsRegIndex, rhsRegIndex)
	case "-":
		fmt.Printf("  %%%d = sub i32 %%%d, %%%d\n", index, lhsRegIndex, rhsRegIndex)
	}

	index++
	resMemIndex := index
	fmt.Printf("  %%%d = alloca i32, align 4\n", index)

	fmt.Printf("  store i32 %%%d, i32* %%%d, align 4\n", resRegIndex, resMemIndex)

	return index
}

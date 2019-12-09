package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimRight(input, "\n")

	fmt.Println("define i32 @main() nounwind {")

	pos := 0

	n := readNumber(input, &pos)

	fmt.Println("  %1 = alloca i32, align 4")
	fmt.Printf("  store i32 %d, i32* %%1\n", n)

	if pos == len(input) {
		fmt.Println("  %2 = load i32, i32* %1, align 4")
		fmt.Println("  ret i32 %2")
		fmt.Println("}")
		os.Exit(1)
	}

	op := input[pos]
	pos ++

	n = readNumber(input, &pos)

	fmt.Println("  %2 = alloca i32, align 4")
	fmt.Printf("  store i32 %d, i32* %%2", n)

	fmt.Println("  %3 = load i32, i32* %1, align 4")
	fmt.Println("  %4 = load i32, i32* %2, align 4")

	switch op {
	case '+':
		fmt.Println("  %5 = add i32 %3, %4")
	case '-':
		fmt.Println("  %5 = sub i32 %3, %4")
	default:
		fmt.Printf("Unexpected char: %c\n", op)
		os.Exit(1)
	}


	fmt.Println("  %6 = alloca i32, align 4")
	fmt.Println("  store i32 %5, i32* %6, align 4")


	fmt.Println("  ret i32 %5")
	fmt.Println("}")

}

func readNumber(input string, pos *int) int {
	readPos := *pos
	for readPos < len(input) && isDigit(input[readPos]) {
		readPos ++
	}

	numLiteral := input[*pos:readPos]
	*pos = readPos

	n, err := strconv.Atoi(numLiteral)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return n
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

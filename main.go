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

	n, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("define i32 @main() nounwind {")
	fmt.Printf("  ret i32 %d\n", n)
	fmt.Println("}")

}

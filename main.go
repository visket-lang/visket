package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/arata-nvm/Solitude/codegen"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/parser"
	"os"
	"strings"
)

func main() {
	isDebug := flag.Bool("debug", false, "for debugging")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimRight(input, "\n")

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors) != 0 {
		for _, e := range p.Errors {
			fmt.Println(e)
		}
		os.Exit(1)
	}

	c := codegen.New(program, *isDebug)

	c.GenerateCode()
}

package main

import (
	"bufio"
	"flag"
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
	c := codegen.New(p, *isDebug)

	c.GenerateCode()
}

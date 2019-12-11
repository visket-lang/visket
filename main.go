package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/arata-nvm/Solitude/codegen"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/parser"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		isDebug = flag.Bool("debug", false, "for debugging")

		output = flag.String("output", "", "specify file to output")
	)
	flag.Parse()

	input := scanInput()
	l := lexer.New(input)

	p := parser.New(l)
	program := p.ParseProgram()
	if *isDebug {
		fmt.Printf("%s\n", program.Inspect())
	}
	printErrors(p)

	w := getWriter(*output)
	c := codegen.New(program, *isDebug, w)
	c.GenerateCode()
}

func scanInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return strings.TrimRight(input, "\n")
}

func printErrors(p *parser.Parser) {
	if len(p.Errors) != 0 {
		for _, e := range p.Errors {
			_, _ = fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}
}

func getWriter(output string) io.Writer {
	if output == "" {
		return os.Stdout
	} else {
		file, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		return file
	}
}

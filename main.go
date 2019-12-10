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

	var w io.Writer
	if *output == "" {
		w = os.Stdout
	} else {
		file, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		w = file
	}

	c := codegen.New(program, *isDebug, w)

	c.GenerateCode()
}

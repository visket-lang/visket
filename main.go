package main

import (
	"flag"
	"fmt"
	"github.com/arata-nvm/Solitude/codegen"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/parser"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var (
		isDebug = flag.Bool("v", false, "Use verbose output")
		output  = flag.String("o", "", "Specify file to output")
	)
	flag.Parse()

	input := scanInput()
	l := lexer.New(input)

	p := parser.New(l)
	program := p.ParseProgram()
	printErrors(p)
	if *isDebug {
		fmt.Printf("%s\n", program.Inspect())
	}

	w := getWriter(*output)
	c := codegen.New(program, *isDebug, w)
	c.GenerateCode()
}

func scanInput() string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
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

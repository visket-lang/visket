package main

import (
	"flag"
	"fmt"
	"github.com/arata-nvm/Solitude/codegen"
	"github.com/arata-nvm/Solitude/lexer"
	"github.com/arata-nvm/Solitude/optimizer"
	"github.com/arata-nvm/Solitude/parser"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var (
		output = flag.String("o", "", "Specify file to output")
	)
	flag.Parse()

	filepath := flag.Arg(0)
	if filepath == "" {
		fmt.Println("Usage: solitude <file>")
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	if *output == "" {
		*output = getFileNameWithoutExt(filepath)
	}

	genIrCode(string(content), *output)
	compile(*output)
}

func genIrCode(code, filename string) {
	l := lexer.New(code)

	p := parser.New(l)
	program := p.ParseProgram()
	printErrors(p)

	o := optimizer.New(program)
	o.Optimize()

	w, err := os.Create(filename + ".ll")
	if err != nil {
		log.Fatal(err)
	}

	c := codegen.New(program, false, w)
	c.GenerateCode()
}

func compile(filename string) {
	llPath := filename + ".ll"
	optLlPath := filename + ".opt.ll"
	asmPath := filename + ".s"

	cmd := exec.Command("opt", "-S", "-O3", llPath, "-o", optLlPath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("llc", optLlPath, "-o", asmPath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("cc", asmPath, "-o", filename, "-no-pie")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(llPath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(optLlPath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(asmPath)
	if err != nil {
		log.Fatal(err)
	}
}

func run(filename string) {
	cmd := exec.Command("./" + filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func printErrors(p *parser.Parser) {
	if len(p.Errors) != 0 {
		for _, e := range p.Errors {
			_, _ = fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

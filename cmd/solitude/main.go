package main

import (
	"flag"
	"fmt"
	"github.com/arata-nvm/Solitude/cmd/solitude/build"
	"log"
	"os"
)

func main() {
	var (
		isDebug  = flag.Bool("v", false, "Emit debug information")
		optimize = flag.Bool("O", false, "Enable optimization")
		output   = flag.String("o", "", "Write output to <filename>")
		emitLLVM = flag.Bool("emit-llvm", false, "Generate output in LLVM formats")
	)
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("Usage: solitude [options] <filename>")
		os.Exit(1)
	}

	if *emitLLVM {
		err := build.EmitLLVM(filename, *output, *isDebug, *optimize)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := build.Build(filename, *output, *isDebug, *optimize)
	if err != nil {
		log.Fatal(err)
	}
}

package errors

import (
	"fmt"
	"os"
	"strings"
)

type ErrorList []string

func (el ErrorList) ShowExit(verbose bool) {
	if len(el) == 0 {
		return
	}

	if !verbose {
		Error(el[0])
		if len(el) > 1 {
			Error(fmt.Sprintf("(and %d more errors)", len(el)-1))
		}
	} else {
		for _, e := range el {
			Error(e)
		}
	}
	os.Exit(1)
}

func Error(msg string) {
	msg = strings.ReplaceAll(msg, "\n", "\n\x1b[31merror\x1b[0m: ")
	fmt.Fprint(os.Stderr, "\x1b[31merror\x1b[0m: ")
	fmt.Fprintln(os.Stderr, msg)
}

func ErrorExit(msg string) {
	msg = strings.ReplaceAll(msg, "\n", "\n\x1b[31merror\x1b[0m: ")
	fmt.Fprint(os.Stderr, "\x1b[31merror\x1b[0m: ")
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

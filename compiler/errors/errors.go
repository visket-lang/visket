package errors

import (
	"fmt"
	"os"
)

type ErrorList []string

func (el ErrorList) ShowExit() {
	if len(el) != 0 {
		for _, e := range el {
			Error(e)
		}
		os.Exit(1)
	}
}

func Error(msg string) {
	fmt.Fprint(os.Stderr, "\x1b[31merror\x1b[0m: ")
	fmt.Fprintln(os.Stderr, msg)
}

func ErrorExit(msg string) {
	fmt.Fprint(os.Stderr, "\x1b[31merror\x1b[0m: ")
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

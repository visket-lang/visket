package errors

import (
	"fmt"
	"os"
	"strings"
)

type ErrorList []string

func (el ErrorList) ShowExit() {
	if len(el) != 0 {
		Error(fmt.Sprintf("parser has %d errors", len(el)))
		for _, e := range el {
			Error(e)
		}
		os.Exit(1)
	}
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

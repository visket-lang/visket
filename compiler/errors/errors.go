package errors

import (
	"fmt"
	"os"
	"strings"
)

type ErrorList []string

// TODO fix
var UseColors = false

const (
	errPrefix        = "error: "
	coloredErrPrefix = "\x1b[31merror\x1b[0m: "
)

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
	var prefix string
	if UseColors {
		prefix = coloredErrPrefix
	} else {
		prefix = errPrefix
	}

	msg = strings.ReplaceAll(msg, "\n", "\n"+prefix)
	fmt.Fprint(os.Stderr, prefix)
	fmt.Fprintln(os.Stderr, msg)
}

func ErrorExit(msg string) {
	Error(msg)
	os.Exit(1)
}

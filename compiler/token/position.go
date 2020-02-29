package token

import "fmt"

type Position struct {
	Filename string
	Line     int
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d", p.Filename, p.Line)
}

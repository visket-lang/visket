package internal

import "fmt"

var (
	labelNum = 0
)

func NextLabel(name string) string {
	labelNum++
	return fmt.Sprintf("%s.%d", name, labelNum)
}

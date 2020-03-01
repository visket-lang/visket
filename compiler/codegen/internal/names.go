package internal

import "fmt"

var (
	labelNum    = 0
	labelForNum = 0
)

func NextLabel(name string) string {
	labelNum++
	return fmt.Sprintf("%s.%d", name, labelNum)
}

func NextForNum(name string) string {
	labelForNum++
	return fmt.Sprintf("%s.%d", name, labelForNum)
}

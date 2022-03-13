package ui

import "fmt"

type Color int
const (
	Red Color = iota
	Green
	Yellow
	Blue
	Purple
	Cyan
	White
)

func (c Color) String() string {
	return []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m",
		"\033[35m", "\033[36m", "\033[37m"}[c]
}

func PrintlnColor(s string, c Color) {
	fmt.Printf("%s%s%s\n", c, s, "\033[0m")
}

func PrintColor(s string, c Color) {
	fmt.Printf("%s%s%s", c, s, "\033[0m")
}
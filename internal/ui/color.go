package ui

import (
	"fmt"
)

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

func PrintlnColor(c Color, a ...interface{}) {
	PrintfColor(c, "%s\n", a...)
}

func PrintfColor(c Color, format string, a ...interface{}) {
	fmt.Printf("%s%s%s", c, fmt.Sprintf(format, a...), "\033[0m")
}

package utils

import (
	"fmt"
	"mangathrV2/internal/utils/ui"
	"syscall"
)

func RaiseError(err error) {
	ui.PrintlnColor(fmt.Sprint(err), ui.Red)
	syscall.Exit(1)
}
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

func FindInSlice(list interface{}, match interface{}) (interface{}, bool) {
	switch list.(type) {
	case []string:
		for _, item := range list.([]string) {
			if item == match.(string) {
				return item, true
			}
		}
		return nil, false
	default:
		fmt.Println("unknown")
	}
	return nil, false
}

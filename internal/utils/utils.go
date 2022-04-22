package utils

import (
	"fmt"
	"mangathrV2/internal/utils/ui"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf8"
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

func PadString(s string, length int) string {
	stringSlice := strings.Split(s, ".")
	s = stringSlice[0]

	if utf8.RuneCountInString(s) >= length {
		return s
	}

	rString := strings.Repeat("0", length-utf8.RuneCountInString(s)) + s
	if len(stringSlice) > 1 {
		rString += "." + stringSlice[1]
	}
	return rString
}

func GetImageExtension(filename string) string {
	return "." + regexp.MustCompile(`.*\.(jpg|jpeg|webp|png)$`).FindAllStringSubmatch(filename, -1)[0][1]
}

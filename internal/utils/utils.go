package utils

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/utils/ui"
	"github.com/schollz/progressbar/v3"
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
		if len(stringSlice) > 1 {
			s += "." + stringSlice[1]
		}
		return s
	}

	rString := strings.Repeat("0", length-utf8.RuneCountInString(s)) + s
	if len(stringSlice) > 1 {
		rString += "." + stringSlice[1]
	}
	return rString
}

func GetImageExtension(filename string) string {
	return "." + regexp.MustCompile(`.*\.(jpg|jpeg|webp|png|gif)$`).FindAllStringSubmatch(filename, -1)[0][1]
}

func CreateProgressBar(length, maxRunes int, chapterNum string) *progressbar.ProgressBar {
	if rc := utf8.RuneCountInString(chapterNum); rc < maxRunes {
		chapterNum += strings.Repeat(" ", maxRunes-rc)
	}

	return progressbar.NewOptions(length,
		progressbar.OptionSetDescription(fmt.Sprintf("Chapter %s", chapterNum)),
	)
}

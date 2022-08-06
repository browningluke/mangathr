package utils

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf8"
)

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

func ExtractDate(dirtyDate string) string {
	re := regexp.MustCompile(`(\d{4})[-_/](\d{2})[-_/](\d{2})(T|$)`)

	result := re.FindAllStringSubmatch(dirtyDate, -1)
	if result == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s-%s", result[0][1], result[0][2], result[0][3])
}

func CreateProgressBar(length, maxRunes int, chapterNum string) *progressbar.ProgressBar {
	if rc := utf8.RuneCountInString(chapterNum); rc < maxRunes {
		chapterNum += strings.Repeat(" ", maxRunes-rc)
	}

	return progressbar.NewOptions(length,
		progressbar.OptionSetDescription(fmt.Sprintf("Chapter %s", chapterNum)),
	)
}

func IsRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false
	}
	return true
}

func CreateSigIntHandler(f func()) {
	SIGINT := make(chan os.Signal, 1)
	signal.Notify(SIGINT, os.Interrupt, syscall.SIGINT)
	go func() {
		<-SIGINT
		f()
		os.Exit(1)
	}()
}

package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf8"
)

func FindInSlice(list interface{}, match interface{}) (interface{}, bool) {
	return findInSlice(list, match, false)
}

func FindInSliceFold(list []string, match string) (string, bool) {
	s, o := findInSlice(list, match, true)

	if !o {
		s = ""
	}

	return s.(string), o
}

func findInSlice(list interface{}, match interface{}, stringFold bool) (interface{}, bool) {
	switch list.(type) {
	case []string:
		for _, item := range list.([]string) {
			if stringFold {
				// Check for fold match if folding
				if strings.EqualFold(item, match.(string)) {
					return item, true
				}
			} else {
				// Check for exact match if no folding
				if item == match.(string) {
					return item, true
				}
			}
		}
		return nil, false
	default:
		fmt.Println("FindInSlice: unknown type")
	}
	return nil, false
}

// MergeSlices merges multiple slices and removes duplicates.
func MergeSlices[K comparable](slices ...[]K) []K {
	// Create a map to keep track of unique elements
	uniqueElements := make(map[K]struct{})

	// Iterate over each slice
	for _, slice := range slices {
		for _, value := range slice {
			uniqueElements[value] = struct{}{}
		}
	}

	// Convert map keys to a slice
	var result []K
	for key := range uniqueElements {
		result = append(result, key)
	}

	return result
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

func ExpandHomePath(path string) string {
	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		return filepath.Join(dirname, path[2:])
	}
	return path
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

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

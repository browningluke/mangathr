package downloader

import (
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"github.com/browningluke/mangathrV2/internal/utils"
	"regexp"
	"strconv"
	"strings"
)

type Templater struct {
	RawTitle string
	Metadata structs.Metadata
}

func (t *Templater) handleNum(options string) string {
	if options == "" {
		return t.Metadata.Num
	} else {
		length, _ := strconv.ParseInt(strings.ReplaceAll(options, ":", ""), 10, 32)
		return utils.PadString(t.Metadata.Num, int(length))
	}
}

func (t *Templater) handleLanguage(options string) string {
	if t.Metadata.Language == "" {
		return ""
	}

	cleanString := strings.ReplaceAll(options, ":", "")
	return strings.ReplaceAll(cleanString, "<.>", t.Metadata.Language)
}

func (t *Templater) handleTitle(options string) string {
	if t.RawTitle == "" {
		return ""
	}

	cleanString := strings.ReplaceAll(options, ":", "")
	return strings.ReplaceAll(cleanString, "<.>", t.RawTitle)
}

func (t *Templater) handleGroups(options string) string {
	groups := strings.Join(t.Metadata.Groups, ", ")

	if groups == "" {
		return ""
	}

	cleanString := strings.ReplaceAll(options, ":", "")
	return strings.ReplaceAll(cleanString, "<.>", groups)
}

func (t *Templater) ExecTemplate(template string) string {
	re := regexp.MustCompile(`{((\w+?)(:.*?)?)}`)

	newString := template
	for _, match := range re.FindAllStringSubmatch(template, -1) {
		replace := match[0]

		varName := match[2]
		switch varName {
		case "num":
			options := ""
			if len(match) > 3 {
				options = match[3]
			}
			replace = t.handleNum(options)
		case "lang":
			replace = t.handleLanguage(match[3])
		case "title":
			replace = t.handleTitle(match[3])
		case "groups":
			replace = t.handleGroups(match[3])
		}

		newString = strings.Replace(newString, match[0], replace, 1)
	}

	return newString
}

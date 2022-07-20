package metadata

import (
	"fmt"
	"strings"
)

type comicInfoAgent struct {
	template string
}

func newComicInfoAgent() *comicInfoAgent {
	t := fmt.Sprintf("<?xml version=\"1.0\"?>\n" +
		"<ComicInfo>\n")

	return &comicInfoAgent{template: t}
	//return &Agent{title: title, num: num}
}

func (a *comicInfoAgent) GenerateMetadataFile() (filename, body string) {
	return "ComicInfo.xml", a.template + "</ComicInfo>"
}

// Setters

func (a *comicInfoAgent) SetTitle(title string) Agent {
	a.template += fmt.Sprintf("<Title>%s</Title>\n", title)
	return a
}

func (a *comicInfoAgent) SetNum(num string) Agent {
	a.template += fmt.Sprintf("<Number>%s</Number>\n", num)
	return a
}

// SetDate in yyyy-mm-dd format
func (a *comicInfoAgent) SetDate(date string) Agent {
	// TODO use a better method for extracting the date
	a.template += fmt.Sprintf("<Year>%s</Year>\n", date[0:4])
	a.template += fmt.Sprintf("<Month>%s</Month>\n", date[5:7])
	a.template += fmt.Sprintf("<Day>%s</Day>\n", date[8:10])
	return a
}

func (a *comicInfoAgent) SetEditors(editors []string) Agent {
	a.template += fmt.Sprintf("<Editor>%s</Editor>\n", strings.Join(editors[:], ", "))
	return a
}

func (a *comicInfoAgent) SetWebLink(link string) Agent {
	a.template += fmt.Sprintf("<Web>%s</Web>\n", link)
	return a
}

func (a *comicInfoAgent) SetPageCount(count int) Agent {
	a.template += fmt.Sprintf("<PageCount>%d</PageCount>\n", count)
	return a
}

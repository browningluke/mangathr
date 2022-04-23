package comicinfo

import (
	"fmt"
)

type Agent struct {
	template string
}

func NewAgent() *Agent {
	t := fmt.Sprintf("<?xml version=\"1.0\"?>\n" +
		"<ComicInfo xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"\n" +
		"xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\">\n")

	return &Agent{template: t}
	//return &Agent{title: title, num: num}
}

func (a *Agent) GenerateMetadataFile() (filename, body string) {
	return "ComicInfo.xml", a.template + "</ComicInfo>"
}

// Setters

func (a *Agent) SetTitle(title string) {
	a.template += fmt.Sprintf("<Title>%s</Title>\n", title)
}

func (a *Agent) SetNum(num string) {
	a.template += fmt.Sprintf("<Number>%s</Number>\n", num)
}

// SetDate in yyyy-mm-dd format
func (a *Agent) SetDate(date string) {
	// TODO use a better method for extracting the date
	a.template += fmt.Sprintf("<Year>%s</Year>\n", date[0:4])
	a.template += fmt.Sprintf("<Month>%s</Month>\n", date[5:7])
	a.template += fmt.Sprintf("<Day>%s</Day>\n", date[8:10])
}

func (a *Agent) SetEditor(editor string) {
	a.template += fmt.Sprintf("<Editor>%s</Editor>\n", editor)
}

func (a *Agent) SetWebLink(link string) {
	a.template += fmt.Sprintf("<Web>%s</Web>\n", link)
}

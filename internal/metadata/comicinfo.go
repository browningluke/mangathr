package metadata

import (
	"encoding/xml"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"strings"
)

type comicInfoXML struct {
	XMLName xml.Name `xml:"ComicInfo"`

	Title string `xml:"Title,omitempty"`
	Num   string `xml:"Number,omitempty"`
	Link  string `xml:"Web,omitempty"`

	Year  string `xml:"Year,omitempty"`
	Month string `xml:"Month,omitempty"`
	Day   string `xml:"Day,omitempty"`

	Editor string `xml:"Editor,omitempty"`

	PageCount int `xml:"PageCount,omitempty"`

	Language string `xml:"Language,omitempty"`
}

type comicInfoAgent struct {
	template comicInfoXML
}

func newComicInfoAgent() *comicInfoAgent {
	return &comicInfoAgent{template: comicInfoXML{}}
}

func (a *comicInfoAgent) GenerateMetadataFile() (filename string, body []byte, err error) {
	data, err := xml.MarshalIndent(a.template, " ", "  ")
	if err != nil {
		return "", []byte{}, err
	}

	data = []byte(xml.Header + string(data))
	return "ComicInfo.xml", data, err
}

// Setters

func (a *comicInfoAgent) SetTitle(title string) Agent {
	a.template.Title = title
	return a
}

func (a *comicInfoAgent) SetNum(num string) Agent {
	a.template.Num = num
	return a
}

// SetDate in yyyy-mm-dd format
func (a *comicInfoAgent) SetDate(date string) Agent {
	// TODO use a better method for extracting the date
	a.template.Year = fmt.Sprint(date[0:4])
	a.template.Month = fmt.Sprint(date[5:7])
	a.template.Day = fmt.Sprint(date[8:10])
	return a
}

func (a *comicInfoAgent) SetEditors(editors []string) Agent {
	a.template.Editor = fmt.Sprint(strings.Join(editors[:], ", "))
	return a
}

func (a *comicInfoAgent) SetWebLink(link string) Agent {
	a.template.Link = link
	return a
}

func (a *comicInfoAgent) SetPageCount(count int) Agent {
	a.template.PageCount += count
	return a
}

// SetFromStruct ingests all data from Metadata struct EXCEPT PageCount
func (a *comicInfoAgent) SetFromStruct(metadata structs.Metadata) Agent {
	a.SetTitle(metadata.Title).
		SetNum(metadata.Num).
		SetWebLink(metadata.Link).
		SetDate(metadata.Date).
		SetEditors(metadata.Groups)
	return a
}

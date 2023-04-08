package metadata

import (
	"encoding/xml"
	"fmt"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/manga"
	"github.com/browningluke/mangathr/internal/utils"
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
	date = utils.ExtractDate(date)
	if date == "" {
		logging.Warningf("Attempted to parse garbage date string: %s\n", date)
		return a
	}

	dateSlice := strings.Split(date, "-")

	a.template.Year = dateSlice[0]
	a.template.Month = dateSlice[1]
	a.template.Day = dateSlice[2]
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
func (a *comicInfoAgent) SetFromStruct(metadata manga.Metadata) Agent {
	a.SetTitle(metadata.Title).
		SetNum(metadata.Num).
		SetWebLink(metadata.Link).
		SetDate(metadata.Date).
		SetEditors(metadata.Groups)
	return a
}

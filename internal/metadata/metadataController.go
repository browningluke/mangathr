package metadata

import "github.com/browningluke/mangathr/internal/manga"

type Agent interface {
	GenerateMetadataFile() (filename string, body []byte, err error)

	SetTitle(title string) Agent
	SetNum(num string) Agent
	SetDate(date string) Agent // MUST BE yyyy-mm-dd
	SetEditors(editors []string) Agent
	SetWebLink(link string) Agent
	SetPageCount(count int) Agent

	SetFromStruct(metadata manga.Metadata) Agent
}

func NewAgent(name string) Agent {
	m := map[string]func() Agent{
		"comicinfo": func() Agent { return newComicInfoAgent() },
	}

	agent, ok := m[name]
	if !ok {
		panic("Passed agent name not in map")
	}
	return agent()
}

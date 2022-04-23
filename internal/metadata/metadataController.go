package metadata

import (
	"mangathrV2/internal/metadata/comicinfo"
)

type Agent interface {
	GenerateMetadataFile() (filename, body string)

	SetTitle(title string)
	SetNum(num string)
	SetDate(date string) // MUST BE yyyy-mm-dd
	SetEditors(editors []string)
	SetWebLink(link string)
}

func NewAgent(name string) Agent {
	m := map[string]func() Agent{
		"comicinfo": func() Agent { return comicinfo.NewAgent() },
	}

	agent, ok := m[name]
	if !ok {
		panic("Passed agent name not in map")
	}
	return agent()
}

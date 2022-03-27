package metadata

import (
	"mangathrV2/internal/metadata/comicinfo"
)

type Agent interface {
	GenerateMetadataFile() (filename, body string)
}

func NewAgent(name, title, num string) Agent {
	m := map[string]func() Agent{
		"comicinfo": func() Agent { return comicinfo.NewAgent(title, num) },
	}

	agent, ok := m[name]
	if !ok {
		panic("Passed agent name not in map")
	}
	return agent()
}

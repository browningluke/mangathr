package comicinfo

import "fmt"

type Agent struct {
	title, num string
}

func NewAgent(title, num string) *Agent {
	fmt.Println("Created a comicinfo agent")
	return &Agent{title: title, num: num}
}

func (a *Agent) GenerateMetadataFile() (filename, body string) {
	b := fmt.Sprintf(`
<?xml version="1.0"?>
<ComicInfo xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	<Number>%s</Number>
	<Title>%s</Title>
</ComicInfo>
	`, a.num, a.title)

	return "ComicInfo.xml", b
}

package cubari

import (
	"fmt"
	"mangathrV2/internal/downloader"
)

type Scraper struct {
	name string
}

func NewScraper() *Scraper {
	fmt.Println("Created a cubari scraper")
	return &Scraper{}
}

func (m *Scraper) Search(query string) interface{} {
	//TODO implement me
	panic("implement me")

}

func (m *Scraper) SearchByID(id string) interface{} {
	//TODO implement me
	panic("implement me")
}

func (m *Scraper) ListChapters() interface{} {
	fmt.Println("List chapters")
	m.name = "test"
	fmt.Println(m.name)

	//TODO implement me
	//panic("implement me")
	return nil
}

func (m *Scraper) SelectChapters() interface{} {
	fmt.Println("Select chapters")
	m.name = "test"
	fmt.Println(m.name)

	//TODO implement me
	//panic("implement me")
	return nil
}

func (m *Scraper) Download(downloader *downloader.Downloader) {
	//TODO implement me
	panic("implement me")
}

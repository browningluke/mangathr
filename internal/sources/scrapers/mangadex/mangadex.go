package mangadex

import (
	"encoding/json"
	"fmt"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/rester"
	"net/url"
)

type Scraper struct {
	name string

	searchResults []searchResult
	allChapters   []chapter

	manga            searchResult
	selectedChapters []chapter
}

type searchResult struct {
	name string
	id   string
}

type chapter struct {
	name string
	num  string // doing this will make it easier
}

func NewScraper() *Scraper {
	fmt.Println("Created a mangadex scraper")
	return &Scraper{}
}

// Search for a Manga, will fill searchResults with 0 or more results
func (m *Scraper) Search(query string) []string {
	jsonString := rester.New().Get(
		"https://api.mangadex.org/manga?limit=10&title="+url.QueryEscape(query),
		map[string]string{})

	var mangaResp mangaResponse

	err := json.Unmarshal([]byte(jsonString), &mangaResp)
	if err != nil {
		panic(err)
	}

	var searchResults []searchResult
	var names []string

	for _, item := range mangaResp.Data {
		searchResults = append(searchResults, searchResult{name: item.Attributes.Title["en"], id: item.Id})
		names = append(names, item.Attributes.Title["en"])
	}

	//fmt.Println(mangaResp)
	m.searchResults = searchResults

	return names
}

// SelectManga from searchResults list
func (m *Scraper) SelectManga(name string) {

	found := false
	for _, item := range m.searchResults {
		if item.name == name {
			m.manga = item
			found = true
			break
		}
	}

	if !found {
		panic("Selected manga not is search result list")
	}

	// Once manga has been selected, clear all search results
	m.searchResults = []searchResult{}
}

// SearchByID for a Manga, will fill searchResults with ONLY 1 result (first result)
func (m *Scraper) SearchByID(id string) interface{} {
	//TODO implement me
	panic("implement me")
}

func (m *Scraper) ListChapters() []string {
	m.allChapters = []chapter{
		{name: "Chapter 100 - something", num: "100"},
		{name: "Chapter 101 - something else", num: "101"},
		{name: "Chapter 101.5 - something extra", num: "101.5"},
	}

	var names []string

	for _, item := range m.allChapters {
		names = append(names, item.name)
	}

	return names // do formatting on names here
}

func (m *Scraper) SelectChapters(titles []string) {
	var chapters []chapter

	for _, chapter := range m.allChapters {
		for _, title := range titles {
			if chapter.name == title {
				chapters = append(chapters, chapter)
			}
		}
	}
	m.selectedChapters = chapters

	// Once chapters have been selected, clear all chapters
	m.allChapters = []chapter{}

	fmt.Println("Selected chapters: ", m.selectedChapters)
}

func (m *Scraper) Download(downloader *downloader.Downloader) {
	//TODO implement me
	panic("implement me")
}

func (m *Scraper) GetChapterTitle() string {
	return m.manga.name
}

func (m *Scraper) GetScraperName() string {
	return "Mangadex"
}

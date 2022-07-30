package mangadex

import (
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
)

const (
	SCRAPERNAME            = "Mangadex"
	APIROOT                = "https://api.mangadex.org"
	CHAPTERSPERMIN         = 40 // set from API docs
	ENFORCECHAPTERDURATION = true
)

type searchResult struct {
	title, id string
}

type chapterResult struct {
	title, // The raw data from MD
	fullTitle, // The formatted title for filename
	id, language string
	sortNum float64

	metadata structs.Metadata
}

type Scraper struct {
	name   string
	config *Config

	searchResults []searchResult
	manga         searchResult

	allChapters, selectedChapters []chapterResult

	// Group queries
	groups []string
}

func NewScraper(config *Config) *Scraper {
	logging.Debugln("Created a Mangadex scraper")
	s := &Scraper{config: config}
	return s
}

/*
	-- Get Chapter data --
*/

func (m *Scraper) Chapters() []structs.Chapter {
	if len(m.allChapters) == 0 && len(m.selectedChapters) == 0 {
		m.scrapeChapters()
	}

	c := m.allChapters

	if len(m.selectedChapters) != 0 {
		c = m.selectedChapters
	}

	var chapters []structs.Chapter
	for _, item := range c {
		chapters = append(chapters,
			structs.Chapter{ID: item.id, Title: item.fullTitle, Metadata: item.metadata})
	}
	return chapters
}

func (m *Scraper) ChapterTitles() []string {
	if len(m.allChapters) == 0 && len(m.selectedChapters) == 0 {
		m.scrapeChapters()
	}

	chapters := m.allChapters

	if len(m.selectedChapters) != 0 {
		chapters = m.selectedChapters
	}

	var titles []string
	for _, item := range chapters {
		titles = append(titles, item.fullTitle)
	}
	return titles
}

/*
	-- Get/Set Group data --
*/

func (m *Scraper) GroupNames() []string {
	if len(m.groups) == 0 {
		m.Chapters()
	}

	var groupNames []string
	for _, val := range m.groups {
		groupNames = append(groupNames, val)
	}
	return groupNames
}

func (m *Scraper) FilterGroups(groups []string) {
	findElemInSlice := func(slice []string, elem string) bool {
		for _, v := range slice {
			if elem == v {
				return true
			}
		}
		return false
	}

	var selectedChapters []chapterResult
	for _, chapter := range m.allChapters { // go through each chapter
		for _, group := range groups { // go through each filtered group
			exists := findElemInSlice(chapter.metadata.Groups, group)
			if exists {
				selectedChapters = append(selectedChapters, chapter)
				break
			}
		}
	}
	m.selectedChapters = selectedChapters
}

/*
	-- Getters --
*/

func (m *Scraper) MangaTitle() string {
	return m.manga.title
}

func (m *Scraper) MangaID() string {
	return m.manga.id
}

func (m *Scraper) ScraperName() string {
	return SCRAPERNAME
}

func (m *Scraper) EnforceChapterDuration() bool {
	return ENFORCECHAPTERDURATION
}

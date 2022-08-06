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

type Scraper struct {
	config *Config

	searchResults []searchResult
	manga         searchResult

	allChapters, selectedChapters []structs.Chapter

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

func (m *Scraper) Chapters() ([]structs.Chapter, *logging.ScraperError) {
	if len(m.allChapters) == 0 && len(m.selectedChapters) == 0 {
		if err := m.scrapeChapters(); err != nil {
			return nil, err
		}
	}

	c := m.allChapters

	if len(m.selectedChapters) != 0 {
		c = m.selectedChapters
	}

	return c, nil
}

// ChapterTitles Returns the full titles of chapters
func (m *Scraper) ChapterTitles() ([]string, *logging.ScraperError) {
	if len(m.allChapters) == 0 && len(m.selectedChapters) == 0 {
		if err := m.scrapeChapters(); err != nil {
			return nil, err
		}
	}

	chapters := m.allChapters
	if len(m.selectedChapters) != 0 {
		chapters = m.selectedChapters
	}

	var titles []string
	for _, item := range chapters {
		titles = append(titles, item.FullTitle)
	}
	return titles, nil
}

/*
	-- Get/Set Group data --
*/

func (m *Scraper) GroupNames() ([]string, *logging.ScraperError) {
	if len(m.groups) == 0 {
		if _, err := m.Chapters(); err != nil {
			return nil, err
		}
	}
	return m.groups, nil
}

func (m *Scraper) FilterGroups(groups []string) *logging.ScraperError {
	findElemInSlice := func(slice []string, elem string) bool {
		for _, v := range slice {
			if elem == v {
				return true
			}
		}
		return false
	}

	var selectedChapters []structs.Chapter
	for _, chapter := range m.allChapters { // go through each chapter
		for _, group := range groups { // go through each filtered group
			exists := findElemInSlice(chapter.Metadata.Groups, group)
			if exists {
				selectedChapters = append(selectedChapters, chapter)
				break
			}
		}
	}
	m.selectedChapters = selectedChapters

	return nil
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

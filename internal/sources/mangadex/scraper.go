package mangadex

import (
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/manga"
)

const (
	SCRAPERNAME            = "Mangadex"
	APIROOT                = "https://api.mangadex.org"
	CHAPTERSPERMIN         = 40 // set from API docs
	ENFORCECHAPTERDURATION = true
	REGISTRABLE            = true
)

type searchResult struct {
	title, id string
}

type Scraper struct {
	searchResults []searchResult
	manga         searchResult

	allChapters, selectedChapters,
	filteredChapters []manga.Chapter

	// Group queries
	groups []string
}

func NewScraper() *Scraper {
	logging.Debugln("Created a Mangadex scraper")
	s := &Scraper{}
	return s
}

/*
	-- Get Chapter data --
*/

func (m *Scraper) Chapters() ([]manga.Chapter, *logging.ScraperError) {
	// If chapters have been filtered, only show the filtered chapters
	if len(m.filteredChapters) != 0 {
		return m.filteredChapters, nil
	}

	// If parsing has already been done, skip repeating it
	if len(m.allChapters) != 0 {
		return m.allChapters, nil
	}

	// Otherwise, parse chapters and return
	var err *logging.ScraperError
	m.allChapters, err = m.scrapeChapters()
	return m.allChapters, err
}

// ChapterTitles Returns the full titles of chapters
func (m *Scraper) ChapterTitles() ([]string, *logging.ScraperError) {
	// Parse chapters if not already done
	chapters, err := m.Chapters()
	if err != nil {
		return []string{}, err
	}

	var chapterTitles []string
	for _, c := range chapters {
		chapterTitles = append(chapterTitles, c.FullTitle)
	}

	return chapterTitles, nil
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
	// Ensure chapters are parsed
	m.Chapters()

	findElemInSlice := func(slice []string, elem string) bool {
		for _, v := range slice {
			if elem == v {
				return true
			}
		}
		return false
	}

	var filteredChapters []manga.Chapter
	for _, chapter := range m.allChapters { // go through each chapter
		for _, group := range groups { // go through each filtered group
			exists := findElemInSlice(chapter.Metadata.Groups, group)
			if exists {
				filteredChapters = append(filteredChapters, chapter)
				break
			}
		}
	}

	m.filteredChapters = filteredChapters

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

func (m *Scraper) Registrable() bool {
	return REGISTRABLE
}

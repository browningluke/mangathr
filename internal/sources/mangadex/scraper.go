package mangadex

import (
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
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
	filtered bool

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
	if m.filtered {
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

func filterGroups(chapters []manga.Chapter, groups []string, exclude bool) []manga.Chapter {
	findElemInSlice := func(slice []string, elem string) bool {
		for _, v := range slice {
			if elem == v {
				return true
			}
		}
		return false
	}

	var filteredChapters []manga.Chapter
	for _, chapter := range chapters { // go through each chapter
		// if excluding: starts as true, inverts if match found
		// if including: start as false, inverts if match found
		addChapter := exclude

		for _, group := range groups { // go through each filtered group
			if exists := findElemInSlice(chapter.Metadata.Groups, group); exists {
				addChapter = !addChapter
				break
			}
		}

		if addChapter {
			filteredChapters = append(filteredChapters, chapter)
		}
	}

	return filteredChapters
}

func (m *Scraper) FilterGroups(includeGroups []string, excludeGroups []string) *logging.ScraperError {
	// Ensure chapters are parsed
	if _, err := m.Chapters(); err != nil {
		return err
	}

	// Start with all
	chaptersToFilter := m.allChapters

	// Include
	if len(includeGroups) > 0 {
		chaptersToFilter = filterGroups(chaptersToFilter, includeGroups, false)
	}

	// Exclude
	if len(excludeGroups) > 0 {
		chaptersToFilter = filterGroups(chaptersToFilter, excludeGroups, true)
	}

	m.filteredChapters = chaptersToFilter

	// Mark filtering done
	m.filtered = true

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

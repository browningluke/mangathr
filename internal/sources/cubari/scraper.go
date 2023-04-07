package cubari

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"github.com/browningluke/mangathrV2/internal/utils"
	"strings"
)

const (
	SCRAPERNAME            = "Cubari"
	APIROOT                = "https://cubari.moe/read/api"
	SITEROOT               = "https://cubari.moe"
	ENFORCECHAPTERDURATION = false
)

type Scraper struct {
	allChapters, selectedChapters,
	filteredChapters []structs.Chapter

	// pages URLs mapped by chapter ID
	pages map[string][]string

	manga    mangaResponse
	provider Provider
}

func NewScraper() *Scraper {
	logging.Debugln("Created a Cubari scraper")
	s := &Scraper{
		pages: make(map[string][]string),
	}
	return s
}

// SelectManga does nothing meaningful with Cubari, since there's only ever 1 manga at a time
func (m *Scraper) SelectManga(_ string) *logging.ScraperError {
	// Only 1 manga is ever handled at a time
	return nil
}

/*
	-- Get Chapter data --
*/

// Chapters returns chapter data from Cubari's API
func (m *Scraper) Chapters() ([]structs.Chapter, *logging.ScraperError) {
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
	m.allChapters, err = m.parseChapters()
	return m.allChapters, err
}

// ChapterTitles returns a list of all chapter titles
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

// GroupNames returns a list of all the group names
func (m *Scraper) GroupNames() ([]string, *logging.ScraperError) {
	// Ensure chapters are parsed
	if _, err := m.Chapters(); err != nil {
		return []string{}, err
	}

	var groupNames []string
	for _, v := range m.manga.Groups {
		v = strings.ReplaceAll(v, "/", ", ")
		groupNames = append(groupNames, v)
	}

	return groupNames, nil
}

// FilterGroups to find all chapters with groups in groups list
func (m *Scraper) FilterGroups(groups []string) *logging.ScraperError {
	// Ensure chapters are parsed
	if _, err := m.Chapters(); err != nil {
		return err
	}

	var filteredChapters []structs.Chapter

	for _, v := range m.allChapters {
		// Assuming chapters must have 1 and only 1 group (as Cubari does with GIST provider)
		if _, ok := utils.FindInSlice(groups, v.Metadata.Groups[0]); ok {
			filteredChapters = append(filteredChapters, v)
		}
	}

	m.filteredChapters = filteredChapters

	return nil
}

/*
	-- Getters --
*/

func (m *Scraper) MangaTitle() string {
	return m.manga.Title
}

func (m *Scraper) MangaID() string {
	return fmt.Sprintf("%s~%s", m.provider.name, m.manga.Slug)
}

func (m *Scraper) ScraperName() string {
	return SCRAPERNAME
}

func (m *Scraper) EnforceChapterDuration() bool {
	return ENFORCECHAPTERDURATION
}

func (m *Scraper) Registrable() bool {
	return m.provider.registrable
}

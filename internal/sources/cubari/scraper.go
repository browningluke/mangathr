package cubari

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
	"github.com/browningluke/mangathr/v2/internal/utils"
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
	filteredChapters []manga.Chapter
	filtered bool

	// pages URLs mapped by chapter ID
	pages map[string][]string

	manga    mangaResponse
	provider Provider
}

func NewScraper() *Scraper {
	logging.Debugln("Created a Cubari scraper")
	s := &Scraper{
		pages:    make(map[string][]string),
		filtered: false,
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

func filterGroups(chapters []manga.Chapter, groups []string, exclude bool) []manga.Chapter {
	var filteredChapters []manga.Chapter

	for _, v := range chapters {
		// Assuming chapters must have 1 and only 1 group (as Cubari does with GIST provider)
		_, ok := utils.FindInSlice(groups, v.Metadata.Groups[0])

		// If we're excluding, invert the answer
		if exclude {
			ok = !ok
		}

		if ok {
			filteredChapters = append(filteredChapters, v)
		}
	}

	return filteredChapters
}

// FilterGroups to find all chapters with groups in groups list
func (m *Scraper) FilterGroups(includeGroups []string, excludeGroups []string) *logging.ScraperError {
	// Ensure chapters are parsed
	if _, err := m.Chapters(); err != nil {
		return err
	}

	// Merge `groups` with config value
	toInclude := utils.MergeSlices(includeGroups, config.Groups.Include)
	toExclude := utils.MergeSlices(excludeGroups, config.Groups.Exclude)

	// Start with all
	chaptersToFilter := m.allChapters

	// Include
	if len(toInclude) > 0 {
		chaptersToFilter = filterGroups(chaptersToFilter, toInclude, false)
	}

	// Exclude
	if len(toExclude) > 0 {
		chaptersToFilter = filterGroups(chaptersToFilter, toExclude, true)
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

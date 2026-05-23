package mangaplus

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
	"github.com/google/uuid"
)

const (
	SCRAPERNAME            = "MangaPlus"
	APIROOT                = "https://jumpg-webapi.tokyo-cdn.com/api"
	REFERERURL             = "https://mangaplus.shueisha.co.jp"
	USERAGENT              = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"
	ENFORCECHAPTERDURATION = false
	REGISTRABLE            = true
)

type searchResult struct {
	titleId uint32
	name    string
}

type Scraper struct {
	searchResults []searchResult
	manga         searchResult

	sessionToken string

	allChapters, selectedChapters,
	filteredChapters []manga.Chapter
	filtered bool
}

func NewScraper() *Scraper {
	logging.Debugln("Created a MangaPlus scraper")
	return &Scraper{sessionToken: uuid.New().String()}
}

func (m *Scraper) headers(referer string) map[string]string {
	return map[string]string{
		"Origin":        REFERERURL,
		"User-Agent":    USERAGENT,
		"Session-Token": m.sessionToken,
		"Referer":       referer,
	}
}

func (m *Scraper) Chapters() ([]manga.Chapter, *logging.ScraperError) {
	if m.filtered {
		return m.filteredChapters, nil
	}
	if len(m.allChapters) != 0 {
		return m.allChapters, nil
	}

	var err *logging.ScraperError
	m.allChapters, err = m.scrapeChapters()
	return m.allChapters, err
}

func (m *Scraper) ChapterTitles() ([]string, *logging.ScraperError) {
	chapters, err := m.Chapters()
	if err != nil {
		return []string{}, err
	}

	var titles []string
	for _, chapter := range chapters {
		titles = append(titles, chapter.FullTitle)
	}
	return titles, nil
}

func (m *Scraper) GroupNames() ([]string, *logging.ScraperError) {
	return []string{}, nil
}

func (m *Scraper) FilterGroups(_, _ []string) *logging.ScraperError {
	return nil
}

func (m *Scraper) MangaTitle() string {
	return m.manga.name
}

func (m *Scraper) MangaID() string {
	return fmt.Sprintf("%d", m.manga.titleId)
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

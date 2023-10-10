package cubari

import (
	"errors"
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
	"github.com/browningluke/mangathr/v2/internal/rester"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"math"
)

func (m *Scraper) addPagesToChapter(chapter *manga.Chapter) *logging.ScraperError {
	pages := m.pages[chapter.ID]

	// Get pages from proxy URL
	// (if using GIST provider)
	if m.provider == GIST {
		jsonResp, _ := rester.New().Get(fmt.Sprintf("%s%s", SITEROOT, pages[0]),
			map[string]string{}, []rester.QueryParam{}).Do(4, "100ms")
		jsonString := jsonResp.(string)

		urls, ok := parseImgurStyle([]byte(jsonString))
		if !ok {
			return &logging.ScraperError{
				Error:   errors.New("failed to get imgur URLs from proxy"),
				Message: "An error occurred while getting pages from imgur",
				Code:    0,
			}
		}

		pages = urls
	}

	// (if using MANGASEE provider)
	if m.provider == MANGASEE {
		jsonResp, _ := rester.New().Get(fmt.Sprintf("%s%s", SITEROOT, pages[0]),
			map[string]string{}, []rester.QueryParam{}).Do(4, "100ms")
		jsonString := jsonResp.(string)

		urls, ok := parseListStyle([]byte(jsonString))
		if !ok {
			return &logging.ScraperError{
				Error:   errors.New("failed to get mangasee URLs from proxy"),
				Message: "An error occurred while getting pages from mangasee",
				Code:    0,
			}
		}

		pages = urls
	}

	digits := int(math.Floor(math.Log10(float64(len(pages)))) + 1)

	for i, url := range pages {
		chapter.AddPage(url, utils.PadString(fmt.Sprintf("%d", i+1), digits))
	}

	return nil
}

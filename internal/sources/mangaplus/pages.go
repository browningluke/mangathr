package mangaplus

import (
	"encoding/hex"
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
	"github.com/browningluke/mangathr/v2/internal/rester"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"math"
)

func hexToBytes(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return decoded
}

func (m *Scraper) addPagesToChapter(chapter *manga.Chapter) *logging.ScraperError {
	respBytes, _ := rester.New().GetBytes(
		fmt.Sprintf("%s/manga_viewer", APIROOT),
		m.headers(fmt.Sprintf("%s/viewer/%s", REFERERURL, chapter.ID)),
		[]rester.QueryParam{
			{Key: "chapter_id", Value: chapter.ID},
			{Key: "split", Value: "yes"},
			{Key: "img_quality", Value: "super_high"},
		},
	).Do(4, "100ms")
	data := respBytes.([]byte)

	resp, err := decodeResponse(data)
	if err != nil {
		return &logging.ScraperError{
			Error:   err,
			Message: "Failed to decode manga viewer response",
			Code:    0,
		}
	}

	viewer := resp.GetSuccess().GetMangaViewer()
	if viewer == nil {
		return &logging.ScraperError{
			Error:   fmt.Errorf("mangaViewer is nil"),
			Message: "Failed to get manga viewer",
			Code:    0,
		}
	}

	pages := viewer.GetPages()
	digits := 1
	if len(pages) > 0 {
		digits = int(math.Floor(math.Log10(float64(len(pages)))) + 1)
	}

	pageNum := 0
	for _, page := range pages {
		mangaPage := page.GetPage()
		if mangaPage == nil {
			continue
		}

		pageNum++
		pageName := utils.PadString(fmt.Sprintf("%d", pageNum), digits)
		encryptionKey := mangaPage.GetEncryptionKey()
		if encryptionKey == "" {
			chapter.AddPage(mangaPage.GetImageUrl(), pageName)
			continue
		}

		key := append([]byte(nil), hexToBytes(encryptionKey)...)
		if len(key) == 0 {
			return &logging.ScraperError{
				Error:   fmt.Errorf("invalid encryption key for page %s", pageName),
				Message: "Failed to parse page encryption key",
				Code:    0,
			}
		}

		chapter.AddPageWithTransform(mangaPage.GetImageUrl(), pageName, func(b []byte) ([]byte, error) {
			result := make([]byte, len(b))
			for i, value := range b {
				result[i] = value ^ key[i%len(key)]
			}
			return result, nil
		})
	}

	return nil
}

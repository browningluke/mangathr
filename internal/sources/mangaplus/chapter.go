package mangaplus

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
	"github.com/browningluke/mangathr/v2/internal/rester"
	mpproto "github.com/browningluke/mangathr/v2/internal/sources/mangaplus/proto"
	"regexp"
	"strconv"
	"strings"
)

var chapterTitleRe = regexp.MustCompile(`.+?\d.? (\w.+?)$`)

func parseChapterNum(chapters []*mpproto.Chapter, i int) float64 {
	parse := func(name string) (float64, bool) {
		value, err := strconv.ParseFloat(strings.TrimPrefix(name, "#"), 64)
		if err != nil {
			return 0, false
		}
		return value, true
	}

	if value, ok := parse(chapters[i].GetName()); ok {
		return value
	}
	if i > 0 {
		if value, ok := parse(chapters[i-1].GetName()); ok {
			return value + 0.5
		}
	}
	if i < len(chapters)-1 {
		if value, ok := parse(chapters[i+1].GetName()); ok {
			return value - 0.5
		}
	}
	return 0
}

func formatNum(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func chapterTitle(num float64, subTitle string) string {
	match := chapterTitleRe.FindStringSubmatch(subTitle)
	if len(match) > 1 {
		return fmt.Sprintf("Chapter %s - %s", formatNum(num), match[1])
	}
	return subTitle
}

func formatChapters(protoChapters []*mpproto.Chapter) []manga.Chapter {
	chapters := make([]manga.Chapter, 0, len(protoChapters))
	for i, chapter := range protoChapters {
		num := parseChapterNum(protoChapters, i)
		numString := formatNum(num)
		title := chapterTitle(num, chapter.GetSubTitle())

		chapters = append(chapters, manga.Chapter{
			ID:        strconv.FormatUint(uint64(chapter.GetChapterId()), 10),
			SortNum:   num,
			RawTitle:  chapter.GetSubTitle(),
			FullTitle: title,
			Metadata: manga.Metadata{
				Title:    title,
				Num:      numString,
				Language: "en",
				Date:     "",
				Link:     fmt.Sprintf("https://mangaplus.shueisha.co.jp/viewer/%d", chapter.GetChapterId()),
				Groups:   []string{"MangaPlus"},
			},
		})
	}
	return chapters
}

func (m *Scraper) scrapeChapters() ([]manga.Chapter, *logging.ScraperError) {
	respBytes, _ := rester.New().GetBytes(
		fmt.Sprintf("%s/title_detailV3", APIROOT),
		m.headers(fmt.Sprintf("%s/titles/%s", REFERERURL, m.MangaID())),
		[]rester.QueryParam{{Key: "title_id", Value: m.MangaID()}},
	).Do(4, "100ms")
	data := respBytes.([]byte)

	resp, err := decodeResponse(data)
	if err != nil {
		return nil, &logging.ScraperError{
			Error:   err,
			Message: "Failed to decode title detail response",
			Code:    0,
		}
	}

	detail := resp.GetSuccess().GetTitleDetailView()
	if detail == nil {
		return nil, &logging.ScraperError{
			Error:   fmt.Errorf("titleDetailView is nil"),
			Message: "Failed to get manga info",
			Code:    0,
		}
	}

	m.manga.name = detail.GetTitle().GetName()

	var chapters []manga.Chapter
	for _, block := range detail.GetChapters() {
		if first := block.GetFirstChapterList(); len(first) > 0 {
			chapters = append(chapters, formatChapters(first)...)
		}
		if last := block.GetLastChapterList(); len(last) > 0 {
			chapters = append(chapters, formatChapters(last)...)
		}
	}

	for i, j := 0, len(chapters)-1; i < j; i, j = i+1, j-1 {
		chapters[i], chapters[j] = chapters[j], chapters[i]
	}

	return chapters, nil
}

func (m *Scraper) SelectChapters(titles []string) *logging.ScraperError {
	var selected []manga.Chapter
	for _, chapter := range m.allChapters {
		for _, title := range titles {
			if chapter.FullTitle == title {
				selected = append(selected, chapter)
				break
			}
		}
	}

	m.selectedChapters = selected
	m.allChapters = []manga.Chapter{}
	return nil
}

func (m *Scraper) SelectNewChapters(chapterIDs []string) ([]manga.Chapter, *logging.ScraperError) {
	chapters, err := m.Chapters()
	if err != nil {
		return []manga.Chapter{}, err
	}

	idSet := make(map[string]struct{}, len(chapterIDs))
	for _, chapterID := range chapterIDs {
		idSet[chapterID] = struct{}{}
	}

	var newChapters []manga.Chapter
	for _, chapter := range chapters {
		if _, exists := idSet[chapter.ID]; !exists {
			newChapters = append(newChapters, chapter)
		}
	}

	m.selectedChapters = newChapters
	logging.Debugln("SelectNewChapters: New chapters: ", newChapters)
	return newChapters, nil
}

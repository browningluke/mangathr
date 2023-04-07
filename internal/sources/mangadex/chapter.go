package mangadex

import (
	"encoding/json"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/manga"
	"github.com/browningluke/mangathrV2/internal/rester"
	"strconv"
	"strings"
)

const feedPageLimit = 500

/*
	-- Feed parsing --
*/

func getMangaFeedPage(id string, params []rester.QueryParam, offset int) (mangaFeedResponse, *logging.ScraperError) {
	jsonResp, _ := rester.New().Get(
		fmt.Sprintf("%s/manga/%s/feed", APIROOT, id),
		map[string]string{},
		append(params,
			rester.QueryParam{Key: "offset", Value: strconv.Itoa(offset), Encode: true},
			rester.QueryParam{Key: "includes[]", Value: "scanlation_group", Encode: true}),
	).Do(4, "100ms")
	jsonString := jsonResp.(string)

	var mangaFeedResp mangaFeedResponse

	err := json.Unmarshal([]byte(jsonString), &mangaFeedResp)
	if err != nil {
		return mangaFeedResponse{}, &logging.ScraperError{
			Error:   err,
			Message: "An error occurred while getting Manga chapters",
			Code:    0,
		}
	}

	return mangaFeedResp, nil
}

// getMangaFeed: handles pagination of Feed API endpoint
func getMangaFeed(mangaID string, languages, ratings []string) ([]mangaFeedResponse, *logging.ScraperError) {
	// Build query params
	queryParams := []rester.QueryParam{
		{Key: "limit", Value: fmt.Sprint(feedPageLimit), Encode: true},
		{Key: "order[chapter]", Value: "desc", Encode: true},
	}

	for _, language := range languages {
		queryParams = append(queryParams, rester.QueryParam{Key: "translatedLanguage[]", Value: language, Encode: true})
	}

	for _, rating := range ratings {
		queryParams = append(queryParams, rester.QueryParam{Key: "contentRating[]", Value: rating, Encode: true})
	}

	// Get all pages of feed

	var mangaFeedRespList []mangaFeedResponse
	initial, err := getMangaFeedPage(mangaID, queryParams, 0)
	if err != nil {
		return nil, err
	}

	mangaFeedRespList = append(mangaFeedRespList, initial)

	// Calculate total offset pages
	// this will be the floor of total/limit (1200/500 = 2)
	// but since we retrieved the first page already, it becomes 1 + (1200/500 = 2) = 1500 > 1200
	pages := initial.Total / feedPageLimit

	for i := 1; i <= pages; i++ {
		page, err := getMangaFeedPage(mangaID, queryParams, feedPageLimit*i)
		if err != nil {
			return nil, err
		}

		mangaFeedRespList = append(mangaFeedRespList, page)
	}

	return mangaFeedRespList, nil
}

/*
	-- Chapter parsing --
*/

// parseChapterNum: parses chapter number as both string and float
func parseChapterNum(chapterNum string) (string, float64, *logging.ScraperError) {
	var numFloat float64

	if chapterNum == "" {
		numFloat = 0
	} else {
		parsedFloat, err := strconv.ParseFloat(chapterNum, 64)
		numFloat = parsedFloat
		if err != nil {
			return "", 0.0, &logging.ScraperError{
				Error:   err,
				Message: "An error occurred while parsing chapter number",
				Code:    0,
			}
		}
	}

	// Extract number
	num := chapterNum
	if chapterNum == "" {
		num = "0"
	}

	return num, numFloat, nil
}

func (m *Scraper) parseGroups(data mangaFeedData) []string {
	var groups []string
	for _, relationship := range data.Relationships {
		if relationship.RelationType == "scanlation_group" {
			groups = append(groups, relationship.Attributes.Name)
		}
	}

	// Add groups to scraper
	for _, group := range groups {
		// Check if group already caught my scraper
		skip := false
		for _, a := range m.groups {
			if a == group {
				skip = true
				break
			}
		}

		// Mark it as caught for future
		if !skip {
			m.groups = append(m.groups, group)
		}
	}

	return groups
}

// generateTitle: returns fullTitle (including group), and metadata title (without group)
func (m *Scraper) generateTitle(chapterTitle, num, lang string, groups []string) (string, string) {
	// Generate title padding
	titlePadding := ""

	if len(config.LanguageFilter) > 1 {
		titlePadding += fmt.Sprintf(" - %s", lang)
	}

	if chapterTitle != "" {
		titlePadding += fmt.Sprintf(" - %s", chapterTitle)
	}

	metadataTitle := fmt.Sprintf("Chapter %s%s", num, titlePadding)

	fullTitle := metadataTitle
	if len(groups) > 0 {
		fullTitle += fmt.Sprintf(" [%s]", strings.Join(groups[:], ", "))
	}

	return fullTitle, metadataTitle
}

func (m *Scraper) scrapeChapters() ([]manga.Chapter, *logging.ScraperError) { // Get entire Manga feed
	mangaFeed, err := getMangaFeed(m.MangaID(), config.LanguageFilter, config.RatingFilter)
	if err != nil {
		return []manga.Chapter{}, err
	}

	var searchResults []manga.Chapter
	// For each page
	for _, mangaFeedResp := range mangaFeed {
		// For each chapter in page
		for _, item := range mangaFeedResp.Data {

			numString, numFloat, err := parseChapterNum(item.Attributes.Chapter)
			if err != nil {
				return []manga.Chapter{}, err
			}

			groups := m.parseGroups(item)

			fullTitle, metadataTitle := m.generateTitle(item.Attributes.Title, numString,
				item.Attributes.TranslatedLanguage, groups)

			searchResults = append(searchResults,
				manga.Chapter{
					ID:      item.Id,
					SortNum: numFloat,

					FullTitle: fullTitle,
					RawTitle:  item.Attributes.Title,

					Metadata: manga.Metadata{
						Title:    metadataTitle,
						Num:      numString,
						Language: item.Attributes.TranslatedLanguage,
						Date:     item.Attributes.CreatedAt[0:11],
						Link:     fmt.Sprintf("https://mangadex.org/chapter/%s", item.Id),
						Groups:   groups,
					},
				})
		}
	}

	searchResults = handleDuplicates(searchResults)
	return searchResults, nil
}

// handleDuplicates: returns a slice of structs.Chapter, with IDs appended to chapters with duplicate titles.
// Handles the rare case where a chapter has the same number, title AND group.
func handleDuplicates(chapters []manga.Chapter) []manga.Chapter {
	allKeys := make(map[string][]int)

	for i, item := range chapters {
		if indexArray, exists := allKeys[item.FullTitle]; exists {
			allKeys[item.FullTitle] = append(indexArray, i)
		} else {
			allKeys[item.FullTitle] = []int{i}
		}
	}

	for _, indexArray := range allKeys {
		if len(indexArray) > 1 {
			// If we're here, we know that there are >1 chapters with the same title
			for _, r := range indexArray {
				chapter := chapters[r]
				shortID := strings.Split(chapter.ID, "-")[0]
				chapter.FullTitle = chapter.FullTitle + " [" + shortID + "]"
				chapter.RawTitle = chapter.RawTitle + " [" + shortID + "]"
				chapters[r] = chapter
			}
		}
	}

	return chapters
}

/*
	-- Chapter selection --
*/

func (m *Scraper) SelectChapters(titles []string) *logging.ScraperError {
	var chapters []manga.Chapter

	for _, chapter := range m.allChapters {
		for _, promptTitle := range titles {
			if chapter.FullTitle == promptTitle {
				chapters = append(chapters, chapter)
			}
		}
	}
	m.selectedChapters = chapters

	// Once chapters have been selected, clear all chapters
	m.allChapters = []manga.Chapter{}

	return nil
}

func (m *Scraper) SelectNewChapters(chapterIDs []string) ([]manga.Chapter, *logging.ScraperError) {
	// Parse chapters if not already done
	chapters, err := m.Chapters()
	if err != nil {
		return []manga.Chapter{}, err
	}

	var diffChapters []manga.Chapter
	for _, newChapter := range chapters {
		exists := false
		for _, oldChapterID := range chapterIDs {
			if oldChapterID == newChapter.ID {
				exists = true
				break
			}
		}
		if !exists {
			diffChapters = append(diffChapters, newChapter)
		}
	}
	m.selectedChapters = diffChapters

	logging.Debugln("SelectNewChapters: New chapters: ", diffChapters)
	return diffChapters, nil
}

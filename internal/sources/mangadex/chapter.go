package mangadex

import (
	"encoding/json"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/rester"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"math"
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

	for i := 1; i <= int(math.Ceil(float64(initial.Total/feedPageLimit))); i++ {
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

	if len(m.config.LanguageFilter) > 1 {
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

func (m *Scraper) scrapeChapters() *logging.ScraperError {
	// Get entire Manga feed
	mangaFeed, err := getMangaFeed(m.MangaID(), m.config.LanguageFilter, m.config.RatingFilter)
	if err != nil {
		return err
	}

	var searchResults []chapterResult
	// For each page
	for _, mangaFeedResp := range mangaFeed {
		// For each chapter in page
		for _, item := range mangaFeedResp.Data {

			numString, numFloat, err := parseChapterNum(item.Attributes.Chapter)
			if err != nil {
				return err
			}

			groups := m.parseGroups(item)

			fullTitle, metadataTitle := m.generateTitle(item.Attributes.Title, numString,
				item.Attributes.TranslatedLanguage, groups)

			searchResults = append(searchResults,
				chapterResult{
					id: item.Id,

					fullTitle: fullTitle,
					title:     item.Attributes.Title,

					sortNum: numFloat,

					language: item.Attributes.TranslatedLanguage,

					metadata: structs.Metadata{
						Title:  metadataTitle,
						Num:    numString,
						Groups: groups,
						Link:   fmt.Sprintf("https://mangadex.org/chapter/%s", item.Id),
						Date:   item.Attributes.CreatedAt[0:11],
					},
				})
		}
	}

	m.allChapters = searchResults
	return nil
}

/*
	-- Chapter selection --
*/

func (m *Scraper) SelectChapters(titles []string) *logging.ScraperError {
	var chapters []chapterResult

	for _, chapter := range m.allChapters {
		for _, promptTitle := range titles {
			if chapter.fullTitle == promptTitle {
				chapters = append(chapters, chapter)
			}
		}
	}
	m.selectedChapters = chapters

	// Once chapters have been selected, clear all chapters
	m.allChapters = []chapterResult{}

	return nil
}

func (m *Scraper) SelectNewChapters(chapters []structs.Chapter) ([]structs.Chapter, *logging.ScraperError) {
	// Populate .allChapters
	if _, err := m.Chapters(); err != nil {
		return nil, err
	}

	var diffChapters []chapterResult
	for _, newChapter := range m.allChapters {
		exists := false
		for _, oldChapter := range chapters {
			if oldChapter.ID == newChapter.id {
				exists = true
				break
			}
		}
		if !exists {
			diffChapters = append(diffChapters, newChapter)
		}
	}
	m.selectedChapters = diffChapters
	m.allChapters = []chapterResult{}

	logging.Debugln("SelectNewChapters: New chapters: ", diffChapters)
	var diffStructChapters []structs.Chapter
	for _, chapter := range diffChapters {
		diffStructChapters = append(diffStructChapters, structs.Chapter{
			ID:       chapter.id,
			Title:    chapter.title,
			Metadata: chapter.metadata,
		})
	}

	return diffStructChapters, nil
}

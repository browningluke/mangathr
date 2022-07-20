package mangadex

import (
	"encoding/json"
	"errors"
	"fmt"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/logging"
	"mangathrV2/internal/rester"
	"mangathrV2/internal/sources/structs"
	"mangathrV2/internal/utils"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Scraper struct {
	name   string
	config *Config

	searchResults []searchResult
	manga         searchResult

	allChapters, selectedChapters []chapterResult

	// Group queries
	groups []string
}

type searchResult struct {
	title, id string
}

type chapterResult struct {
	id, title, prettyTitle, promptTitle,
	num, language string
	sortNum float64

	metadata structs.Metadata
}

func NewScraper(config *Config) *Scraper {
	logging.Debugln("Created a Mangadex scraper")
	s := &Scraper{config: config}
	return s
}

/*
	-- Searching --
*/

// Search for a Manga, will fill searchResults with 0 or more results
func (m *Scraper) Search(query string) []string {
	// Build query params
	queryParams := []rester.QueryParam{
		{Key: "order[relevance]", Value: "desc", Encode: true},
		{Key: "title", Value: query, Encode: true},
	}

	for _, rating := range m.config.RatingFilter {
		queryParams = append(queryParams, rester.QueryParam{Key: "contentRating[]", Value: rating, Encode: true})
	}

	jsonResp, _ := rester.New().Get(
		"https://api.mangadex.org/manga",
		map[string]string{},
		queryParams).Do(4, "100ms")
	jsonString := jsonResp.(string)

	var mangaResp mangaResponse
	err := json.Unmarshal([]byte(jsonString), &mangaResp)
	if err != nil {
		panic(err)
	}

	var searchResults []searchResult
	var names []string

	for _, item := range mangaResp.Data {
		name := item.Attributes.Title["en"]
		if name == "" {
			for _, n := range item.Attributes.Title {
				name = n
				break
			}
		}

		sr := searchResult{title: name, id: item.Id}

		searchResults = append(searchResults, sr)
		names = append(names, name)
	}

	m.searchResults = searchResults

	return names
}

// SearchByID for a Manga, will fill searchResults with ONLY 1 result (first result)
func (m *Scraper) SearchByID(id, title string) error {

	// Test if ID is valid
	_, resp := rester.New().Get(
		fmt.Sprintf("https://api.mangadex.org/manga/%s", id),
		map[string]string{},
		[]rester.QueryParam{}).Do(4, "100ms")

	if resp.StatusCode != 200 {
		return errors.New("SearchByID: validation status code != 200")
	}

	m.manga = searchResult{title: title, id: id}
	return nil
}

/*
	-- Chapter scraping --
*/

func (m *Scraper) scrapeChapters() {
	// Build query params
	queryParams := []rester.QueryParam{
		{Key: "limit", Value: "500", Encode: true},
		{Key: "order[chapter]", Value: "desc", Encode: true},
	}

	for _, language := range m.config.LanguageFilter {
		queryParams = append(queryParams, rester.QueryParam{Key: "translatedLanguage[]", Value: language, Encode: true})
	}

	for _, rating := range m.config.RatingFilter {
		queryParams = append(queryParams, rester.QueryParam{Key: "contentRating[]", Value: rating, Encode: true})
	}

	getMangaFeedResp := func(offset int) mangaFeedResponse {
		jsonResp, _ := rester.New().Get(
			fmt.Sprintf("https://api.mangadex.org/manga/%s/feed", m.manga.id),
			map[string]string{},
			append(queryParams,
				rester.QueryParam{Key: "offset", Value: strconv.Itoa(offset), Encode: true},
				rester.QueryParam{Key: "includes[]", Value: "scanlation_group", Encode: true}),
		).Do(4, "100ms")
		jsonString := jsonResp.(string)

		var mangaFeedResp mangaFeedResponse

		err := json.Unmarshal([]byte(jsonString), &mangaFeedResp)
		if err != nil {
			panic(err)
		}

		return mangaFeedResp
	}

	var mangaFeedRespList []mangaFeedResponse

	initial := getMangaFeedResp(0)
	mangaFeedRespList = append(mangaFeedRespList, initial)
	for i := 1; i <= int(math.Ceil(float64(initial.Total/500))); i++ {
		mangaFeedRespList = append(mangaFeedRespList, getMangaFeedResp(500*i))
	}

	var searchResults []chapterResult

	for _, mangaFeedResp := range mangaFeedRespList {
		for _, item := range mangaFeedResp.Data {
			var f float64

			if item.Attributes.Chapter == "" {
				f = 0
			} else {
				parsedFloat, err := strconv.ParseFloat(item.Attributes.Chapter, 64)
				f = parsedFloat
				if err != nil {
					logging.Errorln(item)
					panic(err)
				}
			}

			num := item.Attributes.Chapter
			if item.Attributes.Chapter == "" {
				num = "0"
			}

			// Extract all group info
			var groups []string
			for _, relationship := range item.Relationships {
				if relationship.RelationType == "scanlation_group" {
					groups = append(groups, relationship.Attributes.Name)
				}
			}

			// Add groups to scraper
			for _, group := range groups {
				skip := false
				for _, a := range m.groups {
					if a == group {
						skip = true
						break
					}
				}
				if !skip {
					m.groups = append(m.groups, group)
				}
			}

			// Generate title padding
			titlePadding := ""

			if len(m.config.LanguageFilter) > 1 {
				titlePadding += fmt.Sprintf(" - %s", item.Attributes.TranslatedLanguage)
			}

			if item.Attributes.Title != "" {
				titlePadding += fmt.Sprintf(" - %s", item.Attributes.Title)
			}

			prettyTitle := fmt.Sprintf("Chapter %s%s",
				num, titlePadding)

			promptTitle := prettyTitle
			if len(groups) > 0 {
				promptTitle += fmt.Sprintf(" [%s]", strings.Join(groups[:], ", "))
			}

			searchResults = append(searchResults,
				chapterResult{
					prettyTitle: prettyTitle,
					promptTitle: promptTitle,
					title:       item.Attributes.Title,
					num:         num,
					sortNum:     f,
					id:          item.Id,
					language:    item.Attributes.TranslatedLanguage,

					metadata: structs.Metadata{
						Groups: groups,
						Link:   fmt.Sprintf("https://mangadex.org/chapter/%s", item.Id),
						Date:   item.Attributes.CreatedAt[0:11],
					},
				})
		}
	}
	m.allChapters = searchResults
}

// SelectManga from searchResults list
func (m *Scraper) SelectManga(name string) {
	found := false
	for _, item := range m.searchResults {
		if item.title == name {
			m.manga = item
			found = true
			break
		}
	}

	if !found {
		panic("Selected manga not is search result list")
	}

	// Once manga has been selected, clear all search results
	m.searchResults = []searchResult{}
}

func (m *Scraper) SelectChapters(titles []string) {
	var chapters []chapterResult

	for _, chapter := range m.allChapters {
		for _, promptTitle := range titles {
			if chapter.promptTitle == promptTitle {
				chapters = append(chapters, chapter)
			}
		}
	}
	m.selectedChapters = chapters

	// Once chapters have been selected, clear all chapters
	m.allChapters = []chapterResult{}
}

func (m *Scraper) SelectNewChapters(chapters []structs.Chapter) []structs.Chapter {
	// Populate .allChapters
	_ = m.Chapters()

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
			Num:      chapter.num,
			Title:    chapter.title,
			Metadata: chapter.metadata,
		})
	}

	return diffStructChapters
}

/*
	-- Chapter data --
*/

// Getters

func (m *Scraper) Chapters() []structs.Chapter {
	if len(m.allChapters) == 0 && len(m.selectedChapters) == 0 {
		m.scrapeChapters()
	}

	c := m.allChapters

	if len(m.selectedChapters) != 0 {
		c = m.selectedChapters
	}

	var chapters []structs.Chapter
	for _, item := range c {
		chapters = append(chapters,
			structs.Chapter{ID: item.id, Title: item.prettyTitle, Num: item.num, Metadata: item.metadata})
	}
	return chapters
}

func (m *Scraper) ChapterTitles() []string {
	if len(m.allChapters) == 0 && len(m.selectedChapters) == 0 {
		m.scrapeChapters()
	}

	chapters := m.allChapters

	if len(m.selectedChapters) != 0 {
		chapters = m.selectedChapters
	}

	var titles []string
	for _, item := range chapters {
		titles = append(titles, item.promptTitle)
	}
	return titles
}

func (m *Scraper) GroupNames() []string {
	if len(m.groups) == 0 {
		m.Chapters()
	}

	var groupNames []string
	for _, val := range m.groups {
		groupNames = append(groupNames, val)
	}
	return groupNames
}

// Setters

func (m *Scraper) FilterGroups(groups []string) {
	findElemInSlice := func(slice []string, elem string) bool {
		for _, v := range slice {
			if elem == v {
				return true
			}
		}
		return false
	}

	var selectedChapters []chapterResult
	for _, chapter := range m.allChapters { // go through each chapter
		for _, group := range groups { // go through each filtered group
			exists := findElemInSlice(chapter.metadata.Groups, group)
			if exists {
				selectedChapters = append(selectedChapters, chapter)
				break
			}
		}
	}
	m.selectedChapters = selectedChapters
}

/*
	-- Pages & Download --
*/

func (m *Scraper) getChapterPages(id string) []downloader.Page {
	resInterface := rester.New().Get(
		fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", id),
		map[string]string{},
		[]rester.QueryParam{},
	).DoWithHelperFunc(4, "200ms", func(res rester.Response, err error) {
		logging.Errorln(res.StatusCode)

		if res.StatusCode == 429 {
			header := res.Headers

			nextRetryTimeInt := header["X-Ratelimit-Retry-After"][0]
			nextRetryTime, err := strconv.ParseInt(nextRetryTimeInt, 10, 64)
			if err != nil {
				panic(err)
			}

			now := time.Now().Unix()
			timeDiff := nextRetryTime - now
			logging.Warningln(fmt.Sprintf("Time right now: %d", now))
			logging.Warningln(fmt.Sprintf("Sleeping for: %d", timeDiff))

			time.Sleep(time.Duration(timeDiff) * time.Second)
		}
	})

	jsonString := resInterface.(string)

	var chapterResp chapterResponse

	err := json.Unmarshal([]byte(jsonString), &chapterResp)
	if err != nil {
		panic(err)
	}

	length := len(chapterResp.Chapter.Data)
	digits := int(math.Floor(math.Log10(float64(length))) + 1)

	logging.Debugln("Chapter pages: ", jsonString)

	getPages := func(slice []string, key string) []downloader.Page {
		var pages []downloader.Page
		for i, chapter := range slice {
			pages = append(pages, downloader.Page{
				Url: fmt.Sprintf("%s/%s/%s/%s",
					chapterResp.BaseUrl, key, chapterResp.Chapter.Hash, chapter),
				Filename: fmt.Sprintf("%s%s",
					utils.PadString(fmt.Sprintf("%d", i+1), digits),
					utils.GetImageExtension(chapter)),
			})
		}
		return pages
	}

	var pages []downloader.Page

	if m.config.DataSaver {
		pages = getPages(chapterResp.Chapter.DataSaver, "data-saver")
	} else {
		pages = getPages(chapterResp.Chapter.Data, "data")
	}

	return pages
}

func (m *Scraper) Download(dl *downloader.Downloader, downloadType string) {
	logging.Debugln("Downloading...")

	chaptersPerMinute := 60 // set from API docs
	duration := int64((chaptersPerMinute * 1000) / 40)
	if numChapters := len(m.selectedChapters); numChapters < chaptersPerMinute {
		duration = int64((numChapters * 1000) / 40)
	}
	dl.SetChapterDuration(duration)

	// downloadType is one of ["download", "update"]
	path := dl.CreateDirectory(m.manga.title, downloadType)
	downloadQueue := make([]downloader.Job, len(m.selectedChapters))

	maxRuneCount := 0 // Used for padding (e.g. Chapter 10 vs Chapter 10.5)
	for i, chapter := range m.selectedChapters {
		language := ""
		if len(m.config.LanguageFilter) > 1 {
			language = fmt.Sprintf("%s", chapter.language)
		}
		chapterFilename := dl.GetNameFromTemplate(m.config.FilenameTemplate,
			chapter.num, chapter.title, language, chapter.metadata.Groups)

		downloadQueue[i] = downloader.Job{
			Title: chapter.prettyTitle, Num: chapter.num, ID: chapter.id,
			Filename: chapterFilename, Metadata: chapter.metadata,
		}

		if runeCount := utf8.RuneCountInString(chapter.num); runeCount > maxRuneCount {
			maxRuneCount = runeCount
		}
	}

	runJob := func(job downloader.Job) {
		pages := m.getChapterPages(job.ID)

		progress := utils.CreateProgressBar(len(pages), maxRuneCount, job.Num)

		mdAgent := dl.MetadataAgent()
		(*mdAgent).SetTitle(job.Title)
		(*mdAgent).SetNum(job.Num)
		(*mdAgent).SetWebLink(job.Metadata.Link)
		(*mdAgent).SetDate(job.Metadata.Date)
		(*mdAgent).SetEditors(job.Metadata.Groups)
		dl.Download(path, job.Filename, pages, progress)

		fmt.Println("") // Create a new bar for each chapter
	}

	// Execute download queue, potential to add workerpool here later
	for _, job := range downloadQueue {
		runJob(job)
	}

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
	return "Mangadex"
}

func (m *Scraper) EnforceChapterLength() bool {
	return true
}

package mangadex

import (
	"encoding/json"
	"fmt"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/logging"
	"mangathrV2/internal/rester"
	"mangathrV2/internal/sources/structs"
	"mangathrV2/internal/utils"
	"mangathrV2/internal/utils/ui"
	"math"
	"strconv"
	"strings"
)

type Scraper struct {
	name   string
	config *Config

	searchResults []searchResult
	manga         searchResult

	allChapters, selectedChapters []chapterResult

	// Group queries
	groups map[string]string // Map groupID -> groupName
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
	logging.Infoln("Created a Mangadex scraper")
	s := &Scraper{config: config}
	s.groups = make(map[string]string)
	return s
}

/*
	-- Manga --
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

	jsonString := rester.New().Get(
		"https://api.mangadex.org/manga",
		map[string]string{},
		queryParams).Do(1).(string)

	var mangaResp mangaResponse
	err := json.Unmarshal([]byte(jsonString), &mangaResp)
	if err != nil {
		panic(err)
	}

	var searchResults []searchResult
	var names []string

	for _, item := range mangaResp.Data {
		searchResults = append(searchResults, searchResult{title: item.Attributes.Title["en"], id: item.Id})
		names = append(names, item.Attributes.Title["en"])
	}

	m.searchResults = searchResults

	return names
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

// SearchByID for a Manga, will fill searchResults with ONLY 1 result (first result)
func (m *Scraper) SearchByID(id string) interface{} {
	//TODO implement me
	panic("implement me")
}

/*
	-- Chapters --
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

	getMangaFeedResp := func(offset int) mangaFeedResponse {
		jsonString := rester.New().Get(
			fmt.Sprintf("https://api.mangadex.org/manga/%s/feed", m.manga.id),
			map[string]string{},
			append(queryParams,
				rester.QueryParam{Key: "offset", Value: strconv.Itoa(offset), Encode: true},
				rester.QueryParam{Key: "includes[]", Value: "scanlation_group", Encode: true}),
		).Do(1).(string)

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

func (m *Scraper) Chapters() []structs.Chapter {
	if len(m.allChapters) == 0 {
		m.scrapeChapters()
	}

	var chapters []structs.Chapter
	for _, item := range m.allChapters {
		chapters = append(chapters, structs.Chapter{
			ID: item.id, Title: item.prettyTitle, Num: item.num, Metadata: item.metadata})
	}
	return chapters
}

func (m *Scraper) ChapterTitles() []string {
	if len(m.allChapters) == 0 {
		m.scrapeChapters()
	}
	var titles []string
	for _, item := range m.allChapters {
		titles = append(titles, item.promptTitle)
	}
	return titles
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

/*
	-- Pages & Download --
*/

func (m *Scraper) getChapterPages(id string) []downloader.Page {
	jsonString := rester.New().Get(
		fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", id),
		map[string]string{},
		[]rester.QueryParam{}).Do(1).(string)

	var chapterResp chapterResponse

	err := json.Unmarshal([]byte(jsonString), &chapterResp)
	if err != nil {
		panic(err)
	}

	length := len(chapterResp.Chapter.Data)
	digits := int(math.Floor(math.Log10(float64(length))) + 1)

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
	// downloadType is one of ["download", "update"]
	path := dl.CreateDirectory(m.manga.title, downloadType)
	progress := ui.CreateProgress()

	downloadQueue := make([]downloader.Job, len(m.selectedChapters))

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
	}

	runJob := func(job downloader.Job) {
		pages := m.getChapterPages(job.ID)

		bar := ui.AddBar(progress, int64(len(pages)),
			fmt.Sprintf("Chapter %s", job.Num))

		mdAgent := dl.MetadataAgent()
		(*mdAgent).SetTitle(job.Title)
		(*mdAgent).SetNum(job.Num)
		(*mdAgent).SetWebLink(job.Metadata.Link)
		(*mdAgent).SetDate(job.Metadata.Date)
		(*mdAgent).SetEditors(job.Metadata.Groups)
		dl.Download(path, job.Filename, pages, bar)
	}

	// Execute download queue, potential to add workerpool here later
	for _, job := range downloadQueue {
		runJob(job)
	}

	progress.Wait()
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

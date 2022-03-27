package mangadex

import (
	"encoding/json"
	"fmt"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/rester"
	"mangathrV2/internal/utils"
	"math"
	"strconv"
)

type Scraper struct {
	name   string
	config *Config

	searchResults []searchResult
	allChapters   []chapterResult

	manga            searchResult
	selectedChapters []reader
}

type searchResult struct {
	title, id string
}

type chapterResult struct {
	id, title, prettyTitle,
	num, language string
	sortNum float64
}

type chapterResultByNum []chapterResult

func (n chapterResultByNum) Len() int           { return len(n) }
func (n chapterResultByNum) Less(i, j int) bool { return n[i].sortNum < n[j].sortNum }
func (n chapterResultByNum) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

type reader struct {
	chapterResult chapterResult
	pages         []downloader.Page
}

func NewScraper(config *Config) *Scraper {
	fmt.Println("Created a mangadex scraper")
	return &Scraper{config: config}
}

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
		queryParams)

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

	//fmt.Println(mangaResp)
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

func (m *Scraper) ListChapters() []string {
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
			append(queryParams, rester.QueryParam{Key: "offset", Value: strconv.Itoa(offset), Encode: true}),
		)

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
	var names []string

	for _, mangaFeedResp := range mangaFeedRespList {
		for _, item := range mangaFeedResp.Data {
			var f float64

			if item.Attributes.Chapter == "" {
				f = 0
			} else {
				parsedFloat, err := strconv.ParseFloat(item.Attributes.Chapter, 64)
				f = parsedFloat
				if err != nil {
					fmt.Println("Error: ", item)
					panic(err)
				}
			}

			num := item.Attributes.Chapter
			if item.Attributes.Chapter == "" {
				num = "0"
			}

			// Generate title padding
			titlePadding := ""

			if len(m.config.LanguageFilter) > 1 {
				titlePadding += fmt.Sprintf(" - %s", item.Attributes.TranslatedLanguage)
			}

			if item.Attributes.Title == "" {
				titlePadding += ""
			} else {
				titlePadding += fmt.Sprintf(" - %s", item.Attributes.Title)
			}

			prettyTitle := fmt.Sprintf("Chapter %s%s",
				num, titlePadding)

			searchResults = append(searchResults,
				chapterResult{
					prettyTitle: prettyTitle,
					title:       item.Attributes.Title,
					num:         num,
					sortNum:     f,
					id:          item.Id,
					language:    item.Attributes.TranslatedLanguage,
				})
		}
	}
	m.allChapters = searchResults

	for _, item := range searchResults {
		names = append(names, item.prettyTitle)
	}

	return names
}

func (m *Scraper) SelectChapters(titles []string) {
	var chapters []reader

	for _, chapter := range m.allChapters {
		for _, prettyTitle := range titles {
			if chapter.prettyTitle == prettyTitle {
				chapters = append(chapters, reader{
					chapterResult: chapter,
					pages:         m.getChapterPages(chapter.id),
				})
			}
		}
	}
	m.selectedChapters = chapters

	// Once chapters have been selected, clear all chapters
	m.allChapters = []chapterResult{}
}

func (m *Scraper) getChapterPages(id string) []downloader.Page {
	fmt.Println(id)
	jsonString := rester.New().Get(
		fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", id),
		map[string]string{},
		[]rester.QueryParam{})

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

func (m *Scraper) Download(downloader *downloader.Downloader, downloadType string) {
	// downloadType is one of ["download", "update"]
	_ = downloader.CreateDirectory(m.manga.title, downloadType)
	for _, chapter := range m.selectedChapters {
		language := ""
		if len(m.config.LanguageFilter) > 1 {
			language = fmt.Sprintf("%s", chapter.chapterResult.language)
		}

		_ = downloader.GetNameFromTemplate(m.config.FilenameTemplate,
			chapter.chapterResult.num, chapter.chapterResult.title, language)

	}
}

func (m *Scraper) GetMangaTitle() string {
	return m.manga.title
}

func (m *Scraper) GetScraperName() string {
	return "Mangadex"
}

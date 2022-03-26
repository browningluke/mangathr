package mangadex

import (
	"encoding/json"
	"fmt"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/rester"
	"math"
	"net/url"
	"strconv"
)

type Scraper struct {
	name string

	searchResults []searchResult
	allChapters   []chapterResult

	manga            searchResult
	selectedChapters []reader
}

type searchResult struct {
	title string
	id    string
}

type chapterResult struct {
	title       string
	prettyTitle string
	num         string // doing this will make it easier
	sortNum     float64
	id          string
}

type chapterResultByNum []chapterResult

func (n chapterResultByNum) Len() int           { return len(n) }
func (n chapterResultByNum) Less(i, j int) bool { return n[i].sortNum < n[j].sortNum }
func (n chapterResultByNum) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

type reader struct {
	chapterResult chapterResult
	pageURLs      []string
}

func NewScraper() *Scraper {
	fmt.Println("Created a mangadex scraper")
	return &Scraper{}
}

// Search for a Manga, will fill searchResults with 0 or more results
func (m *Scraper) Search(query string) []string {
	jsonString := rester.New().Get(
		"https://api.mangadex.org/manga?limit=10&order[relevance]=desc&title="+url.QueryEscape(query),
		map[string]string{})

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

	getMangaFeedResp := func(offset int) mangaFeedResponse {
		jsonString := rester.New().Get(
			fmt.Sprintf("https://api.mangadex.org/manga/%s/feed"+
				"?limit=500&translatedLanguage[]=en&order[chapter]=desc&offset=%d", m.manga.id, offset),
			map[string]string{})

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

			var titlePadding string
			if item.Attributes.Title == "" {
				titlePadding = ""
			} else {
				titlePadding = fmt.Sprintf(" - %s", item.Attributes.Title)
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
				})
		}
	}
	//	ceil(1831/500)
	//fmt.Println(mangaResp)
	//sort.Sort(chapterResultByNum(searchResults))
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
					pageURLs:      getChapterURLs(chapter.id),
				})
			}
		}
	}
	m.selectedChapters = chapters

	// Once chapters have been selected, clear all chapters
	m.allChapters = []chapterResult{}

	fmt.Println("Selected chapters: ", m.selectedChapters)
}

func getChapterURLs(id string) []string {
	return nil
}

func (m *Scraper) Download(downloader *downloader.Downloader) {
	//TODO implement me
	panic("implement me")
}

func (m *Scraper) GetMangaTitle() string {
	return m.manga.title
}

func (m *Scraper) GetScraperName() string {
	return "Mangadex"
}

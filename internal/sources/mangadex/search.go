package mangadex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/rester"
)

/*
	-- Search --
*/

func buildQueryParams(query string, contentRatings []string) []rester.QueryParam {
	queryParams := []rester.QueryParam{
		{Key: "order[relevance]", Value: "desc", Encode: true},
		{Key: "title", Value: query, Encode: true},
	}

	for _, rating := range contentRatings {
		queryParams = append(queryParams, rester.QueryParam{
			Key: "contentRating[]", Value: rating, Encode: true,
		})
	}

	return queryParams
}

func parseSearchResults(mangaResp mangaResponse) ([]searchResult, []string) {
	var searchResults []searchResult
	var names []string

	for _, item := range mangaResp.Data {
		// Default to using English name
		name := item.Attributes.Title["en"]

		// Use next available name, if English is not available
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

	return searchResults, names
}

// Search for a Manga, will fill searchResults with 0 or more results
func (m *Scraper) Search(query string) []string {
	// Build query params
	queryParams := buildQueryParams(query, m.config.RatingFilter)

	// Search for list of Manga
	jsonResp, _ := rester.New().Get(
		fmt.Sprintf("%s/manga", APIROOT),
		map[string]string{},
		queryParams).Do(4, "100ms")
	jsonString := jsonResp.(string)

	var mangaResp mangaResponse
	err := json.Unmarshal([]byte(jsonString), &mangaResp)
	if err != nil {
		panic(err)
	}

	// Parse a list of Manga and a list of their names
	searchResults, names := parseSearchResults(mangaResp)
	m.searchResults = searchResults
	return names
}

/*
	-- SearchByID --
*/

// SearchByID for a Manga, will fill searchResults with ONLY 1 result (first result)
func (m *Scraper) SearchByID(id, title string) error {

	// Test if ID is valid
	_, resp := rester.New().Get(
		fmt.Sprintf("%s/manga/%s", APIROOT, id),
		map[string]string{},
		[]rester.QueryParam{}).Do(4, "100ms")

	if resp.StatusCode != 200 {
		return errors.New("SearchByID: validation status code != 200")
	}

	m.manga = searchResult{title: title, id: id}
	return nil
}

/*
	-- SelectManga --
*/

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

package mangadex

import (
	"encoding/json"
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/rester"
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
func (m *Scraper) Search(query string) ([]string, *logging.ScraperError) {
	// Build query params
	queryParams := buildQueryParams(query, config.RatingFilter)

	// Search for list of Manga
	jsonResp, _ := rester.New().Get(
		fmt.Sprintf("%s/manga", APIROOT),
		map[string]string{},
		queryParams).Do(4, "100ms")
	jsonString := jsonResp.(string)

	var mangaResp mangaResponse
	err := json.Unmarshal([]byte(jsonString), &mangaResp)
	if err != nil {
		return nil, &logging.ScraperError{
			Error: err, Message: "Failed to read data when searching", Code: 0,
		}
	}

	// Parse a list of Manga and a list of their names
	searchResults, names := parseSearchResults(mangaResp)
	m.searchResults = searchResults
	return names, nil
}

/*
	-- SearchByID --
*/

// SearchByID for a Manga, will fill searchResults with ONLY 1 result (first result)
func (m *Scraper) SearchByID(id, title string) *logging.ScraperError {

	// Test if ID is valid
	_, resp := rester.New().Get(
		fmt.Sprintf("%s/manga/%s", APIROOT, id),
		map[string]string{},
		[]rester.QueryParam{}).Do(4, "100ms")

	if resp.StatusCode != 200 {
		return &logging.ScraperError{
			Error:   fmt.Errorf("searching by ID failed with status code: %d", resp.StatusCode),
			Message: "Search for chapter returned non-200 code",
			Code:    0,
		}
	}

	m.manga = searchResult{title: title, id: id}
	return nil
}

/*
	-- SelectManga --
*/

// SelectManga from searchResults list
func (m *Scraper) SelectManga(name string) *logging.ScraperError {
	found := false
	for _, item := range m.searchResults {
		if item.title == name {
			m.manga = item
			found = true
			break
		}
	}

	if !found {
		return &logging.ScraperError{
			Error:   fmt.Errorf("selected manga `%s`not in searchResults", name),
			Message: "An error occurred while selected Manga",
			Code:    0,
		}
	}

	// Once manga has been selected, clear all search results
	m.searchResults = []searchResult{}
	return nil
}

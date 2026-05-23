package mangaplus

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/rester"
	mpproto "github.com/browningluke/mangathr/v2/internal/sources/mangaplus/proto"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
)

func decodeResponse(data []byte) (*mpproto.Response, error) {
	resp := &mpproto.Response{}
	if err := proto.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func fuzzySearch(query, title string) bool {
	return strings.Contains(strings.ToLower(title), strings.ToLower(query))
}

func (m *Scraper) Search(query string) ([]string, *logging.ScraperError) {
	respBytes, _ := rester.New().GetBytes(
		fmt.Sprintf("%s/title_list/allV2", APIROOT),
		m.headers(fmt.Sprintf("%s/manga_list/all", REFERERURL)),
		[]rester.QueryParam{},
	).Do(4, "100ms")
	data := respBytes.([]byte)

	resp, err := decodeResponse(data)
	if err != nil {
		return nil, &logging.ScraperError{
			Error:   err,
			Message: "Failed to decode all titles response",
			Code:    0,
		}
	}

	allTitlesView := resp.GetSuccess().GetAllTitlesViewV2()
	if allTitlesView == nil {
		return nil, &logging.ScraperError{
			Error:   fmt.Errorf("allTitlesViewV2 is nil"),
			Message: "Failed to get all manga data",
			Code:    0,
		}
	}

	var results []searchResult
	var names []string
	for _, group := range allTitlesView.GetAllTitlesGroup() {
		if !fuzzySearch(query, group.GetTheTitle()) {
			continue
		}
		for _, title := range group.GetTitles() {
			if int(title.GetLanguage()) == config.Language {
				result := searchResult{titleId: title.GetTitleId(), name: title.GetName()}
				results = append(results, result)
				names = append(names, title.GetName())
				break
			}
		}
	}

	m.searchResults = results
	return names, nil
}

func (m *Scraper) SearchByID(id, title string) *logging.ScraperError {
	titleId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return &logging.ScraperError{
			Error:   err,
			Message: "Failed to parse MangaPlus title ID",
			Code:    0,
		}
	}

	m.manga = searchResult{titleId: uint32(titleId), name: title}
	return nil
}

func (m *Scraper) SelectManga(name string) *logging.ScraperError {
	for _, result := range m.searchResults {
		if result.name == name {
			m.manga = result
			m.searchResults = []searchResult{}
			return nil
		}
	}

	return &logging.ScraperError{
		Error:   fmt.Errorf("selected manga %q not in search results", name),
		Message: "An error occurred while selecting manga",
		Code:    0,
	}
}

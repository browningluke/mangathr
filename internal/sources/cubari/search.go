package cubari

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/rester"
	"regexp"
	"strings"
)

// searchBySlug for a Manga, should match 1 single manga
func (m *Scraper) searchBySlug(slug string) ([]string, *logging.ScraperError) {
	jsonResp, _ := rester.New().Get(fmt.Sprintf("%s/%s/series/%s", APIROOT, m.provider.name, slug),
		map[string]string{}, []rester.QueryParam{}).Do(4, "100ms")
	jsonString := jsonResp.(string)

	err := json.Unmarshal([]byte(jsonString), &m.manga)
	if err != nil {
		return []string{}, &logging.ScraperError{
			Error:   err,
			Message: "An error occurred while getting Manga chapters",
			Code:    0,
		}
	}

	return []string{m.MangaTitle()}, nil
}

// matchQuery against different provider regex patters.
func (m *Scraper) matchQuery(query, reStr string, provider Provider) string {
	re := regexp.MustCompile(reStr)
	match := re.FindStringSubmatch(query)

	if len(match) > 0 {
		m.provider = provider
		return match[4]
	}
	return ""
}

// Search for a Manga, assumes query is a supported provider URL
func (m *Scraper) Search(query string) ([]string, *logging.ScraperError) {
	var slug string

	// Match query against providers
	for _, provider := range PROVIDERBYSTR {
		// Skip providers without a regex pattern
		if provider.regex == "" {
			continue
		}

		if s := m.matchQuery(query, provider.regex, provider); s != "" {
			slug = s
		}
	}

	// cubari fallback
	cbRe := regexp.MustCompile(`^((http(s)?://)?cubari\.moe/read/(.*?)/)(.*?)/?$`)
	cbMatch := cbRe.FindStringSubmatch(query)

	if len(cbMatch) > 0 {
		if cbMatch[4] == "mangadex" {
			return []string{}, &logging.ScraperError{
				Error:   errors.New("mangadex provider not supported by cubari source"),
				Message: "Cubari source does not support MD links. Please use MD source.",
				Code:    0,
			}
		}

		m.provider = PROVIDERBYSTR[cbMatch[4]]
		slug = cbMatch[5]
	}

	if slug == "" {
		return []string{}, &logging.ScraperError{
			Error:   errors.New("unable to parse query"),
			Message: "Cubari did not recognize the entered query.",
			Code:    0,
		}
	}

	return m.searchBySlug(slug)
}

// SearchByID for a Manga, assumed that it will find a (and only 1) match
func (m *Scraper) SearchByID(id, _ string) *logging.ScraperError {
	ids := strings.Split(id, "~")

	m.provider = PROVIDERBYSTR[ids[0]]
	_, err := m.searchBySlug(ids[1])
	return err
}

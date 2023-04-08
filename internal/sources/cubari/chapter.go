package cubari

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/manga"
	"github.com/browningluke/mangathr/internal/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SelectNewChapters from allChapters that are not in the DB
func (m *Scraper) SelectNewChapters(chapterIDs []string) ([]manga.Chapter, *logging.ScraperError) {
	// Parse chapters if not already done
	chapters, err := m.Chapters()
	if err != nil {
		return []manga.Chapter{}, err
	}

	var newChapters []manga.Chapter

	for _, v := range chapters {
		if _, ok := utils.FindInSlice(chapterIDs, v.ID); !ok {
			newChapters = append(newChapters, v)
		}
	}

	m.selectedChapters = newChapters

	logging.Debugln("SelectNewChapters: New chapters: ", newChapters)
	return newChapters, nil
}

// SelectChapters from allChapters that are to be downloaded
func (m *Scraper) SelectChapters(titles []string) *logging.ScraperError {
	var selectedChaps []manga.Chapter

	for _, v := range m.allChapters {
		_, ok := utils.FindInSlice(titles, v.FullTitle)
		if ok {
			selectedChaps = append(selectedChaps, v)
		}
	}

	m.selectedChapters = selectedChaps

	// Once chapters have been selected, clear all chapters
	m.allChapters = []manga.Chapter{}

	return nil
}

/*
	-- Extracting Chapter data from Cubari --
*/

// parseImgurStyle to get URLs from a raw JSON message
func parseImgurStyle(c json.RawMessage) (urls []string, ok bool) {
	var chapter []imgurPage
	err := json.Unmarshal(c, &chapter)

	if err != nil {
		return []string{}, false
	}

	var pageURLs []string
	for _, v := range chapter {
		pageURLs = append(pageURLs, v.Source)
	}

	return pageURLs, true
}

// parseListStyle to get URLs from a raw JSON message
func parseListStyle(c json.RawMessage) (urls []string, ok bool) {
	var chapter []string
	err := json.Unmarshal(c, &chapter)

	if err != nil {
		return []string{}, false
	}

	return chapter, true
}

// parseProxyStyle to get URLs from a raw JSON message
func parseProxyStyle(c json.RawMessage) (urls []string, ok bool) {
	var proxyURL string
	err := json.Unmarshal(c, &proxyURL)

	if err != nil {
		return []string{}, false
	}

	return []string{proxyURL}, true
}

// Chapters returns chapter data from Cubari's API
func (m *Scraper) parseChapters() ([]manga.Chapter, *logging.ScraperError) {
	var allChaps []manga.Chapter

	// Grab all chapters
	for chapterNum, v := range m.manga.Chapters {
		// Build url list
		groupData := make(map[string]struct {
			URLs        []string
			ReleaseDate string
		})

		for groupNum, chapterData := range v.Groups {
			// Extract URLs
			var urls []string

			if parsedURLs, ok := parseImgurStyle(chapterData); ok {
				// Try convert using imgur style
				urls = parsedURLs
			} else if parsedURLs, ok := parseListStyle(chapterData); ok {
				// Try list of strings style
				urls = parsedURLs
			} else if parsedURLs, ok := parseProxyStyle(chapterData); ok {
				// Try proxy style
				urls = parsedURLs
			} else {
				return []manga.Chapter{}, &logging.ScraperError{
					Error:   errors.New("unable to extract page data for chapter"),
					Message: "An error occurred while getting chapter pages",
					Code:    0,
				}
			}

			// Try extract release_date
			releaseDate := ""
			if v.ReleaseDate != nil {
				var releaseDateTime time.Time
				chapterReleaseDate := v.ReleaseDate[groupNum]

				if strRD, ok := chapterReleaseDate.(string); ok {
					i, err := strconv.ParseInt(strRD, 10, 64)
					if err != nil {
						return []manga.Chapter{}, &logging.ScraperError{
							Error:   err,
							Message: "An error occurred while getting chapters",
							Code:    0,
						}
					}

					releaseDateTime = time.Unix(i, 0)
				} else if intRD, ok := chapterReleaseDate.(int); ok {
					releaseDateTime = time.Unix(int64(intRD), 0)
				} else if fltRD, ok := chapterReleaseDate.(float64); ok {
					releaseDateTime = time.Unix(int64(fltRD), 0)
				}

				releaseDate = fmt.Sprintf("%d-%02d-%02d",
					releaseDateTime.Year(), releaseDateTime.Month(), releaseDateTime.Day())
			}

			groupData[groupNum] = struct {
				URLs        []string
				ReleaseDate string
			}{URLs: urls, ReleaseDate: releaseDate}
		}

		for groupNum, data := range groupData {
			// Build chapter info
			chapterID := fmt.Sprintf("%s~%s~%s~%s", m.provider.name, m.manga.Slug, chapterNum, groupNum)
			title := ""
			fullTitle := fmt.Sprintf("Chapter %s", chapterNum)
			if m.provider == GIST {
				if v.Title != "" {
					title = v.Title
					fullTitle += " - " + v.Title
				}
			}

			// Add group tag to chapter title
			fullTitle += fmt.Sprintf(" [%s]", m.manga.Groups[groupNum])

			sortNum, err := strconv.ParseFloat(chapterNum, 64)
			if err != nil {
				return []manga.Chapter{}, &logging.ScraperError{
					Error:   err,
					Message: "An error occurred while getting chapters",
					Code:    0,
				}
			}

			link := fmt.Sprintf("%s/%s", m.provider.sourceURL, m.manga.Slug)
			if m.provider == GIST {
				link += fmt.Sprintf("/%s/%s", v.Volume, chapterNum)
			}

			groups := []string{strings.ReplaceAll(m.manga.Groups[groupNum], "/", ", ")}

			allChaps = append(allChaps, manga.Chapter{
				ID:        chapterID,
				SortNum:   sortNum,
				RawTitle:  title,
				FullTitle: fullTitle,
				Metadata: manga.Metadata{
					Title:    title,
					Num:      chapterNum,
					Language: "",
					Date:     data.ReleaseDate,
					Link:     link,
					Groups:   groups,
				},
			})
			m.pages[chapterID] = data.URLs
		}
	}

	// Sort chapters by descending SortNum
	sort.SliceStable(allChaps, func(i, j int) bool {
		return allChaps[i].SortNum > allChaps[j].SortNum
	})

	return allChaps, nil
}

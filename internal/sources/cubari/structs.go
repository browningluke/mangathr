package cubari

import "encoding/json"

type mangaResponse struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	SeriesName  string `json:"series_name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Artist      string `json:"artist"`
	Cover       string `json:"cover"`

	Groups map[string]string `json:"groups"` // "1": "Example Group/Another Group"

	Chapters map[string]struct {
		Volume      string                     `json:"volume"`
		Title       string                     `json:"title"`
		Groups      map[string]json.RawMessage `json:"groups"` // either 'string' or '[]imgurChapter'
		ReleaseDate map[string]interface{}     `json:"release_date,omitempty"`
	} `json:"chapters"`
}

type imgurPage struct {
	Description string `json:"description"`
	Source      string `json:"src"`
}

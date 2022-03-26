package mangadex

type mangaResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`

	Data []struct {
		Id   string `json:"id"`
		Type string `json:"type"`

		Attributes struct {
			Title map[string]string `json:"title"`
		} `json:"attributes"`
		//Relationships interface{} `json:"relationships"`
	} `json:"data"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type mangaFeedResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`

	Data []struct {
		Id   string `json:"id"`
		Type string `json:"type"`

		Attributes struct {
			Title              string `json:"title"`
			Chapter            string `json:"chapter"`
			TranslatedLanguage string `json:"translatedLanguage"`
		} `json:"attributes"`
	} `json:"data"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type chapterResponse struct {
	Result  string `json:"result"`
	BaseUrl string `json:"baseUrl"`
	Chapter struct {
		Hash      string   `json:"hash"`
		Data      []string `json:"data"`
		DataSaver []string `json:"dataSaver"`
	} `json:"chapter"`
}

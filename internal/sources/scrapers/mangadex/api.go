package mangadex

type mangaResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`

	Data []manga `json:"data"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type manga struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Title map[string]string `json:"title"`
	} `json:"attributes"`
	//Relationships interface{} `json:"relationships"`
}

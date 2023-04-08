package manga

type Chapter struct {
	ID        string  // Unique identifier for the chapter (can be different per source)
	SortNum   float64 // Number used for sorting chapters
	RawTitle  string  // Title straight from source (without Chapter xx)
	FullTitle string  // Title including groups/language/etc
	Metadata  Metadata

	pages    []Page
	filename string
}

// Pages

func (c *Chapter) AddPage(url, name string) {
	c.pages = append(c.pages, Page{
		Url:  url,
		Name: name,
	})
}

func (c *Chapter) Pages() []Page {
	return c.pages
}

// Filename

func (c *Chapter) SetFilename(name string) {
	c.filename = name
}

func (c *Chapter) Filename() string {
	return c.filename
}

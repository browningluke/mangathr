package structs

type Chapter struct {
	ID        string  // Unique identifier for the chapter (can be different per source)
	SortNum   float64 // Number used for sorting chapters
	RawTitle  string  // Title straight from source (without Chapter xx)
	FullTitle string  // Title including groups/language/etc
	Metadata  Metadata
}

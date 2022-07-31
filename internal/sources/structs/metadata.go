package structs

type Metadata struct {
	Title    string   // Title seen by user in reader (not filename)
	Num      string   // Num string for chapter (eg. 20, 20.5)
	Language string   // Language of chapter (will be converted to, but NOT currently, ISO format)
	Date     string   // Date of release
	Link     string   // Link to chapter
	Groups   []string // Groups that release/scanlated chapter
}

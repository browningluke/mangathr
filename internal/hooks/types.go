package hooks

// Event name constants.
const (
	EventDownloadChapter = "download.chapter"
	EventUpdateSuccess   = "update.success"
	EventUpdateError     = "update.error"
)

// MangaContext holds per-manga template data.
type MangaContext struct {
	Title  string
	Source string
}

// ChapterContext holds per-chapter template data.
type ChapterContext struct {
	Num    string
	Title  string
	Path   string // filesystem path where the chapter was saved
	Lang   string
	Groups string // comma-separated group names
	Count  int    // total chapters downloaded in this run (1 for single-event hooks)
}

// ErrorContext holds error template data; nil when there is no error.
type ErrorContext struct {
	Message string
}

// HookContext is the template context for non-aggregate hook invocations.
type HookContext struct {
	Manga   MangaContext
	Chapter ChapterContext
	Error   *ErrorContext // nil on success
	Event   string
}

// AggregateHookContext is the template context when aggregate: true.
type AggregateHookContext struct {
	Items        []HookContext
	ChapterCount int
	ErrorCount   int
	MangaCount   int
	Event        string
}

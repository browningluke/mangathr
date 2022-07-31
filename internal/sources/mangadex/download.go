package mangadex

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/utils"
	"unicode/utf8"
)

// Download selected chapters. Does not return an error, thus it must handle errors itself
func (m *Scraper) Download(dl *downloader.Downloader, downloadType string) {
	logging.Debugln("Downloading...")

	// 60 seconds / CHAPTERSPERMIN = x = seconds per chapter
	// x * 1000 = milliseconds per chapter

	duration := int64((60 / CHAPTERSPERMIN) * 1000)
	if numChapters := len(m.selectedChapters); numChapters < CHAPTERSPERMIN {
		// Not going to exceed limit during batch, duration doesn't matter
		duration = int64(500)
	}
	dl.SetChapterDuration(duration)

	// downloadType is one of ["download", "update"]
	path := dl.CreateDirectory(m.manga.title, downloadType)
	downloadQueue := make([]downloader.Job, len(m.selectedChapters))

	maxRuneCount := 0 // Used for padding (e.g. Chapter 10 vs Chapter 10.5)
	for i, chapter := range m.selectedChapters {
		language := ""
		if len(m.config.LanguageFilter) > 1 {
			language = fmt.Sprintf("%s", chapter.language)
		}
		chapterFilename := dl.GetNameFromTemplate(m.config.FilenameTemplate,
			chapter.metadata.Num, chapter.title, language, chapter.metadata.Groups)

		downloadQueue[i] = downloader.Job{
			ID: chapter.id, Filename: chapterFilename, Metadata: chapter.metadata,
		}

		if runeCount := utf8.RuneCountInString(chapter.metadata.Num); runeCount > maxRuneCount {
			maxRuneCount = runeCount
		}
	}

	runJob := func(job downloader.Job) *logging.ScraperError {
		pages, err := m.getChapterPages(job.ID)
		if err != nil {
			return err
		}

		progress := utils.CreateProgressBar(len(pages), maxRuneCount, job.Metadata.Num)

		// Set MetadataAgent values
		(*dl.MetadataAgent()).
			SetTitle(job.Metadata.Title).
			SetNum(job.Metadata.Num).
			SetWebLink(job.Metadata.Link).
			SetDate(job.Metadata.Date).
			SetEditors(job.Metadata.Groups).
			SetPageCount(len(pages))
		dl.Download(path, job.Filename, pages, progress)

		fmt.Println("") // Create a new bar for each chapter
		return nil
	}

	// Execute download queue, potential to add workerpool here later
	for _, job := range downloadQueue {
		err := runJob(job)

		// Print error to screen, abandon chapter, and continue
		if err != nil {
			logging.Errorln(err.Error)
			fmt.Printf("%s. Skipping chapter...", err.Message)
		}
	}

}

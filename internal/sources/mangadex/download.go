package mangadex

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
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

	var downloadQueue []downloader.Job
	maxRuneCount := 0 // Used for padding (e.g. Chapter 10 vs Chapter 10.5)
	for _, chapter := range m.selectedChapters {
		downloadQueue = append(downloadQueue, downloader.Job{Chapter: chapter})

		// Check if string length is max in list
		if runeCount := utf8.RuneCountInString(chapter.Metadata.Num); runeCount > maxRuneCount {
			maxRuneCount = runeCount
		}
	}

	runJob := func(job downloader.Job) *logging.ScraperError {
		// Get chapter pages
		pages, err := m.getChapterPages(job.Chapter.ID)
		if err != nil {
			return err
		}

		metadata := job.Chapter.Metadata

		// Initialize progress bar
		progress := utils.CreateProgressBar(len(pages), maxRuneCount, job.Chapter.Metadata.Num)

		// Get chapter filename
		// todo handle language (and user template)
		filename := dl.GetNameFromTemplate(m.config.FilenameTemplate,
			metadata.Num, metadata.Title, metadata.Language, metadata.Groups)

		// Set MetadataAgent values
		(*dl.MetadataAgent()).
			SetFromStruct(job.Chapter.Metadata).
			SetPageCount(len(pages))
		dl.Download(path, filename, pages, progress)

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

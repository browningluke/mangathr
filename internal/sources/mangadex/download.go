package mangadex

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
)

func calculateDuration(numChapters int) int64 {
	// 60 seconds / CHAPTERSPERMIN = x = seconds per chapter
	// x * 1000 = milliseconds per chapter

	duration := int64((60 / CHAPTERSPERMIN) * 1000)
	if numChapters < CHAPTERSPERMIN {
		// Not going to exceed limit during batch, duration doesn't matter
		duration = int64(500)
	}
	return duration
}

func (m *Scraper) runDownloadJob(dl *downloader.Downloader, chapter *manga.Chapter) *logging.ScraperError {
	// Load chapter pages from API
	err := m.addPagesToChapter(chapter)
	if err != nil {
		return err
	}

	// Set chapter filename from template
	dl.SetTemplate(config.FilenameTemplate)
	chapter.SetFilename(dl.GetNameFromTemplate(chapter, m.MangaTitle(), SCRAPERNAME))

	// Set MetadataAgent values
	(*dl.MetadataAgent()).
		SetFromStruct(chapter.Metadata).
		SetPageCount(len(chapter.Pages()))

	// Check if download is possible
	err = dl.CanDownload(chapter)
	if err != nil {
		return err
	}

	downloadErr := dl.Download(chapter)
	if downloadErr != nil {
		return &logging.ScraperError{
			Error:   downloadErr,
			Message: "An error occurred while downloading a page",
			Code:    0,
		}
	}

	return nil
}

// Download selected chapters. Handles errors itself. Returns array of chapters that succeeded
func (m *Scraper) Download(dl *downloader.Downloader, directoryMapping string) []manga.Chapter {
	logging.Debugln("Downloading...")

	directoryName := m.manga.title
	if directoryMapping != "" {
		directoryName = directoryMapping
	}

	// Configure downloader (downloadType is one of ["download", "update"])
	dl.SetChapterDuration(calculateDuration(len(m.selectedChapters))).
		SetPath(dl.CreateDirectory(directoryName)).
		SetMaxRuneCount(m.selectedChapters)

	// Execute download queue, potential to add workerpool here later
	var succeededChapters []manga.Chapter
	for _, chapter := range m.selectedChapters {
		err := m.runDownloadJob(dl, &chapter)

		// Print error to screen, abandon chapter, and continue
		if err != nil {
			logging.Errorln(err.Error)
			fmt.Printf("Chapter %s skipping... reason: %s\n", chapter.Metadata.Num, err.Message)
			continue
		}

		succeededChapters = append(succeededChapters, chapter)
	}

	return succeededChapters
}

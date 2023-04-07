package mangadex

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/manga"
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

	// Get chapter pages
	err := m.addPagesToChapter(chapter)
	if err != nil {
		return err
	}

	// Get chapter filename
	dl.SetTemplate(config.FilenameTemplate)
	filename := dl.GetNameFromTemplate(chapter)

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
		if err := dl.Cleanup(chapter); err != nil {
			logging.Errorln(err)
			fmt.Printf("An error occurred when deleting failed chapter: %s", filename)
		}
		return &logging.ScraperError{
			Error:   downloadErr,
			Message: "An error occurred while downloading a page",
			Code:    0,
		}
	}

	return nil
}

// Download selected chapters. Handles errors itself. Returns array of chapters that succeeded
func (m *Scraper) Download(dl *downloader.Downloader, directoryMapping, downloadType string) []manga.Chapter {
	logging.Debugln("Downloading...")

	dl.SetChapterDuration(calculateDuration(len(m.selectedChapters)))

	directoryName := m.manga.title
	if directoryMapping != "" {
		directoryName = directoryMapping
	}

	// Configure downloader (downloadType is one of ["download", "update"])
	dl.SetPath(dl.CreateDirectory(directoryName, downloadType))
	dl.SetMaxRuneCount(m.selectedChapters)

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

package cubari

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
)

func (m *Scraper) runDownloadJob(dl *downloader.Downloader, chapter *manga.Chapter) *logging.ScraperError {
	// Load chapter pages
	err := m.addPagesToChapter(chapter)
	if err != nil {
		return err
	}

	// Get chapter filename
	dl.SetTemplate(config.FilenameTemplate)
	chapter.SetFilename(dl.GetNameFromTemplate(chapter))

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

func (m *Scraper) Download(dl *downloader.Downloader, directoryMapping string) []manga.Chapter {
	logging.Debugln("Downloading...")

	directoryName := m.manga.Title
	if directoryMapping != "" {
		directoryName = directoryMapping
	}

	// Configure downloader
	dl.SetPath(dl.CreateDirectory(directoryName)).
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

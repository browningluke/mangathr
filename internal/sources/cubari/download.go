package cubari

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/downloader"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/sources/structs"
	"github.com/browningluke/mangathr/internal/utils"
)

func (m *Scraper) runDownloadJob(job downloader.Job, dl *downloader.Downloader,
	path string, maxRuneCount int) *logging.ScraperError {

	// Get chapter pages
	pages, err := m.getChapterPages(job.Chapter.ID)
	if err != nil {
		return err
	}

	// Initialize progress bar
	progress := utils.CreateProgressBar(len(pages), maxRuneCount, job.Chapter.Metadata.Num)

	// Get chapter filename
	dl.SetTemplate("") // Use default
	filename := dl.GetNameFromTemplate(job)

	// Set MetadataAgent values
	(*dl.MetadataAgent()).
		SetFromStruct(job.Chapter.Metadata).
		SetPageCount(len(pages))

	// Check if download is possible
	err = dl.CanDownload(path, filename)
	if err != nil {
		return err
	}

	downloadErr := dl.Download(path, filename, pages, progress)
	if downloadErr != nil {
		if err := dl.Cleanup(path, filename); err != nil {
			logging.Errorln(err)
			fmt.Printf("An error occurred when deleting failed chapter: %s", filename)
		}
		return &logging.ScraperError{
			Error:   downloadErr,
			Message: "An error occurred while downloading a page",
			Code:    0,
		}
	}

	fmt.Println("") // Create a new bar for each chapter
	return nil
}

func (m *Scraper) Download(dl *downloader.Downloader, directoryMapping, downloadType string) []structs.Chapter {
	logging.Debugln("Downloading...")

	directoryName := m.manga.Title
	if directoryMapping != "" {
		directoryName = directoryMapping
	}
	// downloadType is one of ["download", "update"]
	path := dl.CreateDirectory(directoryName, downloadType)

	downloadQueue, maxRuneCount := downloader.BuildDownloadQueue(m.selectedChapters)

	// Execute download queue, potential to add workerpool here later
	var succeededChapters []structs.Chapter
	for _, job := range downloadQueue {
		err := m.runDownloadJob(job, dl, path, maxRuneCount)

		// Print error to screen, abandon chapter, and continue
		if err != nil {
			logging.Errorln(err.Error)
			fmt.Printf("Chapter %s skipping... reason: %s\n", job.Chapter.Metadata.Num, err.Message)
			continue
		}

		succeededChapters = append(succeededChapters, job.Chapter)
	}

	return succeededChapters
}

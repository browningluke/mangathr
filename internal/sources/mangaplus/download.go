package mangaplus

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
)

func (m *Scraper) runDownloadJob(dl *downloader.Downloader, chapter *manga.Chapter) *logging.ScraperError {
	if err := m.addPagesToChapter(chapter); err != nil {
		return err
	}

	dl.SetTemplate(config.FilenameTemplate)
	chapter.SetFilename(dl.GetNameFromTemplate(chapter, m.MangaTitle(), SCRAPERNAME))

	(*dl.MetadataAgent()).
		SetFromStruct(chapter.Metadata).
		SetPageCount(len(chapter.Pages()))

	if err := dl.CanDownload(chapter); err != nil {
		return err
	}

	if err := dl.Download(chapter); err != nil {
		return &logging.ScraperError{
			Error:   err,
			Message: "An error occurred while downloading a page",
			Code:    0,
		}
	}

	return nil
}

func (m *Scraper) Download(dl *downloader.Downloader, directoryMapping string) []manga.Chapter {
	logging.Debugln("Downloading...")

	directoryName := m.manga.name
	if directoryMapping != "" {
		directoryName = directoryMapping
	}

	dl.SetPath(dl.CreateDirectory(directoryName)).
		SetMaxRuneCount(m.selectedChapters)

	var succeeded []manga.Chapter
	for _, chapter := range m.selectedChapters {
		err := m.runDownloadJob(dl, &chapter)
		if err != nil {
			logging.Errorln(err.Error)
			fmt.Printf("Chapter %s skipping... reason: %s\n", chapter.Metadata.Num, err.Message)
			continue
		}
		succeeded = append(succeeded, chapter)
	}

	return succeeded
}

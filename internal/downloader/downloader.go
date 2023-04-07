package downloader

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/downloader/writer"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/manga"
	"github.com/browningluke/mangathrV2/internal/metadata"
	"github.com/browningluke/mangathrV2/internal/utils"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
	"unicode/utf8"
)

type Downloader struct {
	agent metadata.Agent

	updateMode      bool
	destinationPath string

	enforceChapterDuration bool
	chapterDuration        int64
	maxRuneCount           int
}

type Job struct {
	Chapter manga.Chapter
	Bar     *progressbar.ProgressBar
}

func NewDownloader(updateMode bool,
	enforceChapterDuration bool) *Downloader {
	return &Downloader{
		updateMode:             updateMode,
		enforceChapterDuration: enforceChapterDuration,
	}
}

func (d *Downloader) MetadataAgent() *metadata.Agent {
	d.agent = metadata.NewAgent(config.Metadata.Agent)
	return &d.agent
}

func (d *Downloader) SetChapterDuration(duration int64) {
	d.chapterDuration = duration
}

func (d *Downloader) SetTemplate(template string) {
	if template != "" {
		config.Output.FilenameTemplate = template
	}
}

func (d *Downloader) SetMaxRuneCount(chapters []manga.Chapter) {
	maxRC := 0 // Used for padding (e.g. Chapter 10 vs Chapter 10.5)
	for _, chapter := range chapters {
		// Check if string length is max in list
		if runeCount := utf8.RuneCountInString(chapter.Metadata.Num); runeCount > maxRC {
			maxRC = runeCount
		}
	}
	d.maxRuneCount = maxRC
}

func (d *Downloader) SetPath(path string) {
	d.destinationPath = path
}

/*
	-- Chapter Downloading --
*/

func (d *Downloader) CanDownload(filename string) *logging.ScraperError {
	chapterPath := d.GetChapterPath(filename)

	if config.DryRun {
		return &logging.ScraperError{
			Error:   fmt.Errorf("called with dryRun set to true"),
			Message: "DRY RUN",
			Code:    0,
		}
	}

	if _, err := os.Stat(chapterPath); err == nil {
		return &logging.ScraperError{
			Error:   fmt.Errorf("file exists at path %s", chapterPath),
			Message: "chapter already exists",
			Code:    0,
		}
	}

	return nil
}

// Download chapter. Assumes CanDownload() has been called and has returned true
func (d *Downloader) Download(filename string, pages []manga.Page, metadata *manga.Metadata) error {

	// Initialize progress bar
	bar := utils.CreateProgressBar(len(pages), d.maxRuneCount, metadata.Num)

	// Ensure chapter time is correct
	if d.enforceChapterDuration {
		timeStart := time.Now().UnixMilli()
		defer d.waitChapterDuration(timeStart)
	} else {
		// TODO: differentiate between Download & Update delay
		dur, err := time.ParseDuration(config.Delay.Chapter)
		if err != nil {
			return err
		}
		time.Sleep(dur)
	}

	chapterPath := d.GetChapterPath(filename)

	// Set up writer
	var chapterWriter writer.Writer
	if config.Output.Zip {
		chapterWriter = writer.NewZipWriter(chapterPath)
	} else {
		chapterWriter = writer.NewDirWriter(chapterPath)
	}

	// Write metadata to destination
	filename, body, err := d.agent.GenerateMetadataFile()
	if err != nil {
		return err
	}

	err = chapterWriter.Write(body, filename)
	if err != nil {
		return err
	}

	// Build task array
	var tasks []func()
	for _, page := range pages {
		tasks = append(tasks, buildWorkerPoolFunc(page, bar, func(page *manga.Page) error {
			// Write image bytes to zipfile
			return chapterWriter.Write(page.Bytes, page.Filename())
		}))
	}

	// Run tasks on worker pool
	err = runWorkerPool(tasks, config.SimultaneousPages)
	if err != nil {
		if err := bar.Clear(); err != nil {
			// If the progress bar breaks for some reason, we should panic
			panic(err)
		}
		return err
	}

	fmt.Println("") // Create a new bar for each chapter

	return nil
}

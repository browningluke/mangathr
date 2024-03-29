package downloader

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/downloader/workerpool"
	"github.com/browningluke/mangathr/v2/internal/downloader/writer"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/manga"
	"github.com/browningluke/mangathr/v2/internal/metadata"
	"github.com/browningluke/mangathr/v2/internal/rester"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"os"
	"time"
	"unicode/utf8"
)

type DownloadMode int

const (
	DOWNLOAD DownloadMode = iota
	UPDATE
)

type Downloader struct {
	agent metadata.Agent

	downloadMode    DownloadMode
	destinationPath string

	enforceChapterDuration bool
	chapterDuration        int64
	maxRuneCount           int
}

func NewDownloader(mode DownloadMode, enforceChapterDuration bool) *Downloader {
	return &Downloader{
		downloadMode:           mode,
		enforceChapterDuration: enforceChapterDuration,
	}
}

func (d *Downloader) MetadataAgent() *metadata.Agent {
	d.agent = metadata.NewAgent(config.Metadata.Agent)
	return &d.agent
}

func (d *Downloader) DownloadMode() DownloadMode {
	return d.downloadMode
}

func (d *Downloader) SetChapterDuration(duration int64) *Downloader {
	d.chapterDuration = duration
	return d
}

func (d *Downloader) SetTemplate(template string) *Downloader {
	if template != "" {
		config.Output.FilenameTemplate = template
	}
	return d
}

func (d *Downloader) SetMaxRuneCount(chapters []manga.Chapter) *Downloader {
	maxRC := 0 // Used for padding (e.g. Chapter 10 vs Chapter 10.5)
	for _, chapter := range chapters {
		// Check if string length is max in list
		if runeCount := utf8.RuneCountInString(chapter.Metadata.Num); runeCount > maxRC {
			maxRC = runeCount
		}
	}
	d.maxRuneCount = maxRC
	return d
}

func (d *Downloader) SetPath(path string) *Downloader {
	d.destinationPath = path
	return d
}

/*
	-- Chapter Downloading --
*/

func (d *Downloader) CanDownload(chapter *manga.Chapter) *logging.ScraperError {
	chapterPath := d.GetChapterPath(chapter.Filename())

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

func (d *Downloader) DownloadPage(page *manga.Page) ([]byte, error) {
	logging.Debugln("Starting download of page: ", page.Filename())

	// Parse page time delay
	dur, err := time.ParseDuration(config.Delay.Page)
	if err != nil {
		return nil, err
	}

	time.Sleep(dur)

	imageBytesResp, _ := rester.New().GetBytes(page.Url,
		map[string]string{},
		[]rester.QueryParam{}).Do(config.PageRetries, "100ms")
	pageBytes := imageBytesResp.([]byte)

	err = page.GetExtFromBytes(pageBytes)
	if err != nil {
		return nil, err
	}

	logging.Debugln("Downloaded page. Byte length: ", len(pageBytes))

	return pageBytes, nil
}

// Download chapter. Assumes CanDownload() has been called and has returned true
func (d *Downloader) Download(chapter *manga.Chapter) error {

	// Initialize progress bar
	bar := utils.CreateProgressBar(len(chapter.Pages()), d.maxRuneCount, chapter.Metadata.Num)

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

	chapterPath := d.GetChapterPath(chapter.Filename())

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
	pool := workerpool.New(config.SimultaneousPages)
	for i, _ := range chapter.Pages() {
		p := &chapter.Pages()[i]
		pool.AddTask(func() {
			// Get image bytes to write
			pageBytes, err := d.DownloadPage(p)
			if err != nil {
				panic(err)
			}

			// Write bytes to whichever output
			err = chapterWriter.Write(pageBytes, p.Filename())
			if err != nil {
				panic(err)
			}

			// Add 1 to the bar
			err = bar.Add(1)
			if err != nil {
				panic(err)
			}
		})
	}

	// --- Run pool ---
	var poolErr error

	// Handle pool errors if an error occurred
	defer func(err error) {
		if err != nil {
			// If there were errors, cleanup the writer (only if enabled in config)
			if config.CleanupOnError {
				if err := chapterWriter.Cleanup(); err != nil {
					logging.Errorln("Unable to cleanup file. Reason: ", err)
					fmt.Printf("An error occurred when deleting failed chapter: %s", chapter.Filename())
				}
			}

			// Try to clear progressbar
			if err := bar.Clear(); err != nil {
				logging.Errorln("Unable to clear progress bar. Reason: ", err)
			}
		} else {
			// If there were no errors, mark the writer as complete
			if err := chapterWriter.MarkComplete(); err != nil {
				logging.Errorln("Unable to mark file as complete. Reason: ", err)
			}
		}
	}(poolErr)

	// Run tasks on worker pool (blocking call)
	poolErr = pool.Run()

	// Attempt to close writer
	closeErr := chapterWriter.Close()

	// Propagate close error only if pool ran without errors
	if poolErr == nil && closeErr != nil {
		poolErr = closeErr // Ensures cleanup func is run correctly
	}

	fmt.Println("") // Terminate progress bar (creates a new bar for each chapter)

	return poolErr
}

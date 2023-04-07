package downloader

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/browningluke/mangathrV2/internal/downloader/templater"
	"github.com/browningluke/mangathrV2/internal/manga"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"unicode/utf8"
)

/*
	Path / File system
*/

func CleanPath(path string) string {
	re := regexp.MustCompile(`[<>:"\\|/?*]|\.([<>:"\\|/?*]|$)+`)
	return re.ReplaceAllString(path, "")
}

func (d *Downloader) CreateDirectory(title, downloadType string) string {
	var dirname string

	if downloadType == "download" {
		dirname = config.Output.Path
	} else {
		dirname = config.Output.UpdatePath
	}

	newPath := filepath.Join(dirname, CleanPath(title))

	if !config.DryRun {
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return newPath
}

func (d *Downloader) GetNameFromTemplate(job Job) string {
	return templater.New(&job.Chapter).ExecTemplate(config.Output.FilenameTemplate)
}

func (d *Downloader) GetChapterPath(path, filename string) string {
	// Extract file/dir name (depends on config.output.zip)
	filename = CleanPath(filename)
	if config.Output.Zip {
		filename = fmt.Sprintf("%s.cbz", filename)
	}

	return filepath.Join(path, filename)
}

func (d *Downloader) Cleanup(path, filename string) error {
	chapterPath := d.GetChapterPath(path, filename)

	err := os.RemoveAll(chapterPath)
	if err != nil {
		return err
	}
	return nil
}

/*
	Chapter downloading
*/

func (d *Downloader) waitChapterDuration(timeStart int64) {
	timeEnd := time.Now().UnixMilli()
	downloadDuration := timeEnd - timeStart

	if downloadDuration < d.chapterDuration {
		timeDiff := d.chapterDuration - downloadDuration
		time.Sleep(time.Duration(timeDiff) * time.Millisecond)
	}
}

func BuildDownloadQueue(selectedChapters []manga.Chapter) (jobs []Job, maxRuneCount int) {
	var downloadQueue []Job
	maxRC := 0 // Used for padding (e.g. Chapter 10 vs Chapter 10.5)
	for _, chapter := range selectedChapters {
		downloadQueue = append(downloadQueue, Job{Chapter: chapter})

		// Check if string length is max in list
		if runeCount := utf8.RuneCountInString(chapter.Metadata.Num); runeCount > maxRC {
			maxRC = runeCount
		}
	}
	return downloadQueue, maxRC
}

/*
	Worker pool
*/

func buildWorkerPoolFunc(page Page, bar *progressbar.ProgressBar, writeBytes func(*Page) error) func() {
	return func() {
		// Get image bytes to write
		pageD, err := page.download()
		if err != nil {
			panic(err)
		}

		// Write bytes to whichever output
		err = writeBytes(pageD)
		if err != nil {
			panic(err)
		}

		// Add 1 to the bar
		err = bar.Add(1)
		if err != nil {
			panic(err)
		}
	}
}

func runWorkerPool(tasks []func(), simultaneousPages int) error {
	wpErr := make(chan error)
	panicHandler := func(p interface{}) {
		wpErr <- p.(error)
	}
	pool := pond.New(simultaneousPages, 0, pond.PanicHandler(panicHandler))

	for _, task := range tasks {
		pool.Submit(task)
	}

	for {
		select {
		case err := <-wpErr:
			pool.Stop()
			return err
		default:
			if pool.SubmittedTasks() == pool.CompletedTasks() {
				return nil
			}
		}
	}
}

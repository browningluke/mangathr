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

func (d *Downloader) GetNameFromTemplate(chapter *manga.Chapter) string {
	return templater.New(chapter).ExecTemplate(config.Output.FilenameTemplate)
}

func (d *Downloader) GetChapterPath(filename string) string {
	// Extract file/dir name (depends on config.output.zip)
	filename = CleanPath(filename)
	if config.Output.Zip {
		filename = fmt.Sprintf("%s.cbz", filename)
	}

	return filepath.Join(d.destinationPath, filename)
}

func (d *Downloader) Cleanup(chapter *manga.Chapter) error {
	chapterPath := d.GetChapterPath(chapter.Filename())

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

/*
	Worker pool
*/

func buildWorkerPoolFunc(page manga.Page, bar *progressbar.ProgressBar, writeBytes func(*manga.Page) error) func() {
	return func() {
		// Get image bytes to write
		pageD, err := page.Download(config.Delay.Page, config.PageRetries)
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

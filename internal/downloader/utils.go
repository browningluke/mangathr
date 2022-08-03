package downloader

import (
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
		dirname = d.config.Output.Path
	} else {
		dirname = d.config.Output.UpdatePath
	}

	newPath := filepath.Join(dirname, CleanPath(title))
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	return newPath
}

func (d *Downloader) GetNameFromTemplate(job Job) string {
	templater := &Templater{
		RawTitle: job.Chapter.RawTitle,
		Metadata: job.Chapter.Metadata,
	}
	return templater.ExecTemplate(d.config.Output.FilenameTemplate)
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

func buildWorkerPoolFunc(config *Config, page Page, bar *progressbar.ProgressBar, writeBytes func(*Page) error) func() {
	return func() {
		// Get image bytes to write
		pageD, err := page.download(config)
		handleChapterError(err)

		// Write bytes to whichever output
		err = writeBytes(pageD)
		handleChapterError(err)

		// Add 1 to the bar
		err = bar.Add(1)
		handleChapterError(err)
	}
}

func handleChapterError(err error) {
	if err != nil {
		log.Fatalln(err)
		// todo if page fails more than retries, abandon entire chapter
	}
}

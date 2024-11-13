package downloader

import (
	"github.com/browningluke/mangathr/v2/internal/downloader/templater"
	"github.com/browningluke/mangathr/v2/internal/manga"
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

func (d *Downloader) CreateDirectory(title string) string {
	var dirname string

	if d.downloadMode == DOWNLOAD {
		dirname = config.Output.Path
	} else if d.downloadMode == UPDATE {
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

func (d *Downloader) GetNameFromTemplate(chapter *manga.Chapter, mangaTitle, source string) string {
	return templater.New(chapter, mangaTitle, source).ExecTemplate(config.Output.FilenameTemplate)
}

func (d *Downloader) GetChapterPath(filename string) string {
	// Extract file/dir name (depends on config.output.zip)
	return filepath.Join(d.destinationPath, CleanPath(filename))
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

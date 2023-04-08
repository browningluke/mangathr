package downloader

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/downloader/templater"
	"github.com/browningluke/mangathr/internal/manga"
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
	chapterPath := fmt.Sprintf("%s.part", d.GetChapterPath(chapter.Filename()))

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

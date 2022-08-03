package downloader

import (
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path"
)

func writeToFile(fileBytes []byte, filename string) error {
	logging.Debugln("Writing ", filename, " to directory.")

	// Create empty file
	file, err := os.Create(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	if err != nil {
		return err
	}

	_, err = file.Write(fileBytes)
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) downloadDir(pages []Page, chapterPath string, bar *progressbar.ProgressBar) error {
	// Create directory at chapter path
	err := os.MkdirAll(chapterPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// Write metadata to zip
	filename, body, err := d.agent.GenerateMetadataFile()
	if err != nil {
		return err
	}

	err = writeToFile(body, path.Join(chapterPath, filename))
	if err != nil {
		return err
	}

	// Build task array
	var tasks []func()
	for _, page := range pages {
		tasks = append(tasks, buildWorkerPoolFunc(d.config, page, bar, func(page *Page) error {
			// Write image bytes to img file
			return writeToFile(page.bytes, path.Join(chapterPath, page.Filename))
		}))
	}

	// Run tasks on worker pool
	err = runWorkerPool(tasks, d.config.SimultaneousPages)
	if err != nil {
		if err := bar.Clear(); err != nil {
			// If the progress bar breaks for some reason, we should panic
			panic(err)
		}
		return err
	}

	return nil
}

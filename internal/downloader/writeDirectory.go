package downloader

import (
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/gammazero/workerpool"
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

func (d *Downloader) downloadDir(pages []Page, chapterPath string, bar *progressbar.ProgressBar) {
	// Create directory at chapter path
	err := os.MkdirAll(chapterPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	wp := workerpool.New(d.config.SimultaneousPages)

	for _, page := range pages {
		logging.Debugln("Processing " + page.Filename)

		f := buildWorkerPoolFunc(d.config, page, bar, func(page *Page) error {
			// Write image bytes to img file
			return writeToFile(page.bytes, path.Join(chapterPath, page.Filename))
		})

		wp.Submit(f)
	}
	wp.StopWait()
	filename, body, err := d.agent.GenerateMetadataFile()
	if err != nil {
		panic(err)
	}

	err = writeToFile(body, path.Join(chapterPath, filename))
	if err != nil {
		panic(err)
	}
}

package downloader

import (
	"archive/zip"
	"bytes"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/gammazero/workerpool"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"sync"
)

func writeToZip(fileBytes []byte, filename string, writer *zip.Writer, mu *sync.Mutex) error {
	logging.Debugln("Writing ", filename, " to zip.")

	mu.Lock()
	defer mu.Unlock()
	image, err := writer.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(image, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) downloadZip(pages []Page, chapterPath string, bar *progressbar.ProgressBar) {
	// Create empty file
	archive, err := os.Create(chapterPath)
	defer func(archive *os.File) {
		err := archive.Close()
		if err != nil {
			panic(err)
		}
	}(archive)
	if err != nil {
		panic(err)
	}
	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			panic(err)
		}
	}(zipWriter)

	wp := workerpool.New(d.config.SimultaneousPages)
	var mu sync.Mutex

	for _, page := range pages {
		logging.Debugln("Processing " + page.Filename)

		f := buildWorkerPoolFunc(d.config, page, bar, func(page *Page) error {
			// Write image bytes to zipfile
			return writeToZip(page.bytes, page.Filename, zipWriter, &mu)
		})
		wp.Submit(f)
	}
	wp.StopWait()
	filename, body, err := d.agent.GenerateMetadataFile()
	if err != nil {
		panic(err)
	}

	err = writeToZip(body, filename, zipWriter, &mu)
	if err != nil {
		panic(err)
	}
}

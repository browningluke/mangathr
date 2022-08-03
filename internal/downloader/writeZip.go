package downloader

import (
	"archive/zip"
	"bytes"
	"github.com/browningluke/mangathrV2/internal/logging"
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

func (d *Downloader) downloadZip(pages []Page, chapterPath string, bar *progressbar.ProgressBar) error {
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

	// -- Handle writes ---
	var mu sync.Mutex

	// Write metadata to zip
	filename, body, err := d.agent.GenerateMetadataFile()
	if err != nil {
		return err
	}

	err = writeToZip(body, filename, zipWriter, &mu)
	if err != nil {
		return err
	}

	// Build task array
	var tasks []func()
	for _, page := range pages {
		tasks = append(tasks, buildWorkerPoolFunc(d.config, page, bar, func(page *Page) error {
			// Write image bytes to zipfile
			return writeToZip(page.bytes, page.Filename, zipWriter, &mu)
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

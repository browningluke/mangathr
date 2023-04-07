package writer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"io"
	"os"
	"sync"
)

type zipWriter struct {
	mu     sync.Mutex
	writer *zip.Writer
	file   *os.File
}

func NewZipWriter(chapterPath string) Writer {
	// Create empty partial file
	archive, err := os.Create(fmt.Sprintf("%s", chapterPath))
	if err != nil {
		panic(err)
	}

	zipper := &zipWriter{
		mu:     sync.Mutex{},
		writer: zip.NewWriter(archive),
		file:   archive,
	}

	return zipper
}

func (z *zipWriter) Write(fileBytes []byte, filename string) error {
	logging.Debugln("Writing ", filename, " to zip.")

	z.mu.Lock()
	defer z.mu.Unlock()
	image, err := z.writer.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(image, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	return nil
}

func (z *zipWriter) close() error {
	var err error
	err = z.file.Close()
	err = z.writer.Close()
	return err
}

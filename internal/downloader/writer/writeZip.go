package writer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/browningluke/mangathr/internal/logging"
	"io"
	"os"
	"sync"
)

type zipWriter struct {
	filePath string

	mu     sync.Mutex
	writer *zip.Writer
	file   *os.File
}

func NewZipWriter(chapterPath string) Writer {
	// Create empty partial file
	path := fmt.Sprintf("%s.cbz", chapterPath)
	archive, err := os.Create(getPartPath(path))
	if err != nil {
		panic(err)
	}

	zipper := &zipWriter{
		filePath: path,
		mu:       sync.Mutex{},
		writer:   zip.NewWriter(archive),
		file:     archive,
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

// MarkComplete runs any post-processing if everything else (including close) ran without errors
func (z *zipWriter) MarkComplete() error {
	return os.Rename(getPartPath(z.filePath), z.filePath)
}

func (z *zipWriter) Close() error {
	var err error
	err = z.writer.Close()
	err = z.file.Close()
	return err
}

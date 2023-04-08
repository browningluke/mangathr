package writer

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/logging"
	"os"
	"path/filepath"
)

type dirWriter struct {
	dirPath string
}

func NewDirWriter(chapterPath string) Writer {
	// Create directory at chapter path
	path := fmt.Sprintf("%s", chapterPath)
	err := os.MkdirAll(getPartPath(path), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &dirWriter{
		dirPath: path,
	}
}

func (d *dirWriter) Write(fileBytes []byte, filename string) error {
	logging.Debugln("Writing ", filename, " to directory.")

	// Create empty file
	file, err := os.Create(filepath.Join(getPartPath(d.dirPath), filename))
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

// MarkComplete runs any post-processing if everything else (including close) ran without errors
func (d *dirWriter) MarkComplete() error {
	return os.Rename(getPartPath(d.dirPath), d.dirPath)
}

// Cleanup runs any post-processing any error occurred
func (d *dirWriter) Cleanup() error {
	return os.RemoveAll(getPartPath(d.dirPath))
}

func (d *dirWriter) Close() error {
	return nil
}

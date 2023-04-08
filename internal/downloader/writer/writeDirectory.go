package writer

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/logging"
	"os"
	"path/filepath"
)

type dirWriter struct {
	chapterPath string
}

func NewDirWriter(chapterPath string) Writer {
	// Create directory at chapter path
	path := fmt.Sprintf("%s", chapterPath)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &dirWriter{
		chapterPath: path,
	}
}

func (d *dirWriter) Write(fileBytes []byte, filename string) error {
	logging.Debugln("Writing ", filename, " to directory.")

	// Create empty file
	file, err := os.Create(filepath.Join(d.chapterPath, filename))
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

func (d *dirWriter) Close() error {
	return nil
}

package writer

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"os"
)

type dirWriter struct {
}

func NewDirWriter(chapterPath string) Writer {
	// Create directory at chapter path
	err := os.MkdirAll(fmt.Sprintf("%s", chapterPath), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &dirWriter{}
}

func (d *dirWriter) Write(fileBytes []byte, filename string) error {
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

func (d *dirWriter) Close() error {
	return nil
}

package downloader

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gammazero/workerpool"
	"io"
	"log"
	"mangathrV2/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type Downloader struct {
	config *Config
}

type Page struct {
	Url, Filename string
}

func NewDownloader(config *Config) *Downloader {
	return &Downloader{config: config}
}

func (d *Downloader) CreateDirectory(title, downloadType string) string {
	var dirname string

	if downloadType == "download" {
		if d.config.Output.Path == "" {
			d, err := os.UserHomeDir()
			if err != nil {
				log.Fatalln(err)
			}
			dirname = filepath.Join(d, "mangathrV2")
		} else {
			dirname = d.config.Output.Path
		}
	} else {
		if d.config.Output.UpdatePath == "" {
			d, err := os.UserHomeDir()
			if err != nil {
				log.Fatalln(err)
			}
			dirname = filepath.Join(d, "mangathrV2")
		} else {
			dirname = d.config.Output.UpdatePath
		}
	}

	newPath := filepath.Join(dirname, title)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	return newPath
}

func (d *Downloader) GetNameFromTemplate(pluginTemplate, num, title, language string) string {
	var template string
	if pluginTemplate != "" {
		template = pluginTemplate
	} else {
		template = d.config.Output.FilenameTemplate
	}

	// Do template replacement here
	_ = template

	paddedNum := utils.PadString(num, 3)

	conditionalLanguage := ""
	if language != "" {
		conditionalLanguage = fmt.Sprintf(" - %s", language)
	}

	conditionalTitle := ""
	if title != "" {
		conditionalTitle = fmt.Sprintf(" - %s", title)
	}
	return fmt.Sprintf("%s - Chapter %s%s%s.cbz", paddedNum, num, conditionalLanguage, conditionalTitle)
}

func (d *Downloader) Download(path, chapterFilename string, pages []Page) {
	fmt.Println(chapterFilename)

	chapterPath := filepath.Join(path, chapterFilename)

	if _, err := os.Stat(chapterPath); err == nil {
		fmt.Println("Chapter already exists.")
		return
	} else if errors.Is(err, os.ErrNotExist) {
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

		for _, image := range pages {
			//fmt.Println("Processing " + image.Filename + ".png")
			//fmt.Println("Adding ", image.Filename)
			image := image
			zipWriter := zipWriter
			wp.Submit(func() {
				//mu.Lock()
				//defer mu.Unlock()
				if err := downloadImage(image.Url, image.Filename, zipWriter, &mu); err != nil {
					log.Fatalln(err)
				}
			})

		}
		wp.StopWait()

		// TODO: add this (v) code for ComicInfo metadata agent
		//fmt.Println("Saving ComicInfo")
		//comicInfo, err := zipWriter.Create("ComicInfo.xml")
		//if err != nil {
		//	panic(err)
		//}
		//
		//reader := strings.NewReader("<COMIC_INFO_SYNTAX_HERE>")
		//_, err = io.Copy(comicInfo, reader)
		//if err != nil {
		//	panic(err)
		//}

	} else {
		panic(err)
	}
}

func downloadImage(url, filename string, zipWriter *zip.Writer, mu *sync.Mutex) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("Received code: " + strconv.Itoa(resp.StatusCode))
	}

	fmt.Println("Downloading image: ", filename, " Status Code: ", resp.StatusCode)

	mu.Lock()
	defer mu.Unlock()
	image, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(image, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

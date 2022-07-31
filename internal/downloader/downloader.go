package downloader

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/metadata"
	"github.com/browningluke/mangathrV2/internal/rester"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"github.com/browningluke/mangathrV2/internal/utils"
	"github.com/gammazero/workerpool"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Downloader struct {
	config *Config
	agent  metadata.Agent

	updateMode bool

	enforceChapterDuration bool
	chapterDuration        int64
}

type Page struct {
	Url, Filename string
}

type Job struct {
	Chapter structs.Chapter
	Bar     *progressbar.ProgressBar
}

func NewDownloader(config *Config,
	updateMode bool,
	enforceChapterDuration bool) *Downloader {
	return &Downloader{
		config:                 config,
		updateMode:             updateMode,
		enforceChapterDuration: enforceChapterDuration,
	}
}

func (d *Downloader) MetadataAgent() *metadata.Agent {
	d.agent = metadata.NewAgent(d.config.Metadata.Agent)
	return &d.agent
}

func (d *Downloader) SetChapterDuration(duration int64) {
	d.chapterDuration = duration
}

/*
	-- Utils --
*/

func CleanPath(path string) string {
	re := regexp.MustCompile(`[<>:"\\|/?*]|\.([<>:"\\|/?*]|$)+`)
	return re.ReplaceAllString(path, "")
}

func (d *Downloader) CreateDirectory(title, downloadType string) string {
	var dirname string

	if downloadType == "download" {
		dirname = d.config.Output.Path
	} else {
		dirname = d.config.Output.UpdatePath
	}

	newPath := filepath.Join(dirname, CleanPath(title))
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	return newPath
}

func (d *Downloader) GetNameFromTemplate(pluginTemplate, num, title, language string, groups []string) string {
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
		conditionalLanguage = fmt.Sprintf(" [%s]", language)
	}

	conditionalGroups := ""
	if len(groups) > 0 {
		conditionalGroups = fmt.Sprintf(" [%s]", strings.Join(groups[:], ", "))
	}

	conditionalTitle := ""
	if title != "" {
		conditionalTitle = fmt.Sprintf(" - %s", title)
	}
	return fmt.Sprintf("%s - Chapter %s%s%s%s", paddedNum, num, conditionalTitle,
		conditionalLanguage, conditionalGroups)
}

func (d *Downloader) waitChapterDuration(timeStart int64) {
	timeEnd := time.Now().UnixMilli()
	downloadDuration := timeEnd - timeStart

	if downloadDuration < d.chapterDuration {
		timeDiff := d.chapterDuration - downloadDuration
		time.Sleep(time.Duration(timeDiff) * time.Millisecond)
	}
}

/*
	-- Chapter Downloading --
*/

func (d *Downloader) Download(path, chapterFilename string, pages []Page, bar *progressbar.ProgressBar) {

	// Ensure chapter time is correct
	if d.enforceChapterDuration {
		timeStart := time.Now().UnixMilli()
		defer d.waitChapterDuration(timeStart)
	} else {
		// TODO: differentiate between Download & Update delay
		dur, err := time.ParseDuration(d.config.Delay.Chapter)
		if err != nil {
			panic(err)
		}
		time.Sleep(dur)
	}

	chapterPath := filepath.Join(path, fmt.Sprintf("%s.cbz", CleanPath(chapterFilename)))

	if d.config.DryRun {
		fmt.Println("DRY RUN: not downloading")
		if err := bar.Finish(); err != nil {
			panic(err)
		}
	} else if _, err := os.Stat(chapterPath); err == nil {
		fmt.Println("Chapter already exists.")
		if err := bar.Finish(); err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		d.runDownload(pages, chapterPath, bar)
	} else {
		panic(err)
	}
}

func (d *Downloader) runDownload(pages []Page, chapterPath string, bar *progressbar.ProgressBar) {
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
		logging.Debugln("Processing " + image.Filename)

		// Copy variables
		URL := image.Url
		filename := image.Filename

		wp.Submit(func() {
			if err := d.downloadImage(URL, filename, zipWriter, &mu); err != nil {
				log.Fatalln(err)
			}
			err := bar.Add(1)
			if err != nil {
				panic(err)
			}
		})

	}
	wp.StopWait()
	filename, body := d.agent.GenerateMetadataFile()

	comicInfo, err := zipWriter.Create(filename)
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(body)
	_, err = io.Copy(comicInfo, reader)
	if err != nil {
		panic(err)
	}
}

func (d *Downloader) downloadImage(url, filename string, zipWriter *zip.Writer, mu *sync.Mutex) error {
	logging.Debugln("Starting download of page: ", filename)

	dur, err := time.ParseDuration(d.config.Delay.Page)
	if err != nil {
		return err
	}
	time.Sleep(dur)

	imageBytesResp, _ := rester.New().GetBytes(url,
		map[string]string{},
		[]rester.QueryParam{}).Do(d.config.PageRetries, "100ms")
	imageBytes := imageBytesResp.([]byte)

	logging.Debugln("Downloaded page. Byte length: ", len(imageBytes))

	mu.Lock()
	defer mu.Unlock()
	image, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(image, bytes.NewReader(imageBytes))
	if err != nil {
		return err
	}

	return nil
}

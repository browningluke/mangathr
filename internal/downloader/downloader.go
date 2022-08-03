package downloader

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/metadata"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"time"
)

type Downloader struct {
	config *Config
	agent  metadata.Agent

	updateMode bool

	enforceChapterDuration bool
	chapterDuration        int64
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

func (d *Downloader) SetTemplate(template string) {
	if template != "" {
		d.config.Output.FilenameTemplate = template
	}
}

/*
	-- Chapter Downloading --
*/

func (d *Downloader) Download(path, chapterFilename string, pages []Page, bar *progressbar.ProgressBar) error {

	// Ensure chapter time is correct
	if d.enforceChapterDuration {
		timeStart := time.Now().UnixMilli()
		defer d.waitChapterDuration(timeStart)
	} else {
		// TODO: differentiate between Download & Update delay
		dur, err := time.ParseDuration(d.config.Delay.Chapter)
		if err != nil {
			return err
		}
		time.Sleep(dur)
	}

	// Extract file/dir name (depends on config.output.zip)
	filename := CleanPath(chapterFilename)
	if d.config.Output.Zip {
		filename = fmt.Sprintf("%s.cbz", filename)
	}

	chapterPath := filepath.Join(path, filename)

	if d.config.DryRun {
		fmt.Println("DRY RUN: not downloading")
		if err := bar.Finish(); err != nil {
			// If the progress bar breaks for some reason, we should panic
			panic(err)
		}
		return nil
	} else if _, err := os.Stat(chapterPath); err == nil {
		fmt.Println("Chapter already exists.")
		if err := bar.Finish(); err != nil {
			// If the progress bar breaks for some reason, we should panic
			panic(err)
		}
		return nil
	}

	if d.config.Output.Zip {
		return d.downloadZip(pages, chapterPath, bar)
	} else {
		return d.downloadDir(pages, chapterPath, bar)
	}
}

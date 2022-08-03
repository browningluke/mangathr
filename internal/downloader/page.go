package downloader

import (
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/rester"
	"time"
)

type Page struct {
	Url, Filename string
	bytes         []byte
}

func (p *Page) download(config *Config) (*Page, error) {
	logging.Debugln("Starting download of page: ", p.Filename)

	dur, err := time.ParseDuration(config.Delay.Page)
	if err != nil {
		return p, err
	}
	time.Sleep(dur)

	imageBytesResp, _ := rester.New().GetBytes(p.Url,
		map[string]string{},
		[]rester.QueryParam{}).Do(config.PageRetries, "100ms")
	p.bytes = imageBytesResp.([]byte)

	logging.Debugln("Downloaded page. Byte length: ", len(p.bytes))

	return p, nil
}

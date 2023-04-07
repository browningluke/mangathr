package manga

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/rester"
	"mime"
	"net/http"
	"time"
)

type Page struct {
	Url, Name string
	ext       string // file extension from MIME type
	Bytes     []byte
}

func (p *Page) Download(delayStr string, pageRetries int) (*Page, error) {
	logging.Debugln("Starting download of page: ", p.Filename)

	dur, err := time.ParseDuration(delayStr)
	if err != nil {
		return p, err
	}
	time.Sleep(dur)

	imageBytesResp, _ := rester.New().GetBytes(p.Url,
		map[string]string{},
		[]rester.QueryParam{}).Do(pageRetries, "100ms")
	p.Bytes = imageBytesResp.([]byte)

	// Get mime type
	mimeType := http.DetectContentType(p.Bytes)
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return p, err
	}

	p.ext = ext[len(ext)-1] // Hacky way to get image/jpeg to be .jpg, but keep everything else the same

	logging.Debugln("Downloaded page. Byte length: ", len(p.Bytes))

	return p, nil
}

func (p *Page) Filename() string {
	return fmt.Sprintf("%s%s", p.Name, p.ext)
}

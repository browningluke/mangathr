package manga

import (
	"fmt"
	"mime"
	"net/http"
)

type Page struct {
	Url, Name string
	Ext       string // file extension from MIME type
}

func (p *Page) GetExtFromBytes(pageBytes []byte) error {
	// Get mime type
	mimeType := http.DetectContentType(pageBytes)
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return err
	}

	// Set page extension
	p.Ext = ext[len(ext)-1] // Hacky way to get image/jpeg to be .jpg, but keep everything else the same

	return nil
}

func (p *Page) Filename() string {
	return fmt.Sprintf("%s%s", p.Name, p.Ext)
}

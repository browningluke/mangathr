package ui

import (
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func CreateProgress() *mpb.Progress {
	return mpb.New()
}

func AddBar(p *mpb.Progress, length int64, chapterName string) *mpb.Bar {
	bar := p.AddBar(length,
		mpb.PrependDecorators(
			decor.Name(chapterName, decor.WC{W: len(chapterName) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)

	return bar
}

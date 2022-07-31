package logging

import "github.com/browningluke/mangathrV2/internal/ui"

type ScraperError struct {
	Error   error
	Message string
	Code    int
}

// The following functions are temporary (bad) error handling

// ExitIfError TODO make this better
func ExitIfError(err *ScraperError) {
	if err != nil {
		Errorln(err.Error)
		ui.Fatal(err.Message)
	}
}

func ExitIfErrorWithFunc(err *ScraperError, f func()) {
	if err != nil {
		Errorln(err.Error)
		f()
		ui.Fatal(err.Message)
	}
}

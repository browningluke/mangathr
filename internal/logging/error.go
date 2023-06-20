package logging

import "github.com/browningluke/mangathr/internal/ui"

type ScraperError struct {
	Error   error
	Message string
	Code    int
}

// The following functions are temporary (bad) error handling

func ExitIfError(err *ScraperError) {
	ExitIfErrorWithFunc(err, func() {})
}

func ExitIfErrorWithFunc(err *ScraperError, f func()) {
	if err != nil {
		if err.Error.Error() == "interrupt" {
			Errorln("Caught SIGINT, exiting safely")
			f()
			ui.Fatal("Exiting...")
			return
		}

		Errorln(err.Error)
		f()
		ui.Fatal(err.Message)
	}
}

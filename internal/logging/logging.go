package logging

import (
	"fmt"
	"log"
	"os"
)

var (
	loggingLevel = WARNING

	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
)

func Init() {
	debugLogger = nil
	if loggingLevel <= DEBUG {
		debugLogger = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	infoLogger = nil
	if loggingLevel <= INFO {
		infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	warningLogger = nil
	if loggingLevel <= WARNING {
		warningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	errorLogger = nil
	if loggingLevel <= ERROR {
		errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func SetLoggingLevel(level Level) {
	loggingLevel = level
	Init()
}

func output(logger *log.Logger, s string) {
	if logger != nil {
		err := debugLogger.Output(3, s)
		if err != nil {
			panic(err)
		}
	}
}

// Debug

func Debugln(a ...interface{}) {
	output(debugLogger, fmt.Sprintf("%s%s%s\n", "\033[36m", fmt.Sprint(a...), "\033[0m"))
}

func Debugf(format string, a ...interface{}) {
	output(debugLogger, fmt.Sprintf("%s%s%s\n", "\033[36m", fmt.Sprintf(format, a...), "\033[0m"))
}

// Info

func Infoln(a ...interface{}) {
	output(infoLogger, fmt.Sprint(a...))
}

func Infof(format string, a ...interface{}) {
	output(infoLogger, fmt.Sprintf(format, a...))
}

// Warning

func Warningln(a ...interface{}) {
	output(warningLogger, fmt.Sprintf("%s%s%s\n", "\u001B[33m", fmt.Sprint(a...), "\033[0m"))
}

func Warningf(format string, a ...interface{}) {
	output(warningLogger, fmt.Sprintf("%s%s%s\n", "\033[33m", fmt.Sprintf(format, a...), "\033[0m"))
}

// Error

func Errorln(a ...interface{}) {
	output(errorLogger, fmt.Sprintf("%s%s%s\n", "\033[31m", fmt.Sprint(a...), "\033[0m"))
}

func Errorf(format string, a ...interface{}) {
	output(errorLogger, fmt.Sprintf("%s%s%s\n", "\033[31m", fmt.Sprintf(format, a...), "\033[0m"))
}

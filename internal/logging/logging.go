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

// Debug

func Debugln(log string) {
	if debugLogger != nil {
		debugLogger.Printf("%s%s%s\n", "\033[36m", log, "\033[0m")
	}
}

func Debugf(format string, a ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf("%s%s%s\n", "\033[36m", fmt.Sprintf(format, a...), "\033[0m")
	}
}

// Info

func Infoln(log string) {
	if infoLogger != nil {
		infoLogger.Println(log)
	}
}

func Infof(format string, a ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, a...)
	}
}

// Warning

func Warningln(log string) {
	if warningLogger != nil {
		warningLogger.Printf("%s%s%s\n", "\u001B[33m", log, "\033[0m")
	}
}

func Warningf(format string, a ...interface{}) {
	if warningLogger != nil {
		warningLogger.Printf("%s%s%s\n", "\033[33m", fmt.Sprintf(format, a...), "\033[0m")
	}
}

// Error

func Errorln(log string) {
	if errorLogger != nil {
		errorLogger.Printf("%s%s%s\n", "\033[31m", log, "\033[0m")
	}
}

func Errorf(format string, a ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf("%s%s%s\n", "\033[31m", fmt.Sprintf(format, a...), "\033[0m")
	}
}

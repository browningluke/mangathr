package ui

import (
	"os"
)

/*
	Methods designed for pretty printing to user
*/

// Error Show error message to user
func Error(message ...interface{}) {
	PrintlnColor(Red, message...)
	os.Exit(1)
}

// Errorf Show error message (formatted) to user
func Errorf(format string, message ...interface{}) {
	PrintfColor(Red, format, message...)
	os.Exit(1)
}

// Fatal Show error message to user, then exit
func Fatal(message ...interface{}) {
	Error(message...)
	os.Exit(1)
}

// Fatalf Show error message (formatted) to user, then exit
func Fatalf(format string, message ...interface{}) {
	Errorf(format, message...)
	os.Exit(1)
}

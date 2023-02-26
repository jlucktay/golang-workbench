package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
)

var (
	// Info logs INFO events to '<timestamp>.<domain>.info.log'
	Info *log.Logger

	// Error logs ERROR events to '<timestamp>.<domain>.error.log'
	Error *log.Logger

	logFlags int
)

func init() {
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
}

// Sets logs to write out to their respective handles
func createLogFile(urlScheme, logType string) (io.WriteCloser, *log.Logger) {
	filename := sanitize.Name(fileTimestamp + "." + urlScheme + "-" + flagURL +
		"." + logType + ".log")

	f, errOpen := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if errOpen != nil {
		log.Fatalf("Error opening file: %v", errOpen)
	}

	logger := log.New(f, strings.ToUpper(logType)+": ", logFlags)

	return f, logger
}

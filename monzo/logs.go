package main

import (
	"io"
	"log"
	"os"
)

var (
	// Info logs INFO events to '<timestamp>.<domain>.info.log'
	Info *log.Logger

	// Error logs ERROR events to '<timestamp>.<domain>.error.log'
	Error *log.Logger

	fileTimestamp string
	infoFilename  string
	errorFilename string
	logFlags      int
)

func init() {
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
}

func setupLogs(infoHandle io.WriteCloser, errorHandle io.WriteCloser) {
	// Set info and error logs to write out to their respective handles
	Info = log.New(infoHandle, "INFO: ", logFlags)
	Error = log.New(errorHandle, "ERROR: ", logFlags)
}

func createLogFile(filename string) io.WriteCloser {
	f, errOpen := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if errOpen != nil {
		log.Fatalf("Error opening file: %v", errOpen)
	}

	return f
}

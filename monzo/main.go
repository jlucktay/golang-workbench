package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	// Using two maps, with mutexes to protect them from concurrent updates
	// - one to prevent double-fetching
	// - one for persistent storage of parent/child relationships between pages
	// TODO: time permitting, consolidate this down to one map

	fetched = struct {
		m map[url.URL]error
		sync.Mutex
	}{m: make(map[url.URL]error)}

	errFetchInProgress = errors.New("URL retrieval in progress")

	crawled = struct {
		m map[url.URL][]url.URL
		sync.Mutex
	}{m: make(map[url.URL][]url.URL)}
)

var (
	// Various odds and sods that are utilised all over the package

	flagURL string

	pageCrawled       uint
	pageOutsideDomain uint
)

func init() {
	flag.StringVar(&flagURL, "url", "jameslucktaylor.info", "a URL to crawl")
}

func main() {
	flag.Parse()

	// Parse optional URL from command line
	urlTarget, errParse := url.Parse(fmt.Sprintf("https://%s", flagURL))
	if errParse != nil {
		Error.Printf("Couldn't parse URL '%s': %v\n", flagURL, errParse)
		os.Exit(1)
	}

	// Set up log files with timestamp and URL in their names
	fileTimestamp = time.Now().Format("20060102.150405.000000-0700")

	// Use the same prefix for all log file names, so that they are clustered
	infoFilename = fileTimestamp + "." + flagURL + ".info.log"
	errorFilename = fileTimestamp + "." + flagURL + ".error.log"
	jsonFilename = fileTimestamp + "." + flagURL + ".json"

	// Set info and error logs to write out to their respective files
	infoHandle := createLogFile(infoFilename)
	defer infoHandle.Close()
	errorHandle := createLogFile(errorFilename)
	defer errorHandle.Close()

	setupLogs(infoHandle, errorHandle)

	// Start crawling with recursive function
	crawl(*urlTarget)

	// Print stats to stdout
	fmt.Printf("Pages crawled: %d\nPages outside target '%s' domain: %d\n",
		pageCrawled, urlTarget.String(), pageOutsideDomain)

	// Print findings to JSON file
	outputToJSON()
}

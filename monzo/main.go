package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
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

	timestamp string

	logFlags      int
	logFileName   string
	errorLog      *log.Logger
	errorFileName string

	pageCrawled       uint
	pageOutsideDomain uint
)

func init() {
	flag.StringVar(&flagURL, "url", "jameslucktaylor.info", "a URL to crawl")
}

func main() {
	flag.Parse()

	// Set up log files with timestamp and URL in their names
	// Use the same timestamp etc string for all log file names, so that they are grouped together
	timestamp = time.Now().Format("20060102.150405.000000-0700")
	logFileName = timestamp + "." + flagURL + ".log"
	errorFileName = timestamp + "." + flagURL + ".error.log"
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

	// Set info and error logs to write out to their respective files
	// TODO: refactor log file handling into a func, to make multiple streams easier to wrangle
	f, errOpen := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if errOpen != nil {
		log.Fatalf("Error opening file: %v", errOpen)
	}
	defer f.Close()

	log.SetFlags(logFlags)
	log.SetOutput(f)

	e, errOpen := os.OpenFile(errorFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if errOpen != nil {
		log.Fatalf("Error opening file: %v", errOpen)
	}
	defer e.Close()

	errorLog = log.New(e, "ERROR: ", logFlags)

	// Parse optional URL from command line
	urlTarget, errParse := url.Parse(fmt.Sprintf("https://%s", flagURL))
	if errParse != nil {
		log.Fatal(errParse)
	}

	// Start crawling with recursive function
	crawl(*urlTarget)

	// Print stats to stdout
	fmt.Printf("Pages crawled: %d\nPages outside target '%s' domain: %d\n", pageCrawled, urlTarget.String(), pageOutsideDomain)

	outputToJSON()
}

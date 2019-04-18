package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
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

	fileTimestamp string
)

func init() {
	flag.StringVar(&flagURL, "url", "jameslucktaylor.info", "a URL to crawl")

	// Set up all output files with the same timestamp in their names
	fileTimestamp = time.Now().Format("20060102.150405.000000-0700")
}

func main() {
	flag.Parse()

	// Parse optional URL from command line
	urlTarget, errParse := url.Parse(flagURL)
	if errParse != nil {
		Error.Printf("Couldn't parse URL '%s': %v\n", flagURL, errParse)
		os.Exit(1)
	}

	// Fall back on HTTPS if no protocol was specified as an argument
	if urlTarget.Scheme == "" {
		urlTarget.Scheme = "https"
	} else {
		// Need to keep flagURL as just a domain only, without protocol/scheme
		flagURL = strings.TrimPrefix(flagURL, urlTarget.Scheme+"://")
	}

	// If a bare domain is passed in without a protocol, the parser takes it as
	// a path and not a domain/host name, so we need to swap them
	if urlTarget.Host == "" && urlTarget.Path != "" {
		urlTarget.Host, urlTarget.Path = urlTarget.Path, urlTarget.Host
	}

	// Set info and error logs to write out to their respective files
	infoHandle, Info := createLogFile(urlTarget.Scheme, "info")
	defer infoHandle.Close()
	errorHandle, Error := createLogFile(urlTarget.Scheme, "error")
	defer errorHandle.Close()

	//

	Info.Flags()
	Error.Flags()

	//

	// Start crawling with recursive function
	crawl(*urlTarget)

	// Print stats to stdout
	fmt.Printf("Pages crawled: %d\nPages outside target '%s' domain: %d\n",
		pageCrawled, urlTarget.String(), pageOutsideDomain)

	// Print any findings to JSON file
	if len(crawled.m) > 0 {
		outputToJSON(urlTarget.Scheme)
	}
}

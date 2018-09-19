package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"
)

var (
	// Using two maps, with mutexes to protect them from concurrent updates
	// - one to prevent double-fetching
	// - one for persistent storage of parent/child relationships between pages
	// TODO: time-permitting, either:
	//    consolidate this down to one map, or
	//    send the crawled relationships somewhere else
	//       e.g. stream out to disk (marshal into a JSON dump), a document DB, etc

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

var flagURL string

func init() {
	flag.StringVar(&flagURL, "url", "jameslucktaylor.info", "a URL to crawl")
}

func main() {
	flag.Parse()

	fmt.Println("flagURL has value ", flagURL)

	urlTarget, errParse := url.Parse(fmt.Sprintf("https://%s", flagURL))
	if errParse != nil {
		log.Fatal(errParse)
	}

	genCrawl(*urlTarget, flagURL)
}

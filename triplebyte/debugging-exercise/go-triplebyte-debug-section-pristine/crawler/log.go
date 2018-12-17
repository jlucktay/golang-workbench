package crawler

import (
	"fmt"
	"net/url"
)

// Log is an output manager for crawler
type Log interface {
	enqueue(url.URL)
	finalizeCrawl(url.URL)
	spawn(url.URL)
	headRequest(url.URL)
	getRequest(url.URL)
	noteError(string)
}

type verbose struct{}
type silent struct{}

// Verbose log prints during run
func Verbose() Log {
	return &verbose{}
}

func (*verbose) enqueue(u url.URL) {
	fmt.Printf("url enqueued: %s\n", u.String())
}

func (*verbose) finalizeCrawl(u url.URL) {
	fmt.Printf("finalize crawl: %s\n", u.String())
}

func (*verbose) spawn(u url.URL) {
	fmt.Printf("spawn a crawler to look at: %s\n", u.String())
}

func (*verbose) headRequest(u url.URL) {
	fmt.Printf("crawling with a HEAD request: %s\n", u.String())
}

func (*verbose) getRequest(u url.URL) {
	fmt.Printf("crawling with a GET request: %s\n", u.String())
}

func (*verbose) noteError(s string) {
	fmt.Printf("error! %s\n", s)
}

// Silent log suppresses output during run
func Silent() Log {
	return &silent{}
}

func (*silent) enqueue(_ url.URL) {}

func (*silent) finalizeCrawl(_ url.URL) {}

func (*silent) spawn(_ url.URL) {}

func (*silent) headRequest(_ url.URL) {}

func (*silent) getRequest(_ url.URL) {}

func (*silent) noteError(_ string) {}

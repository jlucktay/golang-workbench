package main

import (
	"net/url"
	"testing"

	"../crawler"
)

func TestChallenge3(t *testing.T) {
	// THIS TEST WILL HANG. Don't sit around waiting for it to finish!

	c := crawler.Crawler{5, crawler.Verbose()}
	graph, err := c.Crawl("http://triplebyte.github.io/web-crawler-test-site/test3/", "")
	if err != nil {
		t.Fatal("Broken test, can't run crawl")
	}

	u, err := url.Parse("http://blah.com:7091")
	if err != nil {
		t.Fatal("Broken test, can't parse URL")
	}

	if _, ok := graph.Nodes[*u]; !ok {
		t.Errorf("Not found: %v", u)
	}
}

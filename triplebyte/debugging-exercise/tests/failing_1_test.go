package main

import (
	"net/url"
	"testing"

	"go.jlucktay.dev/golang-workbench/triplebyte/debugging-exercise/crawler"
)

func TestChallenge1(t *testing.T) {
	c := crawler.Crawler{Threads: 5, Log: crawler.Verbose()}
	graph, err := c.Crawl("triplebyte.github.io/web-crawler-test-site/test1/", "")
	if err != nil {
		t.Fatal("Broken test, can't run crawl")
	}

	u, err := url.Parse("http://triplebyte.github.io/web-crawler-test-site/test1/SVG_logo.svg")
	if err != nil {
		t.Fatal("Broken test, can't parse URL")
	}

	if n := graph.Nodes[*u]; n.Rtype != crawler.HEAD {
		t.Errorf("rtype: %v != %v", crawler.HEAD, n.Rtype)
	}
}

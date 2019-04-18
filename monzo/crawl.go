package main

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func crawl(urlTarget url.URL) {
	if urlTarget.Host != flagURL {
		Info.Printf("'%+v' didn't pass the domain filter '%s', returning.\n",
			urlTarget.String(), flagURL)
		pageOutsideDomain++ // Stats
		return
	}

	fetched.Lock()
	if _, ok := fetched.m[urlTarget]; ok {
		fetched.Unlock()
		Info.Printf("Already fetched '%+v', returning.\n", urlTarget.String())
		return
	}

	fetched.m[urlTarget] = errFetchInProgress
	fetched.Unlock()

	resp, errResp := getResponse(urlTarget)
	if errResp != nil {
		Error.Printf("Response: %+v\n", errResp)
		return
	}

	fetched.Lock()
	fetched.m[urlTarget] = errResp
	fetched.Unlock()

	doc, errRead := goquery.NewDocumentFromReader(resp)
	if errRead != nil {
		Error.Printf("Reading from response: %+v\n", errRead)
		return
	}

	Info.Printf("Fetched '%+v'.\n", urlTarget.String())

	// Keeping the child URLs in a separate slice like this is a bit of a hack
	// I don't like it but it got me past some locking issues
	// TODO: learn more about locking and clean this up
	childResults := getLinks(urlTarget, doc)
	// If we've already crawled this page, we're done here
	if childResults == nil {
		return
	}

	// Now start crawlers on all of this page's children
	done := make(chan bool)
	for b, c := range childResults {
		Info.Printf("Crawling child %+v/%+v of %+v: '%+v'\n",
			b+1, len(childResults), urlTarget.String(), c.String())

		go func(u url.URL) {
			crawl(u)
			done <- true
		}(c)
	}

	for x, y := range childResults {
		Info.Printf("<- [%+v] %+v/%+v - waiting for child: %+v\n",
			urlTarget.String(), x+1, len(childResults), y.String())
		<-done
	}

	Info.Printf("Done with '%+v'.\n", urlTarget.String())
	pageCrawled++ // Stats
}

// Get all links on this page, and store them for later reference
func getLinks(urlTarget url.URL, doc *goquery.Document) []url.URL {
	crawled.Lock()

	if _, ok := crawled.m[urlTarget]; ok {
		crawled.Unlock()
		Info.Printf("Already crawled '%+v', returning.\n", urlTarget)
		return nil
	}

	var children []url.URL

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		urlHref := convertURL(urlTarget.String(), href)

		if urlHref != nil && len(urlHref.String()) > 0 {
			Info.Printf("'%+v' is a child of '%v'.\n",
				urlHref, urlTarget.String())
			children = append(children, *urlHref)
		}
	})

	crawled.m[urlTarget] = children
	crawled.Unlock()

	return children
}

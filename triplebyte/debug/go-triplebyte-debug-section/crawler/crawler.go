package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"../htmlhelp"
)

// Crawler describes the parameters for a crawl
type Crawler struct {
	Threads uint
	Log     Log
}

type crawl struct {
	domain url.URL
	max    uint

	wg *sync.WaitGroup

	mx      *sync.Mutex
	running uint
	log     Log
	queue   []url.URL
	graph   *WebsiteGraph
	errors  []string
}

func newCrawl(c *Crawler, domain url.URL) *crawl {
	return &crawl{
		domain: domain,
		max:    c.Threads,

		wg:    &sync.WaitGroup{},
		mx:    &sync.Mutex{},
		log:   c.Log,
		graph: NewGraph(),
	}
}

func (c *crawl) enqueue(u url.URL, from *url.URL) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.graph.Nodes[u]; ok {
		return
	}

	c.graph.Nodes[u] = PageNode{
		URL: u,
	}
	if from != nil {
		f := *from
		c.graph.Outlinks[f] = append(c.graph.Outlinks[f], u)
		c.graph.Inlinks[f] = append(c.graph.Inlinks[f], u)
	}

	c.log.enqueue(u)

	if c.running < c.max {
		c.running++
		c.wg.Add(1)
		c.spawn(u)
		return
	}

	c.queue = append(c.queue, u)
}

func (c *crawl) finalizeCrawl(u url.URL) {
	c.log.finalizeCrawl(u)

	c.mx.Lock()
	defer c.mx.Unlock()

	qlen := len(c.queue)
	if qlen > 0 {
		u := c.queue[len(c.queue)-1]
		c.queue = c.queue[:len(c.queue)-1]
		c.running++
		c.wg.Add(1)
		c.spawn(u)
	}

	c.running--
	c.wg.Done()
}

func (c *crawl) spawn(u url.URL) {
	c.log.spawn(u)
	go func() {
		if c.shouldBeCrawledAsNode(u) {
			c.getRequest(u)
		} else {
			c.headRequest(u)
		}
	}()
}

func (c *crawl) shouldBeCrawledAsNode(u url.URL) bool {
	if u.Host != c.domain.Host {
		return false
	}

	if u.Scheme != c.domain.Scheme {
		return false
	}

	fileTypes := []string{".pdf", ".jpg", ".gif", ".js", ".css", ".png"}

	for _, t := range fileTypes {
		if strings.HasSuffix(u.Path, t) {
			return false
		}
	}

	return true
}

func (c *crawl) getRequest(u url.URL) {
	c.log.getRequest(u)

	// We race on the individual node, which is ok since we're
	// guaranteed to only request a single URL only once
	node := func() PageNode {
		c.mx.Lock()
		defer c.mx.Unlock()

		return c.graph.Nodes[u]
	}()
	defer func() {
		c.mx.Lock()
		defer c.mx.Unlock()

		c.graph.Nodes[u] = node
	}()

	node.Rtype = GET

	cl := http.Client{Timeout: 10 * time.Second}
	resp, err := cl.Get(u.String())

	if err != nil {
		node.Status = FAILURE
		node.Err = err
		c.noteError("When crawling %s, got an error: %v", u.String(), err)
	} else {
		defer resp.Body.Close()
		node.Status = SUCCESS
		node.Code = resp.StatusCode

		if resp.StatusCode >= 400 {
			c.noteError(
				"When crawling %s, got an error: %d %v",
				u.String(),
				resp.StatusCode,
				http.StatusText(resp.StatusCode),
			)
		} else {
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic("can't read response body?") // TODO
			}

			mime := resp.Header.Get("Content-Type")
			if strings.HasPrefix(mime, "text/") {
				neighbors, errors := htmlhelp.Neighbors(string(contents), u)
				for _, e := range errors {
					c.noteError(e)
				}

				for _, n := range neighbors {
					c.enqueue(n, &u)
				}
			}
		}
	}

	c.finalizeCrawl(u)
}

func (c *crawl) headRequest(u url.URL) {
	c.log.headRequest(u)

	// We race on the individual node, which is ok since we're
	// guaranteed to only request a single URL only once

	node := func() PageNode {
		c.mx.Lock()
		defer c.mx.Unlock()

		return c.graph.Nodes[u]
	}()
	defer func() {
		c.mx.Lock()
		defer c.mx.Unlock()

		c.graph.Nodes[u] = node
	}()

	node.Rtype = HEAD

	cl := http.Client{Timeout: 10 * time.Second}
	resp, err := cl.Head(u.String())

	if err != nil {
		node.Status = FAILURE
		node.Err = err

		c.noteError("When crawling %s, got an error: %v", u.String(), err)
	} else {
		defer resp.Body.Close()
		node.Status = SUCCESS
		node.Code = resp.StatusCode

		if resp.StatusCode >= 400 {
			c.noteError(
				"When crawling %s, got an error: %d %v",
				u.String(),
				resp.StatusCode,
				http.StatusText(resp.StatusCode),
			)
		}

		c.finalizeCrawl(u)
	}
}

func (c *crawl) noteError(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	c.log.noteError(msg)
	c.mx.Lock()
	defer c.mx.Unlock()

	c.errors = append(c.errors, msg)
}

func (c *crawl) wait() (graph *WebsiteGraph, errors []string) {
	c.wg.Wait()
	return c.graph, c.errors
}

// Crawl crawls a web site and produces a report
func (crawler *Crawler) Crawl(target, outfile string) (*WebsiteGraph, error) {
	// TODO outfile?
	if !strings.HasPrefix(target, "http://") {
		target = "http://" + target
	}

	initial, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	domain := *initial
	domain.Path = ""
	domain.RawPath = ""
	domain.RawQuery = ""
	domain.Fragment = ""

	cr := newCrawl(crawler, domain)
	cr.enqueue(*initial, nil)

	_, errors := cr.wait()
	if len(errors) == 0 {
		fmt.Printf("No errors found!\n")
	} else {
		fmt.Println("Here are all the complaints found:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	return cr.graph, nil
}

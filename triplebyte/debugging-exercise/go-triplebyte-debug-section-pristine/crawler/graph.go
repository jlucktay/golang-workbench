package crawler

import (
	"fmt"
	"net/url"
)

// RequestType of the crawl
type RequestType int

// Request types
const (
	HEAD RequestType = iota
	GET
)

func (t RequestType) String() string {
	switch t {
	case HEAD:
		return "HEAD"
	case GET:
		return "GET"
	default:
		panic(fmt.Sprintf("Unknown RequestType %d", t))
	}
}

// NodeStatus is the result of the crawl for this node
type NodeStatus int

// Node Statuses
const (
	NONE NodeStatus = iota
	FAILURE
	SUCCESS
)

func (s NodeStatus) String() string {
	switch s {
	case NONE:
		return "NONE"
	case FAILURE:
		return "FAILURE"
	case SUCCESS:
		return "SUCCESS"
	default:
		panic(fmt.Sprintf("Unknown NodeStatus %d", s))
	}
}

// PageNode is a node in the graph
type PageNode struct {
	URL    url.URL
	Rtype  RequestType
	Status NodeStatus
	Code   int
	Err    error
}

// WebsiteGraph represents the website
type WebsiteGraph struct {
	Nodes    map[url.URL]PageNode
	Outlinks map[url.URL][]url.URL
	Inlinks  map[url.URL][]url.URL
}

// NewGraph initalizes a WebsiteGraph
func NewGraph() *WebsiteGraph {
	return &WebsiteGraph{
		Nodes:    make(map[url.URL]PageNode),
		Outlinks: make(map[url.URL][]url.URL),
		Inlinks:  make(map[url.URL][]url.URL),
	}
}

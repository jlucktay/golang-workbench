package main

import "encoding/xml"

type Rss struct {
	XMLName    xml.Name `xml:"rss"`
	Text       string   `xml:",chardata"`
	Version    string   `xml:"version,attr"`
	ITunes     string   `xml:"itunes,attr"`
	Atom       string   `xml:"atom,attr"`
	GooglePlay string   `xml:"googleplay,attr"`
	Channel    Channel  `xml:"channel"`
}

type Channel struct {
	Text          string `xml:",chardata"`
	Title         string `xml:"title"`
	Link          Link   `xml:"link"`
	Description   string `xml:"description"`
	Owner         Owner  `xml:"owner"`
	Author        string `xml:"author"`
	Image         Image  `xml:"image"`
	Block         string `xml:"block"`
	Language      string `xml:"language"`
	PubDate       string `xml:"pubDate"`
	LastBuildDate string `xml:"lastBuildDate"`
	Item          []Item `xml:"item"`
}

type Link struct {
	Text string `xml:",chardata"`
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Owner struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type Image struct {
	Text  string `xml:",chardata"`
	Href  string `xml:"href,attr"`
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Item struct {
	Text        string    `json:"text" xml:",chardata"`
	Title       string    `json:"title" xml:"title"`
	Link        string    `json:"link" xml:"link"`
	Description string    `json:"description" xml:"description"`
	Enclosure   Enclosure `json:"enclosure" xml:"enclosure"`
	Guid        Guid      `json:"guid" xml:"guid"`
	PubDate     string    `json:"pub_date" xml:"pubDate"`
}

type Enclosure struct {
	Text   string `json:"text" xml:",chardata"`
	URL    string `json:"url" xml:"url,attr"`
	Length string `json:"length" xml:"length,attr"`
	Type   string `json:"type" xml:"type,attr"`
}

type Guid struct {
	Text        string `json:"text" xml:",chardata"`
	IsPermaLink string `json:"is_permalink" xml:"isPermaLink,attr"`
}

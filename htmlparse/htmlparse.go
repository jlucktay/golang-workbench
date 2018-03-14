package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("http://adventofcode.com/2017/day/17")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
}

/*
<article class="day-desc">
*/
func getDayDesc(n *html.Node) string {
	var s string

	if n.Type == html.ElementNode && n.Data == "article" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "day-desc") {
				s += ""
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s += getDayDesc(c)
	}

	return s
}

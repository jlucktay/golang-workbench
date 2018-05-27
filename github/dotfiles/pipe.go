package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func genLinks(input *goquery.Document, filter string) (output chan url.URL) {
	output = make(chan url.URL)

	go func() {
		input.Find("a").Each(func(i int, s *goquery.Selection) {
			src := s.AttrOr("href", "")
			if u, _ := url.Parse(src); u.IsAbs() && u.Scheme != "data" && matchStringCompiled(filter, u.String()) {
				output <- *u
			}
		})

		close(output)
	}()

	return
}

func genAPILang(in chan url.URL) (output chan url.URL) {
	output = make(chan url.URL)

	go func() {
		for i := range in {
			p := strings.Split(i.Path, "/")
			langURL, err := url.Parse(fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", p[1], p[2]))
			if err != nil {
				log.Fatal(err)
			}

			output <- *langURL
		}

		close(output)
	}()

	return
}

func genGoRepos(input chan url.URL) (output chan string) {
	// maybe do this with a map instead? use the repo URL as the key?
	output = make(chan string)

	go func() {
		for i := range input {
			// make API request
			// print language(s)

			output <- fmt.Sprint(getResponse(i))
		}

		close(output)
	}()

	return
}

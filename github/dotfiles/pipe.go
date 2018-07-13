package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func genLinks(input *goquery.Document, filter string) (output chan url.URL) {
	output = make(chan url.URL)

	r, err := regexp.Compile(filter)
	if err != nil {
		panic(err)
	}

	go func() {
		input.Find("a").Each(func(i int, s *goquery.Selection) {
			src := s.AttrOr("href", "")
			if u, _ := url.Parse(src); u.IsAbs() &&
				u.Scheme != "data" &&
				r.MatchString(u.String()) {

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

func genFilterGoRepos(input chan url.URL) (output chan url.URL) {
	output = make(chan url.URL)

	go func() {
		for i := range input {
			resp := fmt.Sprint(getResponse(i))
			var responseHolder map[string]int
			json.Unmarshal([]byte(resp), &responseHolder)

			if filterForGoRepos(responseHolder) {
				output <- i
			}
		}

		close(output)
	}()

	return
}

func genGithubURL(input chan url.URL) chan url.URL {
	output := make(chan url.URL)

	go func() {
		for i := range input {
			// Truncate the '/languages' part from the end of the URL
			apiRepo := i.String()[:len(i.String())-10]

			apiURL, apiErr := url.Parse(apiRepo)
			if apiErr != nil {
				log.Fatal(apiErr)
			}

			resp := fmt.Sprint(getResponse(*apiURL))
			var responseHolder map[string]interface{}
			json.Unmarshal([]byte(resp), &responseHolder)

			if htmlURL, urlOK := responseHolder["html_url"]; urlOK {
				if str, strOK := htmlURL.(string); strOK {
					githubURL, githubErr := url.Parse(str)
					if githubErr != nil {
						log.Fatal(githubErr)
					}

					output <- *githubURL
				}
			}
		}

		close(output)
	}()

	return output
}

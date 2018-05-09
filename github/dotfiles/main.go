package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	dotfilesURL, _ := url.Parse("https://dotfiles.github.io")
	ghRegex := `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`

	res, err := http.DefaultClient.Do(newRequest(*dotfilesURL))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ghLinks := map[string]bool{}

	for link := range genLinks(doc) {
		if matchStringCompiled(ghRegex, link.String()) {
			ghLinks[link.String()] = true
		}
	}

	// fmt.Println(len(ghLinks))
	// fmt.Println(ghLinks)

	// Hit this API endpoint: https://api.github.com/repos/jlucktay/adventofcode/languages
}

func newRequest(u url.URL) *http.Request {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "github/dotfiles")

	return req
}

func genLinks(input *goquery.Document) chan url.URL {
	output := make(chan url.URL)

	go func() {
		input.Find("a").Each(func(i int, s *goquery.Selection) {
			src := s.AttrOr("href", "")
			if u, _ := url.Parse(src); u.IsAbs() && u.Scheme != "data" {
				output <- *u
			}
		})

		close(output)
	}()

	return output
}

func matchStringCompiled(needle, haystack string) bool {
	r, err := regexp.Compile(needle)
	if err != nil {
		panic(err)
	}

	return r.MatchString(haystack)
}

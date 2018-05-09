package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	dotfilesURL, _ := url.Parse("https://dotfiles.github.io")
	ghRegex := `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`

	doc, err := goquery.NewDocumentFromReader(getResponse(*dotfilesURL))
	if err != nil {
		log.Fatal(err)
	}

	for x := range genAPILang(genLinks(doc, ghRegex)) {
		fmt.Println(x)
	}
}

func getResponse(get url.URL) io.Reader {
	res, err := http.DefaultClient.Do(newRequest(get))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf
}

func newRequest(u url.URL) *http.Request {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "github/dotfiles")

	return req
}

func genLinks(input *goquery.Document, filter string) chan url.URL {
	output := make(chan url.URL)

	go func() {
		input.Find("a").Each(func(i int, s *goquery.Selection) {
			src := s.AttrOr("href", "")
			if u, _ := url.Parse(src); u.IsAbs() && u.Scheme != "data" && matchStringCompiled(filter, u.String()) {
				output <- *u
			}
		})

		close(output)
	}()

	return output
}

func genAPILang(in chan url.URL) chan url.URL {
	output := make(chan url.URL)

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

	return output
}

func matchStringCompiled(needle, haystack string) bool {
	r, err := regexp.Compile(needle)
	if err != nil {
		panic(err)
	}

	return r.MatchString(haystack)
}

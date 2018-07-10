package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	dotfilesURL, _ := url.Parse("https://dotfiles.github.io")
	ghRegex := `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`

	doc, err := goquery.NewDocumentFromReader(getResponse(*dotfilesURL))
	if err != nil {
		log.Fatal(err)
	}

	for x := range genGithubURL(genFilterGoRepos(genAPILang(genLinks(doc, ghRegex)))) {
		fmt.Println(x.String())
	}
}

func matchStringCompiled(needle, haystack string) bool {
	r, err := regexp.Compile(needle)
	if err != nil {
		panic(err)
	}

	return r.MatchString(haystack)
}

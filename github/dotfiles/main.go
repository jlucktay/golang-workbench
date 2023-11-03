package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	dotfilesURL, _ := url.Parse("https://dotfiles.github.io/utilities/")
	ghRegex := `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`

	doc, err := goquery.NewDocumentFromReader(getResponse(*dotfilesURL))
	if err != nil {
		log.Fatal(err)
	}

	for x := range genGithubURL(genFilterGoRepos(genAPILang(genLinks(doc, ghRegex)))) {
		fmt.Println(x.String())
	}
}

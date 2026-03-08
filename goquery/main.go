// Thank you:
// https://www.progville.com/go/goquery-jquery-html-golang/
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	resp, err := http.DefaultClient.Get("https://blog.golang.org")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".article").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").Text()
		link, _ := s.Find("h3 a").Attr("href")
		fmt.Printf("%d) %s - %s\n", i+1, title, link)
	})
}

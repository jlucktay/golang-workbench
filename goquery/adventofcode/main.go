// Gets stories from AoC
package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	descriptions := make(map[int]map[string]string)

	getAdventDayDescriptions(7, 25, descriptions)

	for a, b := range descriptions {
		fmt.Printf("a: '%v'\n", a)

		for c, d := range b {
			fmt.Printf("c: '%v'\n", c)
			fmt.Printf("d: '%v'\n", d)
		}
	}
}

func getAdventDayDescriptions(firstDay, lastDay int, m map[int]map[string]string) {
	for index := firstDay; index <= lastDay; index++ {
		url := fmt.Sprintf("http://adventofcode.com/2017/day/%d", index)

		fmt.Println("Fetching '" + url + "'...")

		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("article.day-desc").Each(
			func(i int, s *goquery.Selection) {
				m[index] = make(map[string]string)
				m[index][url] = s.Text()
			},
		)
	}
}

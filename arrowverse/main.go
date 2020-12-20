package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const airdateLayout = "January 2, 2006"

func main() {
	if errPE := printEpisodes("https://arrow.fandom.com/wiki/List_of_Arrow_episodes"); errPE != nil {
		fmt.Fprintf(os.Stderr, "could not print Arrow episode list: %v", errPE)
	}
}

func printEpisodes(episodeListURL string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("arrow.fandom.com"),
		colly.MaxDepth(0),
	)

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table.wikitable", func(i int, table *colly.HTMLElement) {
			table.ForEach("tbody tr", func(_ int, tbody *colly.HTMLElement) {
				if tbody.DOM.ChildrenFiltered("th").Length() > 0 { // Skip <th> row
					return
				}

				episodeNum := strings.TrimSpace(tbody.ChildText("td:nth-of-type(2)"))
				episodeName := strings.Trim(strings.TrimSpace(tbody.ChildText("td:nth-of-type(3)")), `"`)
				episodeLink := tbody.Request.AbsoluteURL(tbody.ChildAttr("td:nth-of-type(3) a", "href"))
				episodeAirdate := strings.TrimSpace(tbody.ChildText("td:nth-of-type(4)"))
				ttAirdate, errParse := time.Parse(airdateLayout, episodeAirdate)
				if errParse != nil {
					return
				}

				fmt.Printf("S%dE%02s %-20s\t%-36s\t%s\n",
					i+1, episodeNum, ttAirdate.Format(airdateLayout), episodeName, episodeLink)
			})
		})
	})

	if errVis := c.Visit(episodeListURL); errVis != nil {
		return errVis
	}

	return nil
}

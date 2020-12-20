package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
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

				const layout = "January 2, 2006"
				airDate := strings.TrimSpace(tbody.ChildText("td:nth-of-type(4)"))
				ttAirDate, errParse := time.Parse(layout, airDate)
				if errParse != nil {
					return
				}

				fmt.Printf("S%dE%02s %-20s\t%-36s\t%s\n", i+1, episodeNum, ttAirDate.Format(layout), episodeName, episodeLink)
			})
		})
	})

	if errVis := c.Visit("https://arrow.fandom.com/wiki/List_of_Arrow_episodes"); errVis != nil {
		return
	}
}

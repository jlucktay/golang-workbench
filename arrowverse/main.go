package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	airdateLayout    = "January 2, 2006"
	allowedDomain    = "arrow.fandom.com"
	categoryListsURL = "https://" + allowedDomain + "/wiki/Category:Lists"
)

func main() {
	episodeListURLs, errPE := getEpisodeListURLs()
	if errPE != nil {
		fmt.Fprintf(os.Stderr, "could not get episode list URLs: %v", errPE)
	}

	for show := range episodeListURLs {
		fmt.Printf("%s\n", show)
	}

	fmt.Println()

	for show, elu := range episodeListURLs {
		if errPE := printEpisodes(show, elu); errPE != nil {
			fmt.Fprintf(os.Stderr, "could not print %s episode list: %v", show, errPE)
		}
	}
}

func getEpisodeListURLs() (map[string]string, error) {
	episodeListURLs := map[string]string{}

	const (
		checkPrefix = "List of "
		checkSuffix = " episodes"
	)

	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(1),
	)

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("div.category-page__members "+
			"ul.category-page__members-for-char "+
			"li.category-page__member "+
			"a.category-page__member-link",
			func(_ int, a *colly.HTMLElement) {
				// Only consider 'List of ... episodes' links
				if !strings.HasPrefix(a.Text, checkPrefix) || !strings.HasSuffix(a.Text, checkSuffix) {
					return
				}

				// Need to specifically exclude the list of crossover episodes
				if strings.Contains(a.Text, " crossover ") {
					return
				}

				showName := strings.TrimPrefix(a.Text, checkPrefix)
				showName = strings.TrimSuffix(showName, checkSuffix)

				episodeListURLs[showName] = a.Request.AbsoluteURL(a.Attr("href"))
			})
	})

	if errVis := c.Visit(categoryListsURL); errVis != nil {
		return nil, fmt.Errorf("error while visiting %s: %w", categoryListsURL, errVis)
	}

	return episodeListURLs, nil
}

func printEpisodes(show, episodeListURL string) error {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
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

				fmt.Printf("%s\tS%dE%02s %-20s\t%-36s\t%s\n",
					show, i+1, episodeNum, ttAirdate.Format(airdateLayout), episodeName, episodeLink)
			})
		})
	})

	if errVis := c.Visit(episodeListURL); errVis != nil {
		return fmt.Errorf("error while visiting %s: %w", episodeListURL, errVis)
	}

	return nil
}

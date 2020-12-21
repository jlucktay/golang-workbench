package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cenkalti/backoff/v4"
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

	for s, elu := range episodeListURLs {
		show, errPE := getEpisodes(s, elu)
		if errPE != nil {
			fmt.Fprintf(os.Stderr, "could not print %s episode list: %v", s, errPE)
		}

		fmt.Printf("%s\n", show)
	}
}

func getEpisodeListURLs() (map[string]string, error) {
	const (
		checkPrefix = "List of "
		checkSuffix = " episodes"
	)

	episodeListURLs := map[string]string{}

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

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff with default settings
	operation := func() error {
		return c.Visit(categoryListsURL)
	}

	if errVis := backoff.Retry(operation, backoff.NewExponentialBackOff()); errVis != nil {
		return nil, fmt.Errorf("error while visiting %s: %w", categoryListsURL, errVis)
	}

	return episodeListURLs, nil
}

func getEpisodes(show, episodeListURL string) (*Show, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(0),
	)

	// Create the new show
	s := &Show{Name: show}

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table.wikitable", func(i int, table *colly.HTMLElement) {
			// Add a new season for this wikitable
			s.Seasons = append(s.Seasons, Season{Number: i + 1})

			table.ForEach("tbody tr", func(_ int, tbody *colly.HTMLElement) {
				if tbody.DOM.ChildrenFiltered("th").Length() > 0 { // Skip <th> row
					return
				}

				epOverall := strings.TrimSpace(tbody.ChildText("td:nth-of-type(1)"))
				epOverallConv, errConvES := strconv.Atoi(epOverall)
				if errConvES != nil {
					return
				}

				epSeason := strings.TrimSpace(tbody.ChildText("td:nth-of-type(2)"))
				epSeasonConv, errConvES := strconv.Atoi(epSeason)
				if errConvES != nil {
					return
				}

				epName := strings.Trim(strings.TrimSpace(tbody.ChildText("td:nth-of-type(3)")), `"`)

				epLink := tbody.Request.AbsoluteURL(tbody.ChildAttr("td:nth-of-type(3) a", "href"))
				link, errUParse := url.Parse(epLink)
				if errUParse != nil {
					return
				}

				epAirdate := strings.TrimSpace(strings.Map(mapSpaces, tbody.ChildText("td:nth-of-type(4)")))
				airdate, errTParse := time.Parse(airdateLayout, epAirdate)
				if errTParse != nil {
					return
				}

				// Add this episode to the current season, indexed by 'i' from body.ForEach
				e := Episode{
					Name:           epName,
					EpisodeSeason:  epSeasonConv,
					EpisodeOverall: epOverallConv,
					Airdate:        airdate,
					Link:           link,
				}

				s.Seasons[i].Episodes = append(s.Seasons[i].Episodes, e)
			})
		})
	})

	// Execute the visit to actually make the HTTP request(s), inside an exponential backoff with default settings
	operation := func() error {
		return c.Visit(episodeListURL)
	}

	if errVis := backoff.Retry(operation, backoff.NewExponentialBackOff()); errVis != nil {
		return nil, fmt.Errorf("error while visiting %s: %w", episodeListURL, errVis)
	}

	return s, nil
}

// Show describes an Arrowverse show.
type Show struct {
	Name    string
	Seasons []Season
}

func (s Show) String() string {
	var b strings.Builder

	for _, season := range s.Seasons {
		fmt.Fprintf(&b, "%s, season %d\n", s.Name, season.Number)
		fmt.Fprintf(&b, "%s\n", season)
	}

	return b.String()
}

// Season describes a season of an Arrowverse show.
type Season struct {
	Number   int
	Episodes []Episode
}

func (s Season) String() string {
	var b strings.Builder

	for _, episode := range s.Episodes {
		fmt.Fprintf(&b, "S%02d%s\n", s.Number, episode)
	}

	return b.String()
}

// Episode describes an episode of an Arrowverse show.
type Episode struct {
	Name           string
	EpisodeSeason  int
	EpisodeOverall int
	Airdate        time.Time
	Link           *url.URL
}

func (e Episode) String() string {
	return fmt.Sprintf("E%02d %-70s\t%-20s\t%s", e.EpisodeSeason, e.Name, e.Airdate.Format(airdateLayout), e.Link)
}

// mapSpaces helps us get rid of non-breaking spaces from HTML.
func mapSpaces(input rune) rune {
	if unicode.IsSpace(input) {
		return ' '
	}
	return input
}

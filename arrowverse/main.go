package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cenkalti/backoff/v4"
	"github.com/gocolly/colly/v2"

	"go.jlucktay.dev/golang-workbench/arrowverse/pkg/models"
)

const (
	allowedDomain    = "arrow.fandom.com"
	categoryListsURL = "https://" + allowedDomain + "/wiki/Category:Lists"
)

func main() {
	episodeListURLs, errPE := GetEpisodeListURLs()
	if errPE != nil {
		fmt.Fprintf(os.Stderr, "could not get episode list URLs: %v", errPE)
	}

	shows := []models.Show{}

	for s, elu := range episodeListURLs {
		show, errPE := GetEpisodes(s, elu)
		if errPE != nil {
			fmt.Fprintf(os.Stderr, "could not get episode details for '%s': %v", s, errPE)
		}

		shows = append(shows, *show)
	}

	for i := range shows {
		fmt.Println(shows[i])
	}
}

// GetEpisodeListURLs will retrieve URLs of all of the 'List of ... episodes' for shows that are available on the wiki.
func GetEpisodeListURLs() (map[string]string, error) {
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

// GetEpisodes will retrieve details for all of the given show's episodes from the wiki.
func GetEpisodes(show, episodeListURL string) (*models.Show, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(0),
	)

	// Create the new show
	s := &models.Show{Name: show}

	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table.wikitable", func(i int, table *colly.HTMLElement) {
			// Add a new season for this wikitable
			s.Seasons = append(s.Seasons, models.Season{Number: i + 1})

			table.ForEach("tbody tr", func(_ int, tbody *colly.HTMLElement) {
				if tbody.DOM.ChildrenFiltered("th").Length() > 0 { // Skip <th> row
					return
				}

				var err error
				ep, itSel := models.Episode{}, NewIteratingSelector()

				// Trim citation link suffixes like "[3]"
				checkCiteSuffix := regexp.MustCompile(`"?\[[0-9]+\]$`)

				if tbody.DOM.ChildrenFiltered("td").Length() >= 4 {
					ep.EpisodeOverall, err = strconv.Atoi(strings.TrimSpace(tbody.ChildText(itSel.Next())))
					if err != nil {
						return
					}
				}

				epSeason := tbody.ChildText(itSel.Next())
				epSeason = checkCiteSuffix.ReplaceAllString(epSeason, "")

				// Handle the 'DC's Legends of Tomorrow' season 5 special episode
				if epSeason == `—` && s.Seasons[i].Number == 5 && s.Name == "DC's Legends of Tomorrow" {
					ep.EpisodeSeason = 0
				} else {
					ep.EpisodeSeason, err = strconv.Atoi(strings.TrimSpace(epSeason))
					if err != nil {
						return
					}
				}

				epName := strings.Trim(strings.TrimSpace(tbody.ChildText(itSel.Next())), `"`)
				ep.Name = checkCiteSuffix.ReplaceAllString(epName, "")

				// Get ahead of too much junk data creeping in
				if ep.Name == "TBA" {
					return
				}

				ep.Link, err = url.Parse(tbody.Request.AbsoluteURL(tbody.ChildAttr(itSel.String()+" a", "href")))
				if err != nil {
					return
				}

				epAirdate := strings.TrimSpace(strings.Map(mapSpaces, tbody.ChildText(itSel.Next())))
				epAirdate = checkCiteSuffix.ReplaceAllString(epAirdate, "")

				// Round off 'TBA' airdates into the future 🤷‍♂️
				if epAirdate == "TBA" {
					theFuture := 5252 - time.Now().Year()
					ep.Airdate = time.Now().AddDate(theFuture, 0, 0).Round(time.Hour * 24)
				} else {
					ep.Airdate, err = time.Parse(models.AirdateLayout, epAirdate)
					if err != nil {
						return
					}
				}

				// Add this episode to the current season, indexed by 'i' from body.ForEach
				s.Seasons[i].Episodes = append(s.Seasons[i].Episodes, ep)
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

// mapSpaces helps us get rid of non-breaking spaces from HTML.
func mapSpaces(input rune) rune {
	if unicode.IsSpace(input) {
		return ' '
	}

	return input
}

// IteratingSelector is a helper to get us through those pesky 'td' selectors.
type IteratingSelector struct {
	selectorFmt string
	tdOffset    int
}

// NewIteratingSelector currently has hard-coded values because we only use it in one loop.
func NewIteratingSelector() *IteratingSelector {
	return &IteratingSelector{
		selectorFmt: "td:nth-of-type(%d)",
		tdOffset:    0,
	}
}

func (is *IteratingSelector) String() string {
	return fmt.Sprintf(is.selectorFmt, is.tdOffset)
}

// Current will return the iterator with its current value.
func (is *IteratingSelector) Current() string {
	return fmt.Sprint(is)
}

// Next will first increment the value, and then return the iterator.
func (is *IteratingSelector) Next() string {
	is.tdOffset++

	return fmt.Sprint(is)
}

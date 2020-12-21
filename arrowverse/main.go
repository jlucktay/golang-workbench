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
	episodeListURLs, errPE := GetEpisodeListURLs()
	if errPE != nil {
		fmt.Fprintf(os.Stderr, "could not get episode list URLs: %v", errPE)
	}

	shows := []Show{}

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
func GetEpisodes(show, episodeListURL string) (*Show, error) {
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

				var err error
				ep, itSel := Episode{}, NewIteratingSelector()

				if tbody.DOM.ChildrenFiltered("td").Length() >= 4 {
					ep.EpisodeOverall, err = strconv.Atoi(strings.TrimSpace(tbody.ChildText(itSel.Next())))
					if err != nil {
						return
					}
				}

				ep.EpisodeSeason, err = strconv.Atoi(strings.TrimSpace(tbody.ChildText(itSel.Next())))
				if err != nil {
					return
				}

				ep.Name = strings.Trim(strings.TrimSpace(tbody.ChildText(itSel.Next())), `"`)

				ep.Link, err = url.Parse(tbody.Request.AbsoluteURL(tbody.ChildAttr(itSel.String()+" a", "href")))
				if err != nil {
					return
				}

				epAirdate := strings.TrimSpace(strings.Map(mapSpaces, tbody.ChildText(itSel.Next())))
				ep.Airdate, err = time.Parse(airdateLayout, epAirdate)
				if err != nil {
					return
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

// Show describes an Arrowverse show.
type Show struct {
	// Name of the show.
	Name string

	// Seasons for this show only.
	Seasons []Season
}

func (s Show) String() string {
	var b strings.Builder

	for _, season := range s.Seasons {
		fmt.Fprintf(&b, "%s, season %d/%d (%d episode(s))\n",
			s.Name, season.Number, len(s.Seasons), len(season.Episodes))
		fmt.Fprintf(&b, "%s\n", season)
	}

	return b.String()
}

// Season describes a season of an Arrowverse show.
type Season struct {
	// Number of the season for the show.
	Number int

	// Episodes within this season only.
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
	// Name of the episode.
	Name string

	// EpisodeSeason is the episode number within the current season.
	EpisodeSeason int

	// EpisodeOverall is the episode number in the overall run of the entire show.
	EpisodeOverall int

	// Airdate is when the episode was first broadcast.
	Airdate time.Time

	// Link to a wiki page with episode details.
	Link *url.URL
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

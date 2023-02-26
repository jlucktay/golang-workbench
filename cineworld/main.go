package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	urlDomain = "https://www.cineworld.co.uk"
	urlPath   = "/uk/data-api-service/v1/quickbook/10108/film-events/in-cinema/%s/at-date/%s"

	cinemaID = "073"
)

//nolint:gochecknoglobals // Flags to pass in arguments with.
var days = flag.Int("days", 0, "number of days into the future")

func main() {
	flag.Parse()

	ctx := context.Background()
	localDate := time.Now().AddDate(0, 0, *days).Local().Format("2006-01-02")
	url := urlDomain + fmt.Sprintf(urlPath, cinemaID, localDate)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatal(err)
	}

	var filmEvents Response
	if err := json.Unmarshal(body, &filmEvents); err != nil {
		log.Fatal("error:", err)
	}

	log.Printf("\n\n%s", filmEvents)
}

// Response and its children structs were all generated thanks to [JSON-to-Go].
//
// [JSON-to-Go]: https://mholt.github.io/json-to-go/
type Response struct {
	Body Body `json:"body"`
}

func (r Response) String() string {
	return fmt.Sprint(r.Body)
}

type Body struct {
	Films  []Film  `json:"films"`
	Events []Event `json:"events"`
}

func (b Body) String() string {
	// Closures that order Film structs.
	name := func(f1, f2 *Film) bool {
		return f1.Name < f2.Name
	}
	increasingLength := func(f1, f2 *Film) bool {
		return f1.Length < f2.Length
	}
	// decreasingLength := func(f1, f2 *Film) bool {
	// 	return f1.Length > f2.Length
	// }

	OrderedBy(name, increasingLength).Sort(b.Films)
	// OrderedBy(name, decreasingLength).Sort(b.Films)
	// OrderedBy(decreasingLength, name).Sort(b.Films)

	var sBuilder strings.Builder

	tabW := new(tabwriter.Writer)

	tabW.Init(&sBuilder, 0, 0, 3, ' ', 0) //nolint:gomnd // Arbitrary padding value.

	for _, film := range b.Films {
		fmt.Fprintf(tabW, "%s\t", film)

		firstEvent := true

		for _, event := range b.Events {
			if event.FilmID == film.ID && event.CinemaID == cinemaID {
				if !firstEvent {
					fmt.Fprint(tabW, " ")
				}

				fmt.Fprintf(tabW, "%s", event)

				firstEvent = false
			}
		}

		fmt.Fprintln(tabW, "\t")
	}

	tabW.Flush()

	return sBuilder.String()
}

type Film struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PosterLink   string   `json:"posterLink"`
	VideoLink    string   `json:"videoLink"`
	Link         string   `json:"link"`
	ReleaseYear  string   `json:"releaseYear"`
	AttributeIDs []string `json:"attributeIds"`
	Length       int      `json:"length"`
	Weight       int      `json:"weight"`
}

func (f Film) String() string {
	runningTime, _ := time.ParseDuration(strconv.Itoa(f.Length) + "m")

	return fmt.Sprintf("%s\t%5s", f.Name,
		strings.TrimSuffix(runningTime.Truncate(time.Minute).String(), "0s"))
}

type lessFunc func(p1, p2 *Film) bool

// MultiSorter implements the Sort interface, sorting the changes within.
type MultiSorter struct {
	films     []Film
	lessFuncs []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *MultiSorter) Sort(changes []Film) {
	ms.films = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *MultiSorter {
	return &MultiSorter{
		films:     nil,
		lessFuncs: less,
	}
}

// Len is part of sort.Interface.
func (ms *MultiSorter) Len() int {
	return len(ms.films)
}

// Swap is part of sort.Interface.
func (ms *MultiSorter) Swap(i, j int) {
	ms.films[i], ms.films[j] = ms.films[j], ms.films[i]
}

// Less is part of sort.Interface. It is implemented by looping along the less functions until it finds a comparison
// that discriminates between the two items (one is less than the other). Note that it can call the less functions
// twice per call. We could change the functions to return -1, 0, 1 and reduce the number of calls for greater
// efficiency: an exercise for the reader.
func (ms *MultiSorter) Less(i, j int) bool {
	left, right := &ms.films[i], &ms.films[j]

	// Try all but the last comparison.
	var index int

	for index = 0; index < len(ms.lessFuncs)-1; index++ {
		less := ms.lessFuncs[index]

		switch {
		case less(left, right):
			// left < right, so we have a decision.
			return true

		case less(right, left):
			// left > right, so we have a decision.
			return false
		}
	} // left == right; try the next comparison.

	// All comparisons to here said "equal", so just return whatever the final comparison reports.
	return ms.lessFuncs[index](left, right)
}

type Event struct {
	CompositeBookingLink CompositeBookingLink `json:"compositeBookingLink"`
	ID                   string               `json:"id"`
	FilmID               string               `json:"filmId"`
	CinemaID             string               `json:"cinemaId"`
	BusinessDay          string               `json:"businessDay"`
	EventDateTime        string               `json:"eventDateTime"`
	BookingLink          string               `json:"bookingLink"`
	PresentationCode     string               `json:"presentationCode"`
	Auditorium           string               `json:"auditorium"`
	AuditoriumTinyName   string               `json:"auditoriumTinyName"`
	AttributeIDs         []string             `json:"attributeIds"`
	SoldOut              bool                 `json:"soldOut"`
}

func (e Event) String() string {
	return e.EventDateTime
}

type CompositeBookingLink struct {
	Type                  string      `json:"type"`
	BookingURL            BookingURL  `json:"bookingUrl"`
	ObsoleteBookingURL    string      `json:"obsoleteBookingUrl"`
	BlockOnlineSales      interface{} `json:"blockOnlineSales"`
	BlockOnlineSalesUntil interface{} `json:"blockOnlineSalesUntil"`
	ServiceURL            string      `json:"serviceUrl"`
}

type BookingURL struct {
	URL    string `json:"url"`
	Params Params `json:"params"`
}

type Params struct {
	SiteCode string `json:"sitecode"`
	Site     string `json:"site"`
	ID       string `json:"id"`
	Lang     string `json:"lang"`
}

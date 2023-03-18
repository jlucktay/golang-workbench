package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/slog"
)

const (
	urlDomain = "https://www.cineworld.co.uk"
	urlPath   = "/uk/data-api-service/v1/quickbook/10108/film-events/in-cinema/%s/at-date/%s"

	cinemaID = "073"
)

//nolint:gochecknoglobals // Flags to pass in arguments with.
var (
	futureDays = flag.Int("f", 0, "start listing from this many days into the future")
)

func main() {
	// Parse any flags passed in.
	flag.Parse()

	// Set up logging.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout)))

	localDate := time.Now().AddDate(0, 0, *futureDays).Local().Format("2006-01-02")
	url := urlDomain + fmt.Sprintf(urlPath, cinemaID, localDate)

	res, err := http.Get(url)
	if err != nil {
		slog.Error("getting URL", err, slog.String("url", url))
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode >= http.StatusMultipleChoices {
		slog.Error("response failed", err, slog.Int("status", res.StatusCode), slog.String("body", string(body)))
		return
	}

	if err != nil {
		slog.Error("reading response body", err)
		return
	}

	var filmEvents Response
	if err := json.Unmarshal(body, &filmEvents); err != nil {
		slog.Error("unmarshaling response body", err)
		return
	}

	fmt.Print(filmEvents)
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

	if len(b.Events) >= 1 {
		xedt := strings.Split(b.Events[0].EventDateTime, "T")
		fmt.Fprintf(tabW, "%s\n", xedt[0])
	}

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
	split := strings.Split(e.EventDateTime, "T")

	if len(split) < 2 {
		return e.EventDateTime
	}

	return split[1]
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

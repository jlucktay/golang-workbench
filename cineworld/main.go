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
	"sync"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/slices"
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
	include3d  = flag.Bool("3", false, "include screenings in 3D")
	listDays   = flag.Int("l", 1, "retrieve listings from this many days")
)

func main() {
	// Parse any flags passed in.
	flag.Parse()

	// Set up logging.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout)))

	// Create somewhere to store results.
	rs := responseStorage{
		responses: make(map[time.Time]Response),
		mx:        sync.Mutex{},
	}

	// Range across number of days we want listings from and request them all concurrently.
	wg := sync.WaitGroup{}

	for i := 0; i < *listDays; i++ {
		wg.Add(1)

		go func(j int) {
			defer wg.Done()

			localDate := time.Now().AddDate(0, 0, *futureDays+j).Local()
			fLocalDate := localDate.Format("2006-01-02")
			url := urlDomain + fmt.Sprintf(urlPath, cinemaID, fLocalDate)

			// Derived logger with URL attached.
			slogw := slog.Default().With(slog.String("url", url))

			res, err := http.Get(url)
			if err != nil {
				slogw.Error("getting URL", err)
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if res.StatusCode >= http.StatusMultipleChoices {
				slogw.Error("response failed", err, slog.Int("status", res.StatusCode), slog.String("body", string(body)))
				return
			}
			if err != nil {
				slogw.Error("reading response body", err)
				return
			}

			var filmEvents Response
			if err := json.Unmarshal(body, &filmEvents); err != nil {
				slogw.Error("unmarshaling response body", err)
				return
			}

			// Store result in map.
			rs.mx.Lock()
			rs.responses[localDate] = filmEvents
			rs.mx.Unlock()
		}(i)
	}

	wg.Wait()

	// Get keys from response storage, to iterate through the map in chronological order and print.
	dateKeys := make([]time.Time, len(rs.responses))

	for key := range rs.responses {
		dateKeys = append(dateKeys, key)
	}

	sort.Slice(dateKeys, func(i, j int) bool {
		return dateKeys[i].Before(dateKeys[j])
	})

	for i := 0; i < len(dateKeys); i++ {
		fmt.Print(rs.responses[dateKeys[i]])
	}
}

type responseStorage struct {
	responses map[time.Time]Response
	mx        sync.Mutex
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
		fmt.Fprintf(tabW, "\n%s\n", xedt[0])
	}

	for _, film := range b.Films {
		fmt.Fprintf(tabW, "%s\t", film)

		firstEvent := true

		for _, event := range b.Events {
			// Check if event's string representation is non-zero length before checking further.
			if event.String() == "" {
				continue
			}

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
	if !*include3d && slices.Contains(e.AttributeIDs, "3d") {
		slog.Debug("event is in 3D",
			slog.Any("attributeIDs", e.AttributeIDs),
			slog.String("dateTime", e.EventDateTime))

		return ""
	}

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

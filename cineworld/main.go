package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/slices"
)

const (
	urlDomain  = "https://www.cineworld.co.uk"
	fmtURLPath = "/uk/data-api-service/v1/quickbook/10108/film-events/in-cinema/%s/at-date/%s"
)

//nolint:gochecknoglobals // Flags to pass in arguments with.
var (
	cinemaID   = flag.String("c", "073", "ID of cinema to pull screenings for")
	futureDays = flag.Int("f", 0, "start listing from this many days into the future")
	include3d  = flag.Bool("3", false, "include screenings in 3D")
	listDays   = flag.Int("l", 1, "retrieve listings from this many days")
)

func main() {
	// Parse any flags passed in.
	flag.Parse()

	// Set up logging.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// If we're not currently in the UK, the Cineworld API gets upset at us.
	if err := checkCountry("GB"); err != nil {
		slog.Error("checking country", slog.Any("err", err))
		return
	}

	// Create somewhere to store results.
	respStore := responseStorage{
		responses: make(map[time.Time]Response),
	}

	// Range across number of days we want listings from and request them all concurrently.
	wgDays := sync.WaitGroup{}

	for daysIntoFuture := 0; daysIntoFuture < *listDays; daysIntoFuture++ {
		wgDays.Add(1)

		go func(j int) {
			defer wgDays.Done()

			localDate := time.Now().AddDate(0, 0, *futureDays+j)
			fLocalDate := localDate.Format("2006-01-02")
			url := urlDomain + fmt.Sprintf(fmtURLPath, *cinemaID, fLocalDate)

			// Derived logger with URL attached.
			slogw := slog.Default().With(slog.String("url", url))

			req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
			if err != nil {
				slogw.Error("creating request",
					slog.Any("err", err))

				return
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				slogw.Error("getting URL",
					slog.Any("err", err))

				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if res.StatusCode >= http.StatusMultipleChoices {
				slogw.Error("response failed",
					slog.Any("err", err),
					slog.Int("status", res.StatusCode),
					slog.String("body", string(body)))

				return
			}

			if err != nil {
				slogw.Error("reading response body",
					slog.Any("err", err))

				return
			}

			var filmEvents Response
			if err := json.Unmarshal(body, &filmEvents); err != nil {
				slogw.Error("unmarshaling response body",
					slog.Any("err", err))

				return
			}

			// Store result in map.
			respStore.Lock()
			respStore.responses[localDate] = filmEvents
			respStore.Unlock()
		}(daysIntoFuture)
	}

	wgDays.Wait()

	// Get keys from response storage, to iterate through the map in chronological order and print.
	dateKeys := make([]time.Time, 0)

	for key := range respStore.responses {
		dateKeys = append(dateKeys, key)
	}

	sort.Slice(dateKeys, func(i, j int) bool {
		return dateKeys[i].Before(dateKeys[j])
	})

	for i := 0; i < len(dateKeys); i++ {
		fmt.Fprint(os.Stdout, respStore.responses[dateKeys[i]])
	}
}

var ErrCountryMismatch = errors.New("expected and actual countries do not match")

func checkCountry(expected string) error {
	httpClient := http.Client{
		Timeout: time.Second * 5, //nolint:gomnd,mnd // Five seconds.
	}

	getIPInfo, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "https://ipinfo.io", nil)
	if err != nil {
		return fmt.Errorf("creating IP info request: %w", err)
	}

	resp, err := httpClient.Do(getIPInfo)
	if err != nil {
		return fmt.Errorf("getting IP info: %w", err)
	}
	defer resp.Body.Close()

	type ipInfoResponse struct {
		Country string `json:"country"`
	}

	iir := &ipInfoResponse{}

	if err := json.NewDecoder(resp.Body).Decode(iir); err != nil {
		return fmt.Errorf("decoding IP info response: %w", err)
	}

	slog.Debug("check country",
		slog.String("expected", expected),
		slog.String("actual", iir.Country))

	if !strings.EqualFold(expected, iir.Country) {
		return fmt.Errorf("%w: %s != %s", ErrCountryMismatch, expected, iir.Country)
	}

	return nil
}

type responseStorage struct {
	sync.Mutex

	responses map[time.Time]Response
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

	OrderedBy(name, increasingLength).Sort(b.Films)

	var sBuilder strings.Builder

	tabW := new(tabwriter.Writer)
	tabW.Init(&sBuilder, 0, 0, 3, ' ', 0) //nolint:gomnd // Arbitrary padding value.

	if len(b.Events) >= 1 {
		xedt := strings.Split(b.Events[0].EventDateTime, "T")
		dateHeader := xedt[0]

		parsedTime, err := time.Parse("2006-01-02", xedt[0])
		if err != nil {
			slog.Error("parsing event date time",
				slog.String("input", xedt[0]),
				slog.Any("err", err))
		} else {
			dateHeader = parsedTime.Format("2006-01-02 Monday")
		}

		fmt.Fprintf(tabW, "\n%s\n", dateHeader)
	}

	for _, film := range b.Films {
		fmt.Fprintf(tabW, "%s\t", film)

		firstEvent := true

		for _, event := range b.Events {
			// Check if event's string representation is non-zero length before checking further.
			if event.String() == "" {
				continue
			}

			if event.FilmID == film.ID && event.CinemaID == *cinemaID {
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

	if len(split) < 2 { //nolint:gomnd // If there is no second element, the return below will panic.
		return e.EventDateTime
	}

	return split[1]
}

type CompositeBookingLink struct {
	Type                  string     `json:"type"`
	BookingURL            BookingURL `json:"bookingUrl"`
	ObsoleteBookingURL    string     `json:"obsoleteBookingUrl"`
	BlockOnlineSales      any        `json:"blockOnlineSales"`
	BlockOnlineSalesUntil any        `json:"blockOnlineSalesUntil"`
	ServiceURL            string     `json:"serviceUrl"`
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

package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	baseURL = "https://www.cineworld.co.uk/api/gatsby-source-boxofficeapi/"

	fmtLocalDate   = "2006-01-02T15:04:05"
	fmtURLMovies   = "movies?basic=false"
	fmtURLSchedule = "schedule?from=%s&to=%s&theaters=%s"
)

//nolint:gochecknoglobals // Flags to pass in arguments with, and a static error.
var (
	include3d  = flag.Bool("3", false, "include screenings in 3D")
	cinemaID   = flag.String("c", "X08C8", "ID of cinema to pull screenings for")
	debugLogs  = flag.Bool("d", false, "emit debug logs")
	futureDays = flag.Int("f", 0, "start listing from this many days into the future")
	listDays   = flag.Int("l", 1, "retrieve listings from this many days")

	ErrCountryMismatch = errors.New("expected and actual countries do not match")
)

func main() {
	// Parse any flags passed in.
	flag.Parse()

	// Set up logging.
	var slogHandlerOpts *slog.HandlerOptions = nil

	if *debugLogs {
		slogHandlerOpts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOpts)))

	// If we're not currently in the UK, the Cineworld API gets upset at us.
	if err := checkCountry("GB"); err != nil {
		slog.Error("checking country", slog.Any("err", err))

		return
	}

	// Set up URL and query parameters to request the schedule.
	timeNow := time.Now().Local()
	nowYear, nowMonth, nowDay := timeNow.Date()

	localDate := time.Date(nowYear, nowMonth, nowDay+*futureDays, 0, 0, 0, 0, timeNow.Location())
	paramFrom := localDate.Format(fmtLocalDate)
	paramTo := localDate.AddDate(0, 0, *listDays).Format(fmtLocalDate)

	timeNowZone, _ := timeNow.Zone()
	paramTheaters := `{"id":"` + *cinemaID + `","timeZone":"` + timeNowZone + `"}`

	scheduleReqURL := baseURL + fmt.Sprintf(fmtURLSchedule, paramFrom, paramTo, url.QueryEscape(paramTheaters))

	// Derived logger with URL attached.
	slogw := slog.Default().With(slog.String("requestURL", scheduleReqURL))

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}

	scheduleReq, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, scheduleReqURL, nil)
	if err != nil {
		slogw.Error("creating schedule request",
			slog.Any("err", err))

		return
	}

	slogw.Debug("created schedule request OK")

	scheduleRes, err := httpClient.Do(scheduleReq)
	if err != nil {
		slogw.Error("getting URL",
			slog.Any("err", err))

		return
	}
	defer scheduleRes.Body.Close()

	slogw.Debug("sent schedule request OK")

	scheduleBody, err := io.ReadAll(scheduleRes.Body)
	if scheduleRes.StatusCode >= http.StatusBadRequest {
		slogw.Error("schedule response failed",
			slog.Any("err", err),
			slog.Int("status", scheduleRes.StatusCode),
			slog.String("body", string(scheduleBody)))

		return
	}

	if err != nil {
		slogw.Error("reading schedule response body",
			slog.Any("err", err))

		return
	}

	slogw.Debug("read schedule response OK", slog.String("body", string(scheduleBody)))

	var gatsbyScheduleResponse GatsbyScheduleResponse
	if err := json.Unmarshal(scheduleBody, &gatsbyScheduleResponse); err != nil {
		slogw.Error("unmarshaling schedule response body",
			slog.Any("err", err))

		return
	}

	movieIDCollection := map[GatsbyMovieID]struct{}{}

	for _, cinema := range gatsbyScheduleResponse {
		for movieID := range cinema.Schedule {
			movieIDCollection[movieID] = struct{}{}
		}
	}

	uniqueMovieIDs := slices.Sorted(maps.Keys(movieIDCollection))

	slogw.Debug("movie ID keys",
		slog.Any("raw", uniqueMovieIDs))

	moviesReqURL := baseURL + fmtURLMovies

	for _, uniqueMovieID := range uniqueMovieIDs {
		moviesReqURL += "&ids=" + string(uniqueMovieID)
	}

	// Replace URL attached to derived logger.
	slogw = slog.Default().With(slog.String("requestURL", moviesReqURL))

	moviesReq, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, moviesReqURL, nil)
	if err != nil {
		slogw.Error("creating movies request",
			slog.Any("err", err))

		return
	}

	slogw.Debug("created movie request OK")

	moviesRes, err := httpClient.Do(moviesReq)
	if err != nil {
		slogw.Error("getting URL",
			slog.Any("err", err))

		return
	}
	defer moviesRes.Body.Close()

	slogw.Debug("sent movie request OK")

	moviesBody, err := io.ReadAll(moviesRes.Body)
	if scheduleRes.StatusCode >= http.StatusBadRequest {
		slogw.Error("movie response failed",
			slog.Any("err", err),
			slog.Int("status", scheduleRes.StatusCode),
			slog.String("body", string(scheduleBody)))

		return
	}

	if err != nil {
		slogw.Error("reading movie response body",
			slog.Any("err", err))

		return
	}

	slog.Debug("read movies response OK", slog.String("body", string(moviesBody)))

	var gatsbyMoviesResponse GatsbyMoviesResponse
	if err := json.Unmarshal(moviesBody, &gatsbyMoviesResponse); err != nil {
		slogw.Error("unmarshaling movies response body",
			slog.Any("err", err))

		return
	}

	slog.Debug("marshalled movies response OK",
		slog.Any("struct", gatsbyMoviesResponse))

	collated, err := collate(gatsbyScheduleResponse, gatsbyMoviesResponse)
	if err != nil {
		slogw.Error("collating responses from Gatsby API",
			slog.Any("err", err))

		return
	}

	fmt.Fprint(os.Stdout, collated)
}

func collate(gsr GatsbyScheduleResponse, gmr GatsbyMoviesResponse) (string, error) {
	type (
		collatedGatsbyResponse map[GatsbyScreeningDate]map[GatsbyMovieID][]GatsbyMovieScreening

		gatsbyMovieDetails struct {
			title, certificate, runtime string
		}
	)

	movieDetails := make(map[GatsbyMovieID]gatsbyMovieDetails)

	for _, movie := range gmr {
		runningTime, err := time.ParseDuration(strconv.Itoa(movie.Runtime) + "s")
		if err != nil {
			return "", fmt.Errorf("parsing running time: %w", err)
		}

		rtRounded := runningTime.Round(time.Minute)

		rtHours := rtRounded / time.Hour
		rtRounded -= rtHours * time.Hour

		rtMinutes := rtRounded / time.Minute

		movieDetails[movie.ID] = gatsbyMovieDetails{
			title:       movie.Title,
			certificate: movie.Certificate,
			runtime:     fmt.Sprintf("%dh%02dm", rtHours, rtMinutes),
		}
	}

	cgr := collatedGatsbyResponse{}

	for _, cinema := range gsr {
		for movieID, screeningOnDate := range cinema.Schedule {
			for screeningDate, movieScreenings := range screeningOnDate {
				if _, alreadyExists := cgr[screeningDate]; !alreadyExists {
					cgr[screeningDate] = make(map[GatsbyMovieID][]GatsbyMovieScreening)
				}

				if _, alreadyExists := cgr[screeningDate][movieID]; !alreadyExists {
					cgr[screeningDate][movieID] = make([]GatsbyMovieScreening, 0)
				}

				cgr[screeningDate][movieID] = append(cgr[screeningDate][movieID], movieScreenings...)
			}
		}
	}

	var sBuilder strings.Builder

	tabW := new(tabwriter.Writer)
	tabW.Init(&sBuilder, 0, 0, 3, ' ', 0)

	for _, screeningDate := range slices.Sorted(maps.Keys(cgr)) {
		parsed, err := time.ParseInLocation("2006-01-02", string(screeningDate), time.Now().Local().Location())
		if err != nil {
			return "", fmt.Errorf("parsing screening date in current location: %w", err)
		}

		movieIDsSortedByTitle := slices.SortedFunc(maps.Keys(cgr[screeningDate]), func(gmi1, gmi2 GatsbyMovieID) int {
			return cmp.Compare(movieDetails[gmi1].title, movieDetails[gmi2].title)
		})

		fmt.Fprintf(tabW, "\n%s\n", parsed.Format("2006-01-02 Monday"))

		for _, movieID := range movieIDsSortedByTitle {
			fmt.Fprintf(tabW, "%s\t%3s %s: ", movieDetails[movieID].title, movieDetails[movieID].certificate, movieDetails[movieID].runtime)

			for _, movieScreening := range cgr[screeningDate][movieID] {
				if !*include3d && slices.Contains(movieScreening.Tags, GatsbyTag3D) {
					slog.Debug("screening is tagged as 3D but the '-3'/'include3d' flag is unset/false",
						slog.String("title", movieDetails[movieID].title),
						slog.Time("startsAt", time.Time(movieScreening.StartsAt)),
						slog.String("tags", strings.Join(movieScreening.Tags, ",")))

					continue
				}

				screenNumber, err := strconv.Atoi(movieScreening.Screen.Name)
				if err != nil {
					return "", fmt.Errorf("parsing screen name into number: %w", err)
				}

				fmt.Fprintf(tabW, "%s (%2d) ", time.Time(movieScreening.StartsAt).Format("15:04:05"), screenNumber)
			}

			fmt.Fprintln(tabW)
		}
	}

	tabW.Flush()

	return sBuilder.String(), nil
}

func checkCountry(expected string) error {
	httpClient := http.Client{
		Timeout: time.Second * 5,
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

// The Gatsby family of structs were all generated thanks to [JSON-to-Go].
//
// [JSON-to-Go]: https://mholt.github.io/json-to-go/

type (
	GatsbyCinemaID      string
	GatsbyMovieID       string
	GatsbyScreeningDate string
	GatsbyMovieStartsAt time.Time

	GatsbyScheduleResponse      map[GatsbyCinemaID]GatsbyCinema
	GatsbyMovieScreeningsOnDate map[GatsbyScreeningDate][]GatsbyMovieScreening
)

type GatsbyCinema struct {
	Schedule map[GatsbyMovieID]GatsbyMovieScreeningsOnDate `json:"schedule,omitzero"`
}

type GatsbyMovieScreening struct {
	StartsAt GatsbyMovieStartsAt `json:"startsAt,omitempty,omitzero"`
	Tags     []string            `json:"tags,omitempty,omitzero"`
	Screen   GatsbyScreen        `json:"screen,omitzero"`
}

const GatsbyTag3D string = "Format.Projection.3d"

type GatsbyScreen struct {
	Name string `json:"name,omitempty,omitzero"`
}

type GatsbyMoviesResponse []struct {
	ID    GatsbyMovieID `json:"id,omitempty,omitzero"`
	Title string        `json:"title,omitempty,omitzero"`

	// Certificate is the rating given by the BBFC.
	Certificate string `json:"certificate,omitempty,omitzero"`

	// Runtime is a number of seconds.
	Runtime int `json:"runtime,omitempty,omitzero"`
}

func (gmsa *GatsbyMovieStartsAt) UnmarshalJSON(data []byte) error {
	parsed, err := time.ParseInLocation("2006-01-02T15:04:05", string(bytes.Trim(data, `"`)), time.Now().Local().Location())
	if err != nil {
		return fmt.Errorf("parsing GatsbyMovieStartsAt in current location: %w", err)
	}

	*gmsa = GatsbyMovieStartsAt(parsed)

	return nil
}

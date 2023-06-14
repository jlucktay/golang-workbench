package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"golang.org/x/exp/slog"
)

func main() {
	const (
		biURL = "https://www.mlb.com/live-stream-games/big-inning"
		tz    = "EST5EDT"
	)

	resp, err := soup.Get(biURL)
	if err != nil {
		slog.Error("could not GET", slog.String("url", biURL), slog.Any("err", err))
		os.Exit(1)
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		slog.Error("could not load location", slog.String("tz", tz), slog.Any("err", err))
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)

	things := doc.Find("table", "class", "p-table").Find("tbody").FindAll("tr")

	for _, thing := range things {
		ft := thing.FullText()
		xft := strings.Split(ft, "\n\t")

		if len(xft) < 3 {
			continue
		}

		parseMe := strings.Join([]string{xft[1], xft[2]}, " ")

		startTime, err := time.ParseInLocation("1/2/2006 3:04 PM", parseMe, location)
		if err != nil {
			slog.Error("could not parse time from string", slog.String("input", parseMe), slog.Any("err", err))
			os.Exit(1)
		}

		if startTime.Local().Hour() >= 8 && startTime.Local().Hour() <= 20 {
			slog.Debug("start time", slog.Time("t", startTime), slog.Time("local", startTime.Local()))

			fmt.Printf("%s\n", startTime.Local().Format(time.RFC850))
		}
	}
}

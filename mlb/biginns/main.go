package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	biURL = "https://www.mlb.com/live-stream-games/big-inning"
	tz    = "EST5EDT"
)

func main() {
	// Initialise config from flags, and bind into Viper.
	pflag.IntP("earliest", "e", 8, "The earliest hour of the start time to show, in 24h time")
	pflag.IntP("latest", "l", 20, "The latest hour of the start time to show, in 24h time")
	pflag.BoolP("debug", "d", false, "Enable debug logging")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing command line flags: %v", err)
		os.Exit(1)
	}

	// Set up logging.
	logOpts := &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.RFC3339,

		// https://github.com/lmittmann/tint#automatically-enabledisable-colors
		NoColor: !isatty.IsTerminal(os.Stderr.Fd()),
	}

	if viper.GetBool("debug") {
		logOpts.AddSource = true
		logOpts.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, logOpts)))

	// Grab the webpage to scrape/parse.
	resp, err := soup.Get(biURL)
	if err != nil {
		slog.Error("could not GET", slog.String("url", biURL), tint.Err(err))
		os.Exit(1)
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		slog.Error("could not load location", slog.String("tz", tz), tint.Err(err))
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)

	table := doc.Find("table", "class", "p-table")
	if table.Error != nil {
		slog.Error("can't find table", tint.Err(table.Error))
		os.Exit(1)
	}

	tableBody := table.Find("tbody")
	if tableBody.Error != nil {
		slog.Error("can't find body inside table", tint.Err(tableBody.Error))
		os.Exit(1)
	}

	tableRows := tableBody.FindAll("tr")
	if len(tableRows) == 0 {
		slog.Error("can't find rows in table body")
		os.Exit(1)
	}

	for _, tableRow := range tableRows {
		if tableRow.Error != nil {
			slog.Error("table row", tint.Err(tableRow.Error))
			continue
		}

		ft := tableRow.FullText()
		xft := strings.Split(ft, "\n\t")

		if len(xft) < 3 {
			slog.Warn("full text had less than 3 fields", slog.Any("xft", xft))
			continue
		}

		xxft1 := strings.Split(xft[1], "/")
		if len(xxft1) < 3 {
			slog.Warn("date had less than 3 fields", slog.String("xft[1]", xft[1]))
			continue
		}

		const (
			parseLayoutPrefix = "1/2/"
			parseLayoutSuffix = " 3:04 PM"
		)

		var parseLayout string

		if len(xxft1[2]) == 2 {
			parseLayout = fmt.Sprintf("%s06%s", parseLayoutPrefix, parseLayoutSuffix)
		} else if len(xxft1[2]) == 4 {
			parseLayout = fmt.Sprintf("%s2006%s", parseLayoutPrefix, parseLayoutSuffix)
		} else {
			slog.Warn("year field in date was not 2 nor 4 characters", slog.String("year-field", xxft1[2]))
			continue
		}

		parseMe := strings.Join([]string{xft[1], xft[2]}, " ")

		startTime, err := time.ParseInLocation(parseLayout, parseMe, location)
		if err != nil {
			slog.Error("could not parse time from string", slog.String("input", parseMe), tint.Err(err))
			os.Exit(1)
		}

		if startTime.Local().Hour() >= viper.GetInt("earliest") && startTime.Local().Hour() <= viper.GetInt("latest") {
			slog.Debug("start time", slog.Time("t", startTime), slog.Time("local", startTime.Local()))

			fmt.Printf("%s\n", startTime.Local().Format(time.RFC850))
		}
	}
}

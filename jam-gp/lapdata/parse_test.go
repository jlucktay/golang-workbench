package lapdata_test

import (
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/golang-workbench/jam-gp/lapdata"
)

func loadEventData(t *testing.T) *lapdata.Event {
	is := is.New(t)
	is.Helper()

	e, err := lapdata.NewEvent()
	is.NoErr(err)
	is.True(e != nil)

	return e
}

func TestNewEvent(t *testing.T) {
	t.Parallel()

	loadEventData(t)
}

func TestCompetitorsHaveLapData(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	e := loadEventData(t)
	is.True(len(e.Session.Competitors) > 0) // no competitor data

	for i := range e.Session.Competitors {
		is.True(len(e.Session.Competitors[i].Laps) > 0) // no lap data for competitor
	}
}

func TestLapDataTotalTime(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	e := loadEventData(t)

	const hourInMilliseconds = 60 * 60 * 1000

	for i := range e.Session.Competitors {
		totalTime := 0

		for j := range e.Session.Competitors[i].Laps {
			totalTime += e.Session.Competitors[i].Laps[j].LapTime
		}

		is.True(totalTime >= hourInMilliseconds*1.75) // each competitor should have at least 1h45m of total lap times

		lastLapIndex := len(e.Session.Competitors[i].Laps) - 1
		is.Equal(totalTime, e.Session.Competitors[i].Laps[lastLapIndex].TotalTime) // stored total time != calculated
	}
}

func TestTykLapDataHasFiveSegments(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	e := loadEventData(t)

	competitorsChecked := 0

	for i := range e.Session.Competitors {
		if e.Session.Competitors[i].Name != "Tyks of Hazzard" || e.Session.Competitors[i].Number != "3" {
			continue
		}

		competitorsChecked++

		is.Equal(e.Session.Competitors[i].Laps.Segments(), 5)
	}

	is.Equal(competitorsChecked, 1)
}

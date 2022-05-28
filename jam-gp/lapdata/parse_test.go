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
			totalTime += e.Session.Competitors[i].Laps[j].Lt
		}

		is.True(totalTime >= hourInMilliseconds*1.75) // each competitor should have at least 1h45m of total lap times
	}
}

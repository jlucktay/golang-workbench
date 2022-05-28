package lapdata_test

import (
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/golang-workbench/jam-gp/lapdata"
)

func TestNewEvent(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	e, err := lapdata.NewEvent()
	is.NoErr(err)
	is.True(e != nil)
}

func TestCompetitorsHaveLapData(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	e, err := lapdata.NewEvent()
	is.NoErr(err)
	is.True(e != nil)

	is.True(len(e.Session.Competitors) > 0) // no competitor data

	for i := range e.Session.Competitors {
		is.True(len(e.Session.Competitors[i].Laps) > 0) // no lap data for competitor
	}
}

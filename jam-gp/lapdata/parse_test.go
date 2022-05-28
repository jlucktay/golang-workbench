package lapdata_test

import (
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/golang-workbench/jam-gp/lapdata"
)

func TestNewEvent(t *testing.T) {
	is := is.New(t)

	_, err := lapdata.NewEvent()
	is.NoErr(err)
}

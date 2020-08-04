package main_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"

	hc "go.jlucktay.dev/golang-workbench/giantbomb.com/hour-counter"
)

func TestRun(t *testing.T) {
	is := is.New(t)
	buf := &bytes.Buffer{}

	is.NoErr(hc.Run([]string{`--api-key="example"`}, buf))
	is.True(strings.Contains(buf.String(), "example"))
}

func TestRunError(t *testing.T) {
	is := is.New(t)
	buf := &bytes.Buffer{}

	errRun := hc.Run([]string{}, buf)

	is.True(errRun != nil)
}

package main_test

import (
	"bytes"
	"testing"

	"github.com/matryer/is"

	hc "go.jlucktay.dev/golang-workbench/giantbomb.com/hour-counter"
)

func TestRun(t *testing.T) {
	is := is.New(t)
	buf := &bytes.Buffer{}

	is.NoErr(hc.Run([]string{"--api-key='example'"}, buf))
}

func TestRunError(t *testing.T) {
	is := is.New(t)
	buf := &bytes.Buffer{}

	errRun := hc.Run([]string{}, buf)

	is.True(errRun != nil)
}

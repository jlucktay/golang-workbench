package main_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"

	hc "go.jlucktay.dev/golang-workbench/giantbomb.com/hour-counter"
)

func TestRun(t *testing.T) {
	t.Skip("these tests need some stubs/mocks to have a shot at working correctly")

	is := is.New(t)
	buf := &bytes.Buffer{}

	is.NoErr(hc.Run([]string{`--api-key="example"`}, buf))
	is.True(strings.Contains(buf.String(), "example"))
}

func TestRunError(t *testing.T) {
	t.Skip("these tests need some stubs/mocks to have a shot at working correctly")

	is := is.New(t)
	buf := &bytes.Buffer{}

	errRun := hc.Run([]string{}, buf)
	t.Logf("errRun: '%#v'", errRun)
	is.True(errRun != nil)
}

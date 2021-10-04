package main_test

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/matryer/is"

	main "go.jlucktay.dev/golang-workbench/stdlib/io/fs"
)

var testFS = fstest.MapFS{
	"home_namespace.yaml": &fstest.MapFile{
		Data:    []byte(`hello`),
		Mode:    0o440,
		ModTime: time.Now(),
	},
}

func TestStatHomeNS(t *testing.T) {
	is := is.New(t)

	is.NoErr(main.StatHomeNS(testFS))
}

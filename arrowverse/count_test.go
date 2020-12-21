package main_test

import (
	"testing"

	"github.com/matryer/is"

	arrowverse "go.jlucktay.dev/golang-workbench/arrowverse"
)

func TestGetEpisodeListURLs(t *testing.T) {
	is := is.New(t)

	episodeListURLs, errPE := arrowverse.GetEpisodeListURLs()
	is.NoErr(errPE)
	is.True(len(episodeListURLs) >= 11)
}

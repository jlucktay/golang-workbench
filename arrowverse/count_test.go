package main_test

import (
	"testing"

	"github.com/matryer/is"

	main "go.jlucktay.dev/golang-workbench/arrowverse"
)

func TestGetEpisodeListURLs(t *testing.T) {
	is := is.New(t)

	episodeListURLs, errPE := main.GetEpisodeListURLs()
	is.NoErr(errPE)
	is.True(len(episodeListURLs) >= 11)
}

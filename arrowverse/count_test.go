package main_test

import (
	"testing"
	"time"

	"github.com/matryer/is"

	arrowverse "go.jlucktay.dev/golang-workbench/arrowverse"
)

func TestEpisodeNumbers(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	// Store show/season/episode details for seasons that have finished airing completely.
	showSeasonEpisodes := map[string]map[int]int{
		"Arrow":                     {1: 23, 2: 23, 3: 23, 4: 23, 5: 23, 6: 23, 7: 22, 8: 10},
		"Batwoman":                  {1: 20},
		"Birds of Prey":             {1: 13},
		"Black Lightning":           {1: 13, 2: 16, 3: 16},
		"Constantine":               {1: 13},
		"DC's Legends of Tomorrow":  {1: 16, 2: 17, 3: 18, 4: 16, 5: 15},
		"Freedom Fighters: The Ray": {1: 6, 2: 6},
		"Supergirl":                 {1: 20, 2: 22, 3: 23, 4: 22, 5: 19},
		"The Flash (CBS)":           {1: 22},
		"The Flash (The CW)":        {1: 23, 2: 23, 3: 23, 4: 23, 5: 22, 6: 19},
		"Vixen":                     {1: 6, 2: 6},
	}

	episodeListURLs, errGELU := arrowverse.GetEpisodeListURLs()
	is.NoErr(errGELU)
	is.Equal(len(episodeListURLs), len(showSeasonEpisodes))

	for s, elu := range episodeListURLs {
		// Pin! ref: https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		s, elu := s, elu

		t.Run(s, func(t *testing.T) {
			// Don't use .Parallel() without pinning
			t.Parallel()

			show, errGE := arrowverse.GetEpisodes(s, elu)
			is.NoErr(errGE)

			seasonNumbers, haveShowNumbers := showSeasonEpisodes[s]
			is.True(haveShowNumbers)

			for i := 0; i < len(show.Seasons); i++ {
				// If we have numbers for this season, make sure they match, then move on to the next season
				if episodeCount, haveSeasonNumbers := seasonNumbers[i+1]; haveSeasonNumbers {
					is.Equal(len(show.Seasons[i].Episodes), episodeCount)

					continue
				}

				// Otherwise, make sure we have at least one episode to inspect before going ahead
				lastEpIdx := len(show.Seasons[i].Episodes) - 1

				if lastEpIdx < 0 {
					continue
				}

				// Check the retrieved airdate against today's date
				if show.Seasons[i].Episodes[lastEpIdx].Airdate.Before(time.Now()) {
					t.Fatalf("missing episode count for S%02d of '%s' which has finished airing", i, show.Name)
				}
			}
		})
	}
}

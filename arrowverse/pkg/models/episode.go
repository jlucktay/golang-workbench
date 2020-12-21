package models

import (
	"fmt"
	"net/url"
	"time"
)

// Episode describes an episode of an Arrowverse show.
type Episode struct {
	// Name of the episode.
	Name string

	// EpisodeSeason is the episode number within the current season.
	EpisodeSeason int

	// EpisodeOverall is the episode number in the overall run of the entire show.
	EpisodeOverall int

	// Airdate is when the episode was first broadcast.
	Airdate time.Time

	// Link to a wiki page with episode details.
	Link *url.URL
}

func (e Episode) String() string {
	return fmt.Sprintf("E%02d %-70s\t%-20s\t%s", e.EpisodeSeason, e.Name, e.Airdate.Format(AirdateLayout), e.Link)
}

package models

import (
	"fmt"
	"strings"
)

// Season describes a season of an Arrowverse show.
type Season struct {
	// Number of the season for the show.
	Number int

	// Episodes within this season only.
	Episodes []Episode
}

func (s Season) String() string {
	var b strings.Builder

	for _, episode := range s.Episodes {
		fmt.Fprintf(&b, "S%02d%s\n", s.Number, episode)
	}

	return b.String()
}

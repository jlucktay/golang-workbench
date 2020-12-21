package models

import (
	"fmt"
	"strings"
)

// Show describes an Arrowverse show.
type Show struct {
	// Name of the show.
	Name string

	// Seasons for this show only.
	Seasons []Season
}

func (s Show) String() string {
	var b strings.Builder

	for _, season := range s.Seasons {
		fmt.Fprintf(&b, "%s, season %d/%d (%d episode(s))\n",
			s.Name, season.Number, len(s.Seasons), len(season.Episodes))
		fmt.Fprintf(&b, "%s\n", season)
	}

	return b.String()
}

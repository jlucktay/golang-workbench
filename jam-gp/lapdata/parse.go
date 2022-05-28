package lapdata

import (
	"encoding/json"
	"fmt"
)

// NewEvent will parse the raw data and return a struct containing such.
func NewEvent() (*Event, error) {
	e := &Event{}

	err := json.Unmarshal(jamRaw, e)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal raw data: %w", err)
	}

	return e, nil
}

// SliceOfLaps does what it says on the tin.
type SliceOfLaps []Laps

// Segments will compute the number of turns taken by different drivers.
func (laps SliceOfLaps) Segments() int {
	// Begin with the time from lap one, as lap zero is technically the time taken from being stationary on the grid to
	// crossing the start/finish line for the first time.
	runningAverage := laps[1].LapTime
	segmentCount := 1

	// Longer than usual lap times denote laps where we changed drivers.
	for i := range laps {
		if i <= 1 {
			continue
		}

		// If the current lap time is more than 120% of the current running average, it (probably) signifies a driver
		// change.
		if laps[i].LapTime > (runningAverage*6)/5 {
			segmentCount++

			// Use the previous lap as a new baseline.
			runningAverage = laps[i-1].LapTime
		} else {
			runningAverage = (runningAverage + laps[i].LapTime) / 2
		}
	}

	return segmentCount
}

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

package main

import (
	"encoding/json"
	"strings"
	"time"
)

func (i *Item) ToJSON() (string, error) {
	sb := strings.Builder{}

	enc := json.NewEncoder(&sb)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(i); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}

func (i *Item) CreatedAt() (float64, error) {
	pd, err := time.Parse(time.RFC1123, i.PubDate)
	if err != nil {
		return 0, err
	}

	return float64(pd.Unix()), nil
}

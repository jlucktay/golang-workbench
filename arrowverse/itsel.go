package main

import "fmt"

// IteratingSelector is a helper to get us through those pesky 'td' selectors.
type IteratingSelector struct {
	selectorFmt string
	tdOffset    int
}

// NewIteratingSelector currently has hard-coded values because we only use it in one loop.
func NewIteratingSelector() *IteratingSelector {
	return &IteratingSelector{
		selectorFmt: "td:nth-of-type(%d)",
		tdOffset:    0,
	}
}

func (is *IteratingSelector) String() string {
	return fmt.Sprintf(is.selectorFmt, is.tdOffset)
}

// Current will return the iterator with its current value.
func (is *IteratingSelector) Current() string {
	return is.String()
}

// Next will first increment the value, and then return the iterator.
func (is *IteratingSelector) Next() string {
	is.tdOffset++

	return is.String()
}

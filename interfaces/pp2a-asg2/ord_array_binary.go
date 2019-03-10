package collection

import (
	"fmt"
	"io"
	"sort"
)

// OrdArrayBinary is an ordered slice with linear insert, and binary search
type OrdArrayBinary struct {
	words []string
	size  int
}

// MakeCollection initialises the WordCollection, up to an initial capacity
// defined by WCSIZE. It returns SUCCESS upon successful completion.
func (o *OrdArrayBinary) MakeCollection() int {
	o.words = make([]string, 0, WCSIZE)
	o.size = 0

	return SUCCESS
}

// FreeCollection frees the memory dynamically allocated to the WordCollection.
func (o *OrdArrayBinary) FreeCollection() {
	o.words = nil
	o.size = 0
}

// AddCollection adds the string given by the parameter to the WordCollection.
// The string is added so that the WordCollection is in alphabetical order.
func (o *OrdArrayBinary) AddCollection(word string) int {
	i := o.size

	// Find the right index to insert at
	for i > 0 && o.words[i-1] > word {
		i--
	}

	// Shuffle the slice items in front of our desired index forward
	o.words = append(o.words, "")
	copy(o.words[i+1:], o.words[i:])

	// Insert new word at desired index, to maintain alphabetical order
	o.words[i] = word

	// Increment word count
	o.size++

	return SUCCESS
}

// SearchCollection searches for the string given by the parameter in the
// WordCollection. This utilises a binary search algorithm, and returns SUCCESS
// or FAILURE depending upon the outcome of the search.
func (o *OrdArrayBinary) SearchCollection(needle string) int {
	i := sort.Search(
		len(o.words),
		func(i int) bool {
			return o.words[i] >= needle
		},
	)
	if i < len(o.words) && o.words[i] == needle {
		return SUCCESS
	}

	return FAILURE
}

// SizeCollection returns the number of words in the WordCollection.
func (o *OrdArrayBinary) SizeCollection() int {
	return o.size
}

// DisplayCollection prints the contents of the WordCollection to the Writer
// given by the parameter.
func (o *OrdArrayBinary) DisplayCollection(w io.Writer) {
	for i := 0; i < o.size; i++ {
		fmt.Fprintf(w, "Element %d:\t%s\n", i, o.words[i])
	}
}

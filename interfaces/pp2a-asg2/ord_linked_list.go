package collection

import "io"

// OrdLinkedList is an ordered linked list with linear insert and search
type OrdLinkedList struct {
	head *ordLinkedListNode
	size int
}

// ordLinkedListNode is a node used in OrdLinkedList
type ordLinkedListNode struct {
	word string
	next *ordLinkedListNode
}

// MakeCollection initialises the WordCollection with one initial nil entry. It
// returns SUCCESS upon successful completion.
func (o *OrdLinkedList) MakeCollection() int {
	o.head = nil
	o.size = 0

	return SUCCESS
}

// FreeCollection frees the memory dynamically allocated to the WordCollection.
func (o *OrdLinkedList) FreeCollection() {
	o.head = nil
	o.size = 0
}

// AddCollection adds the string given by the parameter to the WordCollection.
// The string is added so that the WordCollection is in alphabetical order.
func (o *OrdLinkedList) AddCollection(word string) int {
	return FAILURE
}

// SearchCollection searches for the string given by the parameter in the
// WordCollection. This utilises a linear search algorithm, and returns SUCCESS
// or FAILURE depending upon the outcome of the search.
func (o *OrdLinkedList) SearchCollection(word string) int {
	return FAILURE
}

// SizeCollection returns the number of words in the WordCollection.
func (o *OrdLinkedList) SizeCollection() int {
	return FAILURE
}

// DisplayCollection prints the contents of the WordCollection to the Writer
// given by the parameter.
func (o *OrdLinkedList) DisplayCollection(w io.Writer) {}

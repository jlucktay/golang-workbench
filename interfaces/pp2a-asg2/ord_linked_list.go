package collection

import (
	"fmt"
	"io"
)

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
	newNode := &ordLinkedListNode{}
	current := o.head
	var previous *ordLinkedListNode

	for current != nil && current.word < word {
		previous = current
		current = current.next
	}

	newNode.word = word
	newNode.next = current
	o.size++

	if previous == nil {
		o.head = newNode
	} else {
		previous.next = newNode
	}

	return SUCCESS
}

// SearchCollection searches for the string given by the parameter in the
// WordCollection. This utilises a linear search algorithm, and returns SUCCESS
// or FAILURE depending upon the outcome of the search.
func (o *OrdLinkedList) SearchCollection(needle string) int {
	current := o.head

	for current != nil {
		if current.word == needle {
			return SUCCESS
		} else if current.word > needle {
			return FAILURE
		} else {
			current = current.next
		}
	}

	return FAILURE
}

// SizeCollection returns the number of words in the WordCollection.
func (o *OrdLinkedList) SizeCollection() int {
	return o.size
}

// DisplayCollection prints the contents of the WordCollection to the Writer
// given by the parameter.
func (o *OrdLinkedList) DisplayCollection(w io.Writer) {
	i := 0
	current := o.head

	for current != nil {
		fmt.Fprintf(w, "Element %d:\t%s\n", i, current.word)
		i++
		current = current.next
	}
}

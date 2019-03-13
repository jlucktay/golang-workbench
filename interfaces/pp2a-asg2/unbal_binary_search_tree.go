package collection

import (
	"fmt"
	"io"
)

// UnbalBinarySearchTree is an unbalanced binary search tree.
type UnbalBinarySearchTree struct {
	root *ubstNode
	size int
}

// ubstNode is a node used in UnbalBinarySearchTree.
type ubstNode struct {
	word        string
	left, right *ubstNode
}

// MakeCollection initialises the WordCollection with an initial nil entry. It
// returns SUCCESS upon successful completion.
func (u *UnbalBinarySearchTree) MakeCollection() int {
	u.root = nil
	u.size = 0

	return SUCCESS
}

// FreeCollection frees the memory dynamically allocated to the WordCollection.
// It achieves this by setting the root node pointer to nil, so that the
// garbage collector can clean up all of the now-unlinked nodes.
func (u *UnbalBinarySearchTree) FreeCollection() {
	u.root = nil
	u.size = 0
}

// AddCollection adds the string to the WordCollection.
func (u *UnbalBinarySearchTree) AddCollection(word string) int {
	current, newNode := u.root, &ubstNode{}
	var previous *ubstNode

	for current != nil {
		previous = current

		if current.word > word {
			current = current.left
		} else {
			current = current.right
		}
	}

	newNode.left, newNode.right = nil, nil
	newNode.word = word
	u.size++

	if previous == nil {
		u.root = newNode
		return SUCCESS
	}

	if previous.word > word {
		previous.left = newNode
	} else {
		previous.right = newNode
	}

	return SUCCESS
}

// SearchCollection searches for the string in the WordCollection. This
// utilises a binary search algorithm (which is what the structure of this
// implementation is based around), and returns SUCCESS or FAILURE depending
// upon the outcome of the search.
func (u *UnbalBinarySearchTree) SearchCollection(needle string) int {
	current := u.root

	for current != nil {
		if current.word == needle {
			return SUCCESS
		} else if current.word > needle {
			current = current.left
		} else {
			current = current.right
		}
	}

	return FAILURE
}

// SizeCollection returns the number of words in the WordCollection.
func (u *UnbalBinarySearchTree) SizeCollection() int {
	return u.size
}

// DisplayCollection prints the contents of the WordCollection to the given
// Writer, utilising the recursive function displayBst.
func (u *UnbalBinarySearchTree) DisplayCollection(w io.Writer) {
	var element int
	u.root.displayBst(w, &element)
}

// displayBst is a recursive function that prints the words in the unbalanced
// binary search tree to the given Writer.
func (n *ubstNode) displayBst(w io.Writer, count *int) {
	if n != nil {
		n.left.displayBst(w, count)
		fmt.Fprintf(w, "Element %d:\t%s\n", *count, n.word)
		*count++
		n.right.displayBst(w, count)
	}
}

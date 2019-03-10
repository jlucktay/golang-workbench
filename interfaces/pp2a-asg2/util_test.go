package collection_test

import (
	"os"
)

const (
	SUCCESS = iota
	FAILURE
)

const (
	dictionary = "dictionary.txt"

	// A Tale of Two Cities, by Charles Dickens
	book1url = "https://www.gutenberg.org/files/98/98-0.txt"

	// Pride and Prejudice, by Jane Austen
	book2url = "http://www.gutenberg.org/files/1342/1342-0.txt"

	// Frankenstein, by Mary Wollstonecraft (Godwin) Shelley
	book3url = "https://www.gutenberg.org/files/84/84-0.txt"
)

// mustOpen helps with inlining benchmark functions
func mustOpen(filename string) (fp *os.File) {
	fp, errOpen := os.Open(filename)
	if errOpen != nil {
		panic(errOpen)
	}
	return
}

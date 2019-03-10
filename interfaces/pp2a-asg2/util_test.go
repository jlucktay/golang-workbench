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
	book1 = "98-0.txt"

	// Pride and Prejudice, by Jane Austen
	book2 = "1342-0.txt"

	// Frankenstein, by Mary Wollstonecraft (Godwin) Shelley
	book3 = "84-0.txt"
)

// mustOpen helps with inlining benchmark functions
func mustOpen(filename string) (fp *os.File) {
	fp, errOpen := os.Open(filename)
	if errOpen != nil {
		panic(errOpen)
	}
	return
}

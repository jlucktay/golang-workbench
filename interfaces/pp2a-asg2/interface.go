package collection

import "io"

const (
	SUCCESS = iota
	FAILURE
)

const WCSIZE = 250000

type WordCollection interface {
	MakeCollection() int
	FreeCollection()
	AddCollection(word string) int
	SearchCollection(word string) int
	SizeCollection() int
	DisplayCollection(stdout io.Writer)
}

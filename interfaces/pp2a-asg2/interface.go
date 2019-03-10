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
	AddCollection(string) int
	SearchCollection(string) int
	SizeCollection() int
	DisplayCollection(io.Writer)
}

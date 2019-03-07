package main

type WordCollection interface {
	MakeCollection() int
	FreeCollection()
	AddCollection(string) int
	SearchCollection(string) int
	SizeCollection() int
	DisplayCollection()
}

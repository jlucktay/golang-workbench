package main

import (
	"sync"
	"testing"
)

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

var pool = sync.Pool{
	New: func() any {
		return &Book{}
	},
}

func BenchmarkNoPool(b *testing.B) {
	var book *Book

	for b.Loop() {
		book = &Book{
			Title:  "The Art of Computer Programming, Vol. 1",
			Author: "Donald E. Knuth",
			Pages:  672,
		}
	}

	_ = book
}

func BenchmarkPool(b *testing.B) {
	for b.Loop() {
		book, ok := pool.Get().(*Book)
		if !ok {
			b.Fatal("casting by type was not OK")
		}

		book.Title = "The Art of Computer Programming, Vol. 1"
		book.Author = "Donald E. Knuth"
		book.Pages = 672

		pool.Put(book)
	}
}

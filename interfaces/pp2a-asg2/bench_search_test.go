package collection_test

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	p2 "github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2"
)

func BenchmarkSearchOAL(b *testing.B) {
	wc := &p2.OrdArrayLinear{}
	fillCollection(wc, mustOpen(dictionary), b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		searchCollection(wc, mustOpen(book1), b)
		searchCollection(wc, mustOpen(book2), b)
		searchCollection(wc, mustOpen(book3), b)
	}

	b.StopTimer()
	wc.FreeCollection()
}

func searchCollection(wc p2.WordCollection, book *os.File, b *testing.B) {
	defer book.Close()

	b.Logf("Reading '%s' (searching)...", book.Name())
	startTime := time.Now().UnixNano()
	found, notFound, wordTotal, lineTotal := 0, 0, 0, 0

	scanner := bufio.NewScanner(book)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		words := strings.FieldsFunc(strings.ToLower(scanner.Text()), split)

		for _, w := range words {
			if wc.SearchCollection(w) == SUCCESS {
				found++
			} else {
				notFound++
			}

			wordTotal++
		}

		lineTotal++
	}
	if errScan := scanner.Err(); errScan != nil {
		b.Fatal(errScan)
	}

	stopTime := time.Now().UnixNano()
	finalTime := stopTime - startTime

	b.Logf("%d words found on %d lines, %d words not found (total %d searches in %dμs)", found, lineTotal, notFound, wordTotal, finalTime/1e3)
}

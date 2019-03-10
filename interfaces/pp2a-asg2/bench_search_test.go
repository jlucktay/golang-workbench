package collection_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	p2 "github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2"
)

func BenchmarkSearchOAL(b *testing.B) {
	b.StopTimer()
	runSearchBenchmark(&p2.OrdArrayLinear{}, b)
}

func BenchmarkSearchOAB(b *testing.B) {
	b.StopTimer()
	runSearchBenchmark(&p2.OrdArrayBinary{}, b)
}

func runSearchBenchmark(wc p2.WordCollection, b *testing.B) {
	fillCollection(wc, mustOpen(dictionary), b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		searchCollection(wc, mustOpen(book1), b)
		searchCollection(wc, mustOpen(book2), b)
		searchCollection(wc, mustOpen(book3), b)
	}

	wc.FreeCollection()
}

// searchCollection takes three arguments:
// 1) an initialised WordCollection containing dictionary words
// 2) a pointer to an open text file descriptor
// 3) a pointer to a testing benchmark struct
//
// searchCollection searches the WordCollection for each of the words in the
// text file, where a 'word' is defined as what is delimited/tokenised on each
// line by the 'delims' constant and split() function.
// searchCollection also logs some timings of its own, in addition to the
// standard benchmark timings.
func searchCollection(wc p2.WordCollection, book *os.File, b *testing.B) {
	defer book.Close()

	found, notFound, wordTotal, lineTotal := 0, 0, 0, 0
	scanner := bufio.NewScanner(book)
	b.StartTimer() // code to be timed begins below here

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

	b.StopTimer() // timing ends here
	b.Logf("%d words found on %d lines, %d words not found (total %d)", found, lineTotal, notFound, wordTotal)
}

package collection_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	p2 "github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2"
)

func BenchmarkFillOAL(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		wc := &p2.OrdArrayLinear{}
		fillCollection(wc, mustOpen(dictionary), b)
		wc.FreeCollection()
	}
}

func BenchmarkFillOAB(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		wc := &p2.OrdArrayBinary{}
		fillCollection(wc, mustOpen(dictionary), b)
		wc.FreeCollection()
	}
}

func BenchmarkFillOLL(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		wc := &p2.OrdLinkedList{}
		fillCollection(wc, mustOpen(dictionary), b)
		wc.FreeCollection()
	}
}

// fillCollection takes three arguments:
// 1) a brand new uninitialised WordCollection
// 2) a pointer to an open text file descriptor
// 3) a pointer to a testing benchmark struct
//
// fillCollection fills the WordCollection with the words from the text file,
// and logs some timings of its own, in addition to the standard benchmark
// timings.
func fillCollection(wc p2.WordCollection, dictionary *os.File, b *testing.B) {
	defer dictionary.Close()

	if wc.MakeCollection() == FAILURE {
		b.Fatal("Unable to initialise WordCollection")
	}

	b.StartTimer() // code to be timed begins below here
	scanner := bufio.NewScanner(dictionary)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if wc.AddCollection(strings.ToLower(scanner.Text())) == FAILURE {
			b.Fatal("AddCollection failed")
		}
	}
	if errScan := scanner.Err(); errScan != nil {
		b.Fatal(errScan)
	}
	b.StopTimer() // timing ends here
}

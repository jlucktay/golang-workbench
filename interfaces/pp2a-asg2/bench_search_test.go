package collection_test

import (
	"os"
	"reflect"
	"testing"

	p2 "github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2"
)

func BenchmarkSearchOAL(b *testing.B) {
	wc := &p2.OrdArrayLinear{}
	fillCollection(wc, mustOpen(dictionary), b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		searchCollection(wc, nil, b)
	}

	b.StopTimer()
	wc.FreeCollection()
}

func searchCollection(wc p2.WordCollection, book *os.File, b *testing.B) {
	b.Logf("Implementation: %s", reflect.TypeOf(wc))
}

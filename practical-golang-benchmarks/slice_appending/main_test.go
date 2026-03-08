package main

import (
	"testing"
)

var numItems = 1000000

func BenchmarkSliceAppend(b *testing.B) {
	s := make([]byte, 0)

	i := 0

	for b.Loop() {
		s = append(s, 1) //nolint:staticcheck // Part of the purpose of this benchmark.

		i++
		if i == numItems {
			b.StopTimer()
			i = 0
			s = make([]byte, 0)
			b.StartTimer()
		}
	}
}

func BenchmarkSliceAppendPrealloc(b *testing.B) {
	s := make([]byte, 0, numItems)

	i := 0

	for b.Loop() {
		s = append(s, 1) //nolint:staticcheck // Part of the purpose of this benchmark.

		i++
		if i == numItems {
			b.StopTimer()
			i = 0
			s = make([]byte, 0, numItems)
			b.StartTimer()
		}
	}
}

package main

import (
	"math/rand"
	"sort"
	"testing"
)

func generateSlice(n int) []int {
	s := make([]int, 0, n)
	for range n {
		s = append(s, rand.Intn(1e9))
	}

	return s
}

func BenchmarkSort1000(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		s := generateSlice(1000)
		b.StartTimer()
		sort.Ints(s)
	}
}

func BenchmarkSort10000(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		s := generateSlice(10000)
		b.StartTimer()
		sort.Ints(s)
	}
}

func BenchmarkSort100000(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		s := generateSlice(100000)
		b.StartTimer()
		sort.Ints(s)
	}
}

func BenchmarkSort1000000(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		s := generateSlice(1000000)
		b.StartTimer()
		sort.Ints(s)
	}
}

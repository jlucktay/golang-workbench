package main

import (
	"math/rand"
	"strconv"
	"testing"
)

var NumItems = 1000000

func BenchmarkMapStringKeys(b *testing.B) {
	m := make(map[string]string)
	k := make([]string, 0)

	for i := 0; i < NumItems; i++ {
		key := strconv.Itoa(rand.Intn(NumItems))
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	i := 0
	l := len(m)

	for b.Loop() {
		if _, ok := m[k[i]]; ok { //nolint:staticcheck // Part of the purpose of this benchmark.
		}

		i++
		if i >= l {
			i = 0
		}
	}
}

func BenchmarkMapIntKeys(b *testing.B) {
	m := make(map[int]string)
	k := make([]int, 0)

	for i := 0; i < NumItems; i++ {
		key := rand.Intn(NumItems)
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	i := 0
	l := len(m)

	for b.Loop() {
		if _, ok := m[k[i]]; ok { //nolint:staticcheck // Part of the purpose of this benchmark.
		}

		i++
		if i >= l {
			i = 0
		}
	}
}

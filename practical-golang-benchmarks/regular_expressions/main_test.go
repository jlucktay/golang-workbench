package main

import (
	"regexp"
	"testing"
)

var testRegexp = `^[A-Za-z0-9._%+-][email protected][A-Za-z0-9.-]+\.[A-Za-z]+$`

func BenchmarkMatchString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := regexp.MatchString(testRegexp, "[email protected]")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMatchStringCompiled(b *testing.B) {
	r, err := regexp.Compile(testRegexp)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.MatchString("[email protected]")
	}
}

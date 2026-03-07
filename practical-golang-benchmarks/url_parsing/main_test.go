package main

import (
	"net/url"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	testURL := "https://www.example.com/path/file.html?param1=value1&param2=123"

	for b.Loop() {
		_, err := url.Parse(testURL)
		if err != nil {
			panic(err)
		}
	}
}

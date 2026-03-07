package main

import (
	"encoding/base64"
	"math/rand"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	data := make([]byte, 1024)
	rand.Read(data)

	for b.Loop() {
		base64.StdEncoding.EncodeToString([]byte(data))
	}
}

func BenchmarkDecode(b *testing.B) {
	data := make([]byte, 1024)
	rand.Read(data)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	for b.Loop() {
		_, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			panic(err)
		}
	}
}

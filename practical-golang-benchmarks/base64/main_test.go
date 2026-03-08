package main

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	data := make([]byte, 1024)

	if _, err := rand.Read(data); err != nil {
		b.Fatalf("reading random data: %v", err)
	}

	for b.Loop() {
		base64.StdEncoding.EncodeToString([]byte(data))
	}
}

func BenchmarkDecode(b *testing.B) {
	data := make([]byte, 1024)

	if _, err := rand.Read(data); err != nil {
		b.Fatalf("reading random data: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	for b.Loop() {
		_, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			panic(err)
		}
	}
}

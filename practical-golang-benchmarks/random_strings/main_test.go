package main

import (
	"crypto/rand"
	"testing"
)

// 64 letters
const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return bytes
}

func randomString(bytes []byte) string {
	for i, b := range bytes {
		bytes[i] = letters[b%64]
	}

	return string(bytes)
}

func BenchmarkCryptoRandString(b *testing.B) {
	for b.Loop() {
		randomString(randomBytes(16))
	}
}

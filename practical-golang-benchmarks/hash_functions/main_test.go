package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"hash"
	"testing"

	"golang.org/x/crypto/blake2b"
)

func benchmarkHash(b *testing.B, h hash.Hash) {
	data := make([]byte, 1024)

	if _, err := rand.Read(data); err != nil {
		b.Fatalf("reading random data: %v", err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		h.Write(data)
		h.Sum(nil)
	}
}

func BenchmarkMD5(b *testing.B) {
	benchmarkHash(b, md5.New())
}

func BenchmarkSHA1(b *testing.B) {
	benchmarkHash(b, sha1.New())
}

func BenchmarkSHA256(b *testing.B) {
	benchmarkHash(b, sha256.New())
}

func BenchmarkSHA512(b *testing.B) {
	benchmarkHash(b, sha512.New())
}

func BenchmarkSHA3256(b *testing.B) {
	benchmarkHash(b, hash.Hash(sha3.New256()))
}

func BenchmarkSHA3512(b *testing.B) {
	benchmarkHash(b, hash.Hash(sha3.New512()))
}

func BenchmarkBLAKE2b256(b *testing.B) {
	h, _ := blake2b.New256(nil)
	benchmarkHash(b, h)
}

func BenchmarkBLAKE2b512(b *testing.B) {
	h, _ := blake2b.New512(nil)
	benchmarkHash(b, h)
}

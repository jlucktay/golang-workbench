package prime

import (
	"math/big"
	"testing"
)

// const num = 65533051
const num = 1<<19 - 1

// const num = 2017133
func TestPrime(t *testing.T) {
	prime := uint32(num)
	if !i32Prime(prime) {
		t.Error(prime, "failed prime test")
	}
}

func TestPPrime(t *testing.T) {
	prime := big.NewInt(num)
	if !prime.ProbablyPrime(0) {
		t.Error(prime, "failed prime test")
	}
}

func TestIPrime(t *testing.T) {
	prime := uint64(num)
	if !isPrime(prime) {
		t.Error(prime, "failed prime test")
	}
}

func BenchmarkPrime(b *testing.B) {
	prime := uint32(num)
	for n := 0; n < b.N; n++ {
		i32Prime(prime)
	}
}

func BenchmarkPPrime(b *testing.B) {
	prime := big.NewInt(num)
	for n := 0; n < b.N; n++ {
		prime.ProbablyPrime(0)
	}
}

func BenchmarkIPrime(b *testing.B) {
	prime := uint64(num)
	for n := 0; n < b.N; n++ {
		isPrime(prime)
	}
}

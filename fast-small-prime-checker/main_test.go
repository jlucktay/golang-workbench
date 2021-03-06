package prime

import (
	"math"
	"math/big"
	"testing"
)

func i32Prime(n uint32) bool {
	//    if (n==2)||(n==3) {return true;}
	if n%2 == 0 {
		return false
	}
	if n%3 == 0 {
		return false
	}
	sqrt := uint32(math.Sqrt(float64(n)))
	for i := uint32(5); i <= sqrt; i += 6 {
		if n%i == 0 {
			return false
		}
		if n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func isPrime(value uint64) bool {
	for i := uint64(2); i <= uint64(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}

	return value > 1
}

//const num = 65533051
const num = 1<<19 - 1

//const num = 2017133

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

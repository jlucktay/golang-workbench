package prime

import "math"

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

// Find a problem at projecteuler.net then create the solution.
// Add a comment beneath your solution that includes a description of the problem.
// Upload your solution to github.
// Tweet me a link to your solution.
package main

import (
	"fmt"
	"math"
)

// TODO: learn concurrency and run this faster. Max out all of the cores!
func main() {
	var sum uint64

	for i := uint64(0); i < 2000000; i++ {
		if isPrime(i) {
			fmt.Print(i, "...")
			sum += i
		}
	}

	fmt.Println("\n\nSum of primes below 2 million:", sum)
}

// Blagged the basis for this function (and this function only!) with thanks from:
// https://www.thepolyglotdeveloper.com/2016/12/determine-number-prime-using-golang/
func isPrime(value uint64) bool {
	for i := uint64(2); i <= uint64(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}

	return value > 1
}

// Summation of primes
// Problem 10
// The sum of the primes below 10 is 2 + 3 + 5 + 7 = 17.
// Find the sum of all the primes below two million.

// Sum of primes below 2 million: 142913828922

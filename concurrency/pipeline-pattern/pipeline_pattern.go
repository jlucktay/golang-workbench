// An example of the pipeline pattern that I wrote as part of a training course.
package main

import (
	"math/rand"
	"time"

	"golang.org/x/text/message"
)

func main() {
	p := message.NewPrinter(message.MatchLanguage("en"))

	for x := range pipe(gen(20, 100)) {
		p.Printf("%42d\n", x)
	}

	// // Test the limits of what a uint64 can hold
	// for index := 1; index <= 32; index++ {
	// 	p.Printf("%2d!: %37d\n", index, factorial(uint64(index)))
	// }
	//
	// p.Printf("math.MaxUint64: %26d\n", uint64(math.MaxUint64))
}

func pipe(input chan uint64) chan uint64 {
	out := make(chan uint64)

	go func() {
		for j := range input {
			out <- factorial(j)
		}

		close(out)
	}()

	return out
}

func factorial(in uint64) uint64 {
	total := uint64(1)

	for i := in; i > 1; i-- {
		total *= i
	}

	return total
}

func gen(max int, num uint64) chan uint64 {
	out := make(chan uint64)

	go func() {
		for index := uint64(0); index < num; index++ {
			// Re-roll the random seed each time
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			out <- uint64(r.Intn(max) + 1) // Get [1,max] rather than [0,max)
		}

		close(out)
	}()

	return out
}

/*
CHALLENGE:
-- Change the code to execute 1,000 factorial computations concurrently and in parallel.
-- Use the "fan out / fan in" pattern
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// Generate 1,000 inputs for the factorial function, randomised from 1 up to 20
	in := gen(20, 1e3)

	// Get a slice of 16x channels ready to go to work
	chans := make([]<-chan int, 16)

	// Fan out
	for a := range chans {
		// Launch 16x goroutines to calculate factorials concurrently
		chans[a] = factorial(in)
	}

	fmt.Printf("Running on %v channels.\n", len(chans))

	// Fan in
	for n := range merge(chans...) {
		// Bring the results together into one stream
		fmt.Printf("%v/", n)
	}

	fmt.Printf("\nRan on %v channels.\n", len(chans))
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}

		wg.Done()
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func gen(max, num int) <-chan int {
	out := make(chan int, 100)

	go func() {
		for index := 0; index < num; index++ {
			// Re-roll the random seed each time
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			out <- r.Intn(max) + 1 // Get [1,max] rather than [0,max)
		}

		close(out)
	}()

	return out
}

func factorial(in <-chan int) <-chan int {
	out := make(chan int, 10)

	go func() {
		for n := range in {
			out <- fact(n)
		}

		close(out)
	}()

	return out
}

func fact(n int) int {
	total := 1

	for i := n; i > 0; i-- {
		total *= i
	}

	return total
}

// Playing around with this prime sieve example from the golang.org front page.
package main

import (
	"fmt"
	"math/big"

	"golang.org/x/text/message"
)

func main() {
	p := message.NewPrinter(message.MatchLanguage("en"))

	for x := range sieve(gen(1e3)) {
		p.Printf("%24d\n", x)
	}
}

func sieve(in <-chan int64) <-chan int64 {
	fmt.Println("sieve/outer: start")

	out := make(chan int64)

	go func() {
		fmt.Println("sieve/inner: start")

		for x := range in {
			i := big.NewInt(x)

			if i.ProbablyPrime(10) {
				out <- i.Int64()
			}
		}

		close(out)

		fmt.Println("sieve/inner: finish")
	}()

	fmt.Println("sieve/outer: finish")

	return out
}

func gen(limit int64) <-chan int64 {
	fmt.Println("gen/outer: start")

	out := make(chan int64)

	go func() {
		fmt.Println("gen/inner: start")

		for index := int64(0); index < limit; index++ {
			out <- index
		}

		close(out)

		fmt.Println("gen/inner: finish")
	}()

	fmt.Println("gen/outer: finish")

	return out
}

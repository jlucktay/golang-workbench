// In Go, all loops are 'for' loops. Whoa.
package main

import (
	"fmt"
)

func doForLoop() {
	sum := 0

	for i := 0; i < 10; i++ {
		sum += i
	}

	fmt.Println(sum)

	// init and post statements are optional
	sum = 1

	for sum < 1000 {
		sum += sum
	}

	fmt.Println(sum)
}

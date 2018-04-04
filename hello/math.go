// The 'math' package in the standard library contains lots of nice things.
package main

import (
	"fmt"
	"math"
)

func doMath() {
	fmt.Println(sqrt(2), sqrt(-4))

	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)

	fmt.Println("Sqrt:", mySqrt(2))
	fmt.Println("math:", math.Sqrt(2))
}

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}

	return fmt.Sprint(math.Sqrt(x))
}

func pow(x, n, lim float64) float64 {
	v := math.Pow(x, n)

	if v < lim {
		return v
	}

	fmt.Printf("%g >= %g\n", v, lim)

	return lim
}

func mySqrt(x float64) float64 {
	z := 1.0
	y := 0.0

	for i := 0; math.Abs(z-y) > 0.00000000001; i++ {
		y = z
		z -= (z*z - x) / (2 * z)

		fmt.Println(i, z)
	}

	return z
}

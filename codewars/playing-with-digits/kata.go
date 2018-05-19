package kata

import (
	"fmt"
	"math"
)

// DigPow ...
// Some numbers have funny properties. For example:
//
// 89 --> 8¹ + 9² = 89 * 1
//
// 695 --> 6² + 9³ + 5⁴= 1390 = 695 * 2
//
// 46288 --> 4³ + 6⁴+ 2⁵ + 8⁶ + 8⁷ = 2360688 = 46288 * 51
//
// Given a positive integer n written as abcd... (a, b, c, d... being digits) and a positive integer p we want to find a positive integer k, if it exists, such as the sum of the digits of n taken to the successive powers of p is equal to k * n. In other words:
//
// Is there an integer k such as : (a ^ p + b ^ (p+1) + c ^(p+2) + d ^ (p+3) + ...) = n * k
//
// If it is the case we will return k, if not return -1.
//
// Note: n, p will always be given as strictly positive integers.
func DigPow(n, p uint) int {
	fmt.Printf("splitIntoDigits(%v) == %v\n", n, splitIntoDigits(n))
	fmt.Printf("successivePowers(%v, %v) == %v\n", n, p, successivePowers(n, p))

	return 0
}

func splitIntoDigits(in uint) (out []uint) {
	out = make([]uint, 0, 1)

	for ; in > 0; in /= 10 {
		out = append(out, in%10)
	}

	for i := len(out)/2 - 1; i >= 0; i-- {
		opp := len(out) - 1 - i
		out[i], out[opp] = out[opp], out[i]
	}

	return
}

func successivePowers(n, p uint) (total uint) {
	digits := splitIntoDigits(n)
	exponent := float64(p)
	// fmt.Println("n:", n)

	for _, digit := range digits {
		// fmt.Println("total:", total)
		// fmt.Println("digit:", digit)
		// fmt.Println("exponent:", exponent)
		total += uint(math.Pow(float64(digit), exponent))
		exponent++
	}

	return
}

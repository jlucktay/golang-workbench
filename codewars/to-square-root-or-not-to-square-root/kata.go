// Package kata - https://www.codewars.com/kata/to-square-root-or-not-to-square-root/train/go
package kata

import (
	"math"
)

// SquareOrSquareRoot is a method that will get an integer array as parameter and will process every number from this array.
// Return a new array where, if the number has an integer square root, take this, otherwise square the number.
// [4,3,9,7,2,1] -> [2,9,3,49,4,1]
// The input array will always contain only positive numbers and will never be empty or null.
// The input array should not be modified!
func SquareOrSquareRoot(arr []int) []int {
	out := make([]int, 0, 1)

	for _, value := range arr {
		if result := math.Sqrt(float64(value)); result == math.Floor(result) {
			out = append(out, int(math.Sqrt(float64(value))))
		} else {
			out = append(out, value*value)
		}
	}

	return out
}

package kata

import (
	"strings"
)

// Write a function that will return the count of distinct case-insensitive alphabetic characters and numeric digits that occur more than once in the input string.
// The input string can be assumed to contain only alphabets (both uppercase and lowercase) and numeric digits.
func duplicateCount(s string) (result int) {
	parse := map[rune]uint{}

	for _, r := range strings.ToLower(s) {
		if parse[r]++; parse[r] == 2 {
			result++
		}
	}

	return
}

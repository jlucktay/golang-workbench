package kata

import (
	"unicode"
)

// Write a function that will return the count of distinct case-insensitive alphabetic characters and numeric digits that occur more than once in the input string.
// The input string can be assumed to contain only alphabets (both uppercase and lowercase) and numeric digits.
func duplicateCount(s1 string) int {
	parse := map[rune]uint{}

	for _, r := range []rune(s1) {
		parse[unicode.ToLower(r)]++
	}

	result := 0

	for _, x := range parse {
		if x > 1 {
			result++
		}
	}

	return result
}

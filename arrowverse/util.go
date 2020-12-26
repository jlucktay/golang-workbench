package main

import "unicode"

// mapSpaces helps us get rid of non-breaking spaces from HTML.
func mapSpaces(input rune) rune {
	if unicode.IsSpace(input) {
		return ' '
	}

	return input
}

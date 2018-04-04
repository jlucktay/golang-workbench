// Package fizzbuzz implements a basic function for this lovely interview question.
package fizzbuzz

import (
	"errors"
	"strconv"
)

// Fizzbuzz: See https://wikipedia.org/wiki/Fizz_buzz
func fizzbuzz(number int) (string, error) {
	if number <= 0 {
		return "", errors.New("undefined")
	}

	switch {
	case number%15 == 0:
		return "FizzBuzz", nil
	case number%3 == 0:
		return "Fizz", nil
	case number%5 == 0:
		return "Buzz", nil
	}

	return strconv.Itoa(number), nil
}

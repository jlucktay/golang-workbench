package main

import (
	"testing"
)

func TestFactorial(t *testing.T) {
	testCases := []struct {
		desc string
		in   uint64
		want uint64
	}{
		{"0!", 0, 1},
		{"1!", 1, 1},
		{"2!", 2, 2},
		{"3!", 3, 6},
		{"4!", 4, 24},
		{"5!", 5, 120},
		{"6!", 6, 720},
		{"7!", 7, 5040},
		{"8!", 8, 40320},
		{"9!", 9, 362880},
		{"10!", 10, 3628800},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if result := factorial(tC.in); result != tC.want {
				t.Fatalf("factorial(%v) == '%v', wanted '%v'", tC.in, result, tC.want)
			}
		})
	}
}

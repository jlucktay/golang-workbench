package main

import "testing"

func TestCountingValleys(t *testing.T) {
	testCases := []struct {
		n, want int32
		s       string
	}{
		{
			n:    8,
			s:    "UDDDUDUU",
			want: 1,
		},
		{
			n:    12,
			s:    "DDUUDDUDUUUD",
			want: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.s, func(t *testing.T) {
			if result := countingValleys(tC.n, tC.s); result != tC.want {
				t.Fatalf("countingValleys(%v, %v) == '%v', wanted '%v'", tC.n, tC.s, result, tC.want)
			}
		})
	}
}

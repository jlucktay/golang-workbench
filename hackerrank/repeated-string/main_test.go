package main

import (
	"fmt"
	"testing"
)

func TestRepeatedString(t *testing.T) {
	testCases := []struct {
		s       string
		n, want int64
	}{
		{
			s:    "aba",
			n:    10,
			want: 7,
		},
		{
			s:    "a",
			n:    1000000000000,
			want: 1000000000000,
		},
		{
			s:    "epsxyyflvrrrxzvnoenvpegvuonodjoxfwdmcvwctmekpsnamchznsoxaklzjgrqruyzavshfbmuhdwwmpbkwcuomqhiyvuztwvq",
			n:    549382313570,
			want: 16481469408,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v x %v", tC.s, tC.n), func(t *testing.T) {
			if result := repeatedString(tC.s, tC.n); result != tC.want {
				t.Fatalf("repeatedString(%v, %v) == '%v', wanted '%v'", tC.s, tC.n, result, tC.want)
			}
		})
	}
}

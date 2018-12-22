package main

import (
	"fmt"
	"testing"
)

func TestJumpingOnClouds(t *testing.T) {
	testCases := []struct {
		c    []int32
		want int32
	}{
		{
			c:    []int32{0, 0, 1, 0, 0, 1, 0},
			want: 4,
		},
		{
			c:    []int32{0, 0, 0, 1, 0, 0},
			want: 3,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.c), func(t *testing.T) {
			if result := jumpingOnClouds(tC.c); result != tC.want {
				t.Fatalf("countingValleys(%v) == '%v', wanted '%v'", tC.c, result, tC.want)
			}
		})
	}
}

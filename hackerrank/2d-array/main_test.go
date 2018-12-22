package main

import (
	"fmt"
	"testing"
)

func TestHourglassSum(t *testing.T) {
	testCases := []struct {
		arr  [][]int32
		want int32
	}{
		{
			arr: [][]int32{
				{1, 1, 1, 0, 0, 0},
				{0, 1, 0, 0, 0, 0},
				{1, 1, 1, 0, 0, 0},
				{0, 0, 2, 4, 4, 0},
				{0, 0, 0, 2, 0, 0},
				{0, 0, 1, 2, 4, 0},
			},
			want: 19,
		},
		{
			arr: [][]int32{
				{1, 1, 1, 0, 0, 0},
				{0, 1, 0, 0, 0, 0},
				{1, 1, 1, 0, 0, 0},
				{0, 9, 2, -4, -4, 0},
				{0, 0, 0, -2, 0, 0},
				{0, 0, -1, -2, -4, 0},
			},
			want: 13,
		},
		{
			arr: [][]int32{
				{-9, -9, -9, 1, 1, 1},
				{0, -9, 0, 4, 3, 2},
				{-9, -9, -9, 1, 2, 3},
				{0, 0, 8, 6, 6, 0},
				{0, 0, 0, -2, 0, 0},
				{0, 0, 1, 2, 4, 0},
			},
			want: 28,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.arr), func(t *testing.T) {
			if result := hourglassSum(tC.arr); result != tC.want {
				t.Fatalf("hourglassSum(%v) == '%v', wanted '%v'", tC.arr, result, tC.want)
			}
		})
	}
}

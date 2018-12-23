package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMinimumSwaps(t *testing.T) {
	testCases := []struct {
		arr  []int32
		want int32
	}{
		{
			arr:  []int32{4, 3, 1, 2},
			want: 3,
		},
		{
			arr:  []int32{2, 3, 4, 1, 5},
			want: 3,
		},
		{
			arr:  []int32{1, 3, 5, 2, 4, 6, 7},
			want: 3,
		},
	}

	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.arr), func(t *testing.T) {
			if result := minimumSwaps(tC.arr); !reflect.DeepEqual(result, tC.want) {
				t.Fatalf("minimumSwaps(%v) == '%v', wanted '%v'", tC.arr, result, tC.want)
			}
		})
	}
}

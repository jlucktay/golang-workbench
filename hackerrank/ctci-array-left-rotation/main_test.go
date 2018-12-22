package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRotLeft(t *testing.T) {
	testCases := []struct {
		a    []int32
		d    int32
		want []int32
	}{
		{
			a:    []int32{1, 2, 3, 4, 5},
			d:    1,
			want: []int32{2, 3, 4, 5, 1},
		},
		{
			a:    []int32{1, 2, 3, 4, 5},
			d:    4,
			want: []int32{5, 1, 2, 3, 4},
		},
		{
			a:    []int32{41, 73, 89, 7, 10, 1, 59, 58, 84, 77, 77, 97, 58, 1, 86, 58, 26, 10, 86, 51},
			d:    10,
			want: []int32{77, 97, 58, 1, 86, 58, 26, 10, 86, 51, 41, 73, 89, 7, 10, 1, 59, 58, 84, 77},
		},
		{
			a:    []int32{33, 47, 70, 37, 8, 53, 13, 93, 71, 72, 51, 100, 60, 87, 97},
			d:    13,
			want: []int32{87, 97, 33, 47, 70, 37, 8, 53, 13, 93, 71, 72, 51, 100, 60},
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v <- %v", tC.a, tC.d), func(t *testing.T) {
			if result := rotLeft(tC.a, tC.d); !reflect.DeepEqual(result, tC.want) {
				t.Fatalf("rotLeft(%v, %v) == '%v', wanted '%v'", tC.a, tC.d, result, tC.want)
			}
		})
	}
}

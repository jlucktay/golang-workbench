package kata

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSquareOrSquareRoot(t *testing.T) {
	tests := [][][]int{
		{{4, 3, 9, 7, 2, 1}, {2, 9, 3, 49, 4, 1}},
		{{100, 101, 5, 5, 1, 1}, {10, 10201, 25, 25, 1, 1}},
		{{1, 2, 3, 4, 5, 6}, {1, 4, 9, 2, 25, 36}},
	}

	for i, tC := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			if result := SquareOrSquareRoot(tC[0]); !reflect.DeepEqual(result, tC[1]) {
				t.Fatalf("SquareOrSquareRoot(%v) == '%v', wanted '%v'", tC[0], result, tC[1])
			}
		})
	}
}

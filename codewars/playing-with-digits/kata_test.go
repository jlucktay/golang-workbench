package kata

import (
	"fmt"
	"testing"
)

// digPow(89, 1) should return 1 since 8¹ + 9² = 89 = 89 * 1
// digPow(92, 1) should return -1 since there is no k such as 9¹ + 2² equals 92 * k
// digPow(695, 2) should return 2 since 6² + 9³ + 5⁴= 1390 = 695 * 2
// digPow(46288, 3) should return 51 since 4³ + 6⁴+ 2⁵ + 8⁶ + 8⁷ = 2360688 = 46288 * 51

func Test(t *testing.T) {
	t.Skip()

	testCases := []struct {
		n uint
		p uint
		k int
	}{
		{89, 1, 1},
		{92, 1, -1},
		{695, 2, 2},
		{46288, 3, 51},
	}

	for _, tC := range testCases {
		t.Run(fmt.Sprintf("n: %v, p: %v, k: %v", tC.n, tC.p, tC.k), func(t *testing.T) {
			if result := DigPow(tC.n, tC.p); result != tC.k {
				t.Fatalf("DigPow(%v, %v) == '%v', wanted '%v'", tC.n, tC.p, result, tC.k)
			}
		})
	}
}

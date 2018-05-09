package kata

import (
	"reflect"
	"testing"
)

func TestNbMonths(t *testing.T) {
	testCases := []struct {
		desc               string
		startPriceOld      int
		startPriceNew      int
		savingperMonth     int
		percentLossByMonth float64
		exp                [2]int
	}{
		{"One", 2000, 8000, 1000, 1.5, [2]int{6, 766}},
		{"Two", 12000, 8000, 1000, 1.5, [2]int{0, 4000}},
		{"Three", 8000, 12000, 500, 1.0, [2]int{8, 597}},
		{"Four", 18000, 32000, 1500, 1.25, [2]int{8, 332}},
		{"Five", 7500, 32000, 300, 1.55, [2]int{25, 122}},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ans := NbMonths(tC.startPriceOld, tC.startPriceNew, tC.savingperMonth, tC.percentLossByMonth)

			if !reflect.DeepEqual(ans, tC.exp) {
			}
		})
	}
}

package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestStorage(t *testing.T) {
	testCases := []struct {
		desc string
		ps   PaymentStorage
	}{
		{
			desc: "Dummy storage in a map; won't persist across app restarts",
			ps:   &dummyStorage{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.ps))
			i := is.New(t)

			i.NoErr(tC.ps.Init())

			dummyPay := &Payment{
				Amount: decimal.New(12345, 1),
			}

			// C
			newId, errCreate := tC.ps.Create(*dummyPay)
			i.NoErr(errCreate)

			// R
			readPay, errRead := tC.ps.Read(newId)
			i.NoErr(errRead)

			i.True(reflect.DeepEqual(dummyPay, &readPay))

			// U
			dummyPay.Reference = "ref"
			i.NoErr(tC.ps.Update(newId, *dummyPay))
			updatedPay, _ := tC.ps.Read(newId)
			i.True(reflect.DeepEqual(dummyPay, &updatedPay))

			// D
			i.NoErr(tC.ps.Delete(newId))
			_, errDeleted := tC.ps.Read(newId)
			i.Equal(
				errDeleted.Error(),
				fmt.Sprintf("Payment ID '%s' not found.", newId),
			)
		})
	}
}

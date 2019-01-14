package main

import (
	"reflect"
	"testing"
)

func newDSM() dummyStoreMap {
	return dummyStoreMap{
		map[string]uint64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
}

func TestCreateDSM(t *testing.T) {
	testDSM := newDSM()
	e := testDSM.Create(widget{"four", 4})
	if e != nil {
		t.Fatalf("error '%v' from Create()", e)
	}

	targetDSM := dummyStoreMap{
		map[string]uint64{
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  4,
		},
	}
	if !reflect.DeepEqual(testDSM, targetDSM) {
		t.Fatalf("'%v' and '%v' don't match", testDSM, targetDSM)
	}
}

func TestReadDSM(t *testing.T) {
	testDSM := newDSM()
	if r, e := testDSM.Read("one"); r.data != 1 || e != nil {
		t.Fatalf("Unexpected '%v' and/or '%v'", r, e)
	}
}

func TestUpdateDSM(t *testing.T) {
	testDSM := newDSM()
	e := testDSM.Update("one", 11)
	if e != nil {
		t.Fatalf("error '%v' from Update()", e)
	}

	targetDSM := dummyStoreMap{
		map[string]uint64{
			"one":   11,
			"two":   2,
			"three": 3,
		},
	}
	if !reflect.DeepEqual(testDSM, targetDSM) {
		t.Fatalf("'%v' and '%v' don't match", testDSM, targetDSM)
	}
}

func TestDeleteDSM(t *testing.T) {
	testDSM := newDSM()
	e := testDSM.Delete("two")
	if e != nil {
		t.Fatalf("error '%v' from Delete()", e)
	}

	targetDSM := dummyStoreMap{
		map[string]uint64{
			"one":   1,
			"three": 3,
		},
	}
	if !reflect.DeepEqual(testDSM, targetDSM) {
		t.Fatalf("'%v' and '%v' don't match", testDSM, targetDSM)
	}
}

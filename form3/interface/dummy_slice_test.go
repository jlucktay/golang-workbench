package main

import (
	"reflect"
	"testing"
)

func newDSS() dummyStoreSlice {
	return dummyStoreSlice{
		[]widget{
			{"one", 1},
			{"two", 2},
			{"three", 3},
		},
	}
}

func TestCreateDSS(t *testing.T) {
	actual := newDSS()
	e := actual.Create(widget{"four", 4})
	if e != nil {
		t.Fatalf("error '%v' from Create()", e)
	}

	expected := dummyStoreSlice{
		[]widget{
			{"one", 1},
			{"two", 2},
			{"three", 3},
			{"four", 4},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Create() failed!\nExpected: %v\nActual %v", expected, actual)
	}
}

func TestReadDSS(t *testing.T) {
	testDSS := newDSS()
	if r, e := testDSS.Read("one"); r.data != 1 || e != nil {
		t.Fatalf("Unexpected '%v' and/or '%v'", r, e)
	}
}

func TestUpdateDSS(t *testing.T) {
	actual := newDSS()
	e := actual.Update("one", 11)
	if e != nil {
		t.Fatalf("error '%v' from Update()", e)
	}

	expected := dummyStoreSlice{
		[]widget{
			{"one", 11},
			{"two", 2},
			{"three", 3},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Update() failed!\nExpected: %v\nActual %v", expected, actual)
	}
}

func TestDeleteDSS(t *testing.T) {
	actual := newDSS()
	e := actual.Delete("two")
	if e != nil {
		t.Fatalf("error '%v' from Delete()", e)
	}

	expected := dummyStoreSlice{
		[]widget{
			{"one", 1},
			{"three", 3},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Delete() failed!\nExpected: %v\nActual %v", expected, actual)
	}
}

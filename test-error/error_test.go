package main

import "testing"

// Testing to get error message
func TestReturnSomeErr(t *testing.T) {
	Expected := "this is error message"
	actual := returnSomeErr(-1)

	if actual.Error() != Expected {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, Expected)
	}
}

package testerr_test

import (
	"testing"

	testerr "go.jlucktay.dev/golang-workbench/test-error"
)

// Testing to get error message
func TestReturnSomeErr(t *testing.T) {
	Expected := "this is error message"
	actual := testerr.ReturnSomeErr(-1)

	if actual.Error() != Expected {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, Expected)
	}
}

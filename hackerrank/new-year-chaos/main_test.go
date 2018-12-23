package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMinimumBribes(t *testing.T) {
	testCases := []struct {
		q    []int32
		want string
	}{
		{
			q:    []int32{2, 1, 5, 3, 4},
			want: "3\n",
		},
		{
			q:    []int32{2, 5, 1, 3, 4},
			want: "Too chaotic\n",
		},
		{
			q:    []int32{5, 1, 2, 3, 7, 8, 6, 4},
			want: "Too chaotic\n",
		},
		{
			q:    []int32{1, 2, 5, 3, 7, 8, 6, 4},
			want: "7\n",
		},
		{
			q:    []int32{1, 2, 5, 3, 4, 7, 8, 6},
			want: "4\n",
		},
	}

	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.q), func(t *testing.T) {
			if result := catchStdOut(t, minimumBribes, tC.q); result != tC.want {
				t.Fatalf("minimumBribes(%v) == '%v', wanted '%v'",
					tC.q, strings.TrimSpace(result), strings.TrimSpace(tC.want))
			}
		})
	}
}

// Thank you:
// https://groups.google.com/forum/#!topic/golang-nuts/hVUtoeyNL7Y
//
// Returns output to `os.Stdout` from `runnable` as string.
func catchStdOut(t *testing.T, runnable func([]int32), input []int32) string {
	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()

	r, fakeStdout, errPipeOpen := os.Pipe()
	if errPipeOpen != nil {
		t.Fatal(errPipeOpen)
	}

	os.Stdout = fakeStdout
	runnable(input)

	// need to close here, otherwise ReadAll never gets "EOF".
	if errStdoutClose := fakeStdout.Close(); errStdoutClose != nil {
		t.Fatal(errStdoutClose)
	}

	newOutBytes, errRead := ioutil.ReadAll(r)
	if errRead != nil {
		t.Fatal(errRead)
	}

	if errPipeClose := r.Close(); errPipeClose != nil {
		t.Fatal(errPipeClose)
	}

	return string(newOutBytes)
}

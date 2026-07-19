package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	var out bytes.Buffer
	if err := run(&out); err != nil {
		t.Errorf("run() returned an error: %v", err)
	}

	expected := "Hello from sekigochi!\n"
	if out.String() != expected {
		t.Errorf("run() output = %q, want %q", out.String(), expected)
	}
}

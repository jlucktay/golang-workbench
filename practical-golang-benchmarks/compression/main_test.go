package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	data, err := os.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	zw := gzip.NewWriter(io.Discard)

	for b.Loop() {
		_, err = zw.Write(data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	data, err := os.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(data)
	if err != nil {
		panic(err)
	}

	err = zw.Close()
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(buf.Bytes())
	zr, _ := gzip.NewReader(r)

	for b.Loop() {
		r.Reset(buf.Bytes())

		if err := zr.Reset(r); err != nil {
			panic(err)
		}

		_, err := io.ReadAll(zr)
		if err != nil {
			panic(err)
		}
	}
}

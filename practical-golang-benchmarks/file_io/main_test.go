package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkWriteFile(b *testing.B) {
	for b.Loop() {
		f, err := os.Create(filepath.Clean("./test.txt"))
		if err != nil {
			panic(err)
		}

		for range 100000 {
			f.WriteString("some text!\n")
		}

		f.Close()
	}
}

func BenchmarkWriteFileBuffered(b *testing.B) {
	for b.Loop() {
		f, err := os.Create(filepath.Clean("./test.txt"))
		if err != nil {
			panic(err)
		}

		w := bufio.NewWriter(f)

		for range 100000 {
			w.WriteString("some text!\n")
		}

		w.Flush()
		f.Close()
	}
}

func BenchmarkReadFile(b *testing.B) {
	for b.Loop() {
		f, err := os.Open(filepath.Clean("./test.txt"))
		if err != nil {
			panic(err)
		}

		b := make([]byte, 10)

		_, err = f.Read(b)
		for err == nil {
			_, err = f.Read(b)
		}
		if err != io.EOF {
			panic(err)
		}

		f.Close()
	}
}

func BenchmarkReadFileBuffered(b *testing.B) {
	for b.Loop() {
		f, err := os.Open(filepath.Clean("./test.txt"))
		if err != nil {
			panic(err)
		}

		r := bufio.NewReader(f)

		_, err = r.ReadString('\n')
		for err == nil {
			_, err = r.ReadString('\n')
		}
		if err != io.EOF {
			panic(err)
		}

		f.Close()
	}
}

package ex1_3_test

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func argsJoin(w io.Writer, args []string) {
	fmt.Fprintf(w, strings.Join(args[0:], " "))
}

func argsRange(w io.Writer, args []string) {
	for i := 0; i < len(args); i++ {
		fmt.Fprintf(w, "%d: %s\n", i, args[i])
	}
}

func benchmarkArgs(b *testing.B, size int, benchMe func(io.Writer, []string)) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		a := generateSlice(size)
		buf := &bytes.Buffer{}
		b.StartTimer()
		benchMe(buf, a)
	}
}

func BenchmarkArgsJoin2(b *testing.B)      { benchmarkArgs(b, 2, argsJoin) }
func BenchmarkArgsJoin10(b *testing.B)     { benchmarkArgs(b, 10, argsJoin) }
func BenchmarkArgsJoin100(b *testing.B)    { benchmarkArgs(b, 100, argsJoin) }
func BenchmarkArgsJoin1000(b *testing.B)   { benchmarkArgs(b, 1000, argsJoin) }
func BenchmarkArgsJoin10000(b *testing.B)  { benchmarkArgs(b, 10000, argsJoin) }
func BenchmarkArgsJoin100000(b *testing.B) { benchmarkArgs(b, 100000, argsJoin) }

func BenchmarkArgsRange2(b *testing.B)      { benchmarkArgs(b, 2, argsRange) }
func BenchmarkArgsRange10(b *testing.B)     { benchmarkArgs(b, 10, argsRange) }
func BenchmarkArgsRange100(b *testing.B)    { benchmarkArgs(b, 100, argsRange) }
func BenchmarkArgsRange1000(b *testing.B)   { benchmarkArgs(b, 1000, argsRange) }
func BenchmarkArgsRange10000(b *testing.B)  { benchmarkArgs(b, 10000, argsRange) }
func BenchmarkArgsRange100000(b *testing.B) { benchmarkArgs(b, 100000, argsRange) }

func generateSlice(n int) []string {
	s := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, RandStringBytesMaskImprSrcSB(rand.Intn(20)))
	}
	return s
}

package stringconcat

import "testing"

const TestString = "test"

func benchmarkConcat(size int, SelfConcat func(string, int) string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		SelfConcat(TestString, size)
	}
}

func BenchmarkConcatOperator2(b *testing.B)      { benchmarkConcat(2, selfConcatOperator, b) }
func BenchmarkConcatOperator10(b *testing.B)     { benchmarkConcat(10, selfConcatOperator, b) }
func BenchmarkConcatOperator100(b *testing.B)    { benchmarkConcat(100, selfConcatOperator, b) }
func BenchmarkConcatOperator1000(b *testing.B)   { benchmarkConcat(1000, selfConcatOperator, b) }
func BenchmarkConcatOperator10000(b *testing.B)  { benchmarkConcat(10000, selfConcatOperator, b) }
func BenchmarkConcatOperator100000(b *testing.B) { benchmarkConcat(100000, selfConcatOperator, b) }

func BenchmarkConcatBuffer2(b *testing.B)      { benchmarkConcat(2, selfConcatBuffer, b) }
func BenchmarkConcatBuffer10(b *testing.B)     { benchmarkConcat(10, selfConcatBuffer, b) }
func BenchmarkConcatBuffer100(b *testing.B)    { benchmarkConcat(100, selfConcatBuffer, b) }
func BenchmarkConcatBuffer1000(b *testing.B)   { benchmarkConcat(1000, selfConcatBuffer, b) }
func BenchmarkConcatBuffer10000(b *testing.B)  { benchmarkConcat(10000, selfConcatBuffer, b) }
func BenchmarkConcatBuffer100000(b *testing.B) { benchmarkConcat(100000, selfConcatBuffer, b) }

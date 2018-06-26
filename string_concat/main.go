package stringconcat

import (
	"bytes"
)

func concatOperator(original *string, concat string) {
	// This could be written as 'return *original + concat' but wanted to confirm no special
	// compiler optimizations existed for concatenating a string to itself.
	*original = *original + concat
}

func selfConcatOperator(input string, n int) string {
	output := ""
	for i := 0; i < n; i++ {
		concatOperator(&output, input)
	}
	return output
}

func concatBuffer(original *bytes.Buffer, concat string) {
	original.WriteString(concat)
}

func selfConcatBuffer(input string, n int) string {
	var output bytes.Buffer
	for i := 0; i < n; i++ {
		concatBuffer(&output, input)
	}
	return string(output.Bytes())
}

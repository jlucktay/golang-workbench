// Spitting out all available environment variables
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for _, e := range os.Environ() {
		v := strings.Split(e, "=")

		// Variables containing 'termcap' have some annoying colour codes, so we filter them out
		if !strings.Contains(strings.ToLower(v[0]), "termcap") {
			fmt.Printf("'%s' = '%s'\n", v[0], v[1:])
		}
	}
}

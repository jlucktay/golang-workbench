// TODO: replace with your own tests (TDD). An example to get you started is included below.
// Ginkgo BDD Testing Framework <http://onsi.github.io/ginkgo></http:>
// Gomega Matcher Library <http://onsi.github.io/gomega></http:>

package kata

import (
	"testing"
)

func Test(t *testing.T) {
	testCases := []struct {
		in   string
		want int
	}{
		{"aA11", 2},
		{"aabbcde", 2},
		{"aabBcde", 2},
		{"ABBA", 2},
		{"abcde", 0},
		{"abcdea", 1},
		{"abcdeaB11", 3},
		{"Indivisibilities", 2},
		{"indivisibility", 1},
	}
	for _, tC := range testCases {
		t.Run(tC.in, func(t *testing.T) {
			if result := duplicateCount(tC.in); result != tC.want {
				t.Fatalf("duplicateCount(%v) == '%v', wanted '%v'", tC.in, result, tC.want)
			}
		})
	}
}

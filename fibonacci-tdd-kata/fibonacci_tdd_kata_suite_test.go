package fibonacci_tdd_kata_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFibonacciTddKata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FibonacciTddKata Suite")
}

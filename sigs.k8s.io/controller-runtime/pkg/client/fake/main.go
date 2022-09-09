package main

import (
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func main() {
	_ = fake.NewFakeClient()
}

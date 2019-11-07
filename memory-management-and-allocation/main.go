package main

type smallStruct struct {
	a, b int64
	c, d float64
}

func main() {
	smallAllocation()
}

// The annotation //go:noinline will disable in-lining that would optimize the code by removing the function and,
// therefore, end up with no allocation.

//go:noinline
func smallAllocation() *smallStruct {
	return &smallStruct{}
}

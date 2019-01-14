package main

type widget struct {
	name string
	data uint64
}

func main() {
	ms := dummyStoreMap{
		w: map[string]uint64{},
	}

	_ = ms
}

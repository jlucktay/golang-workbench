package main

//go-sumtype:decl MySumType

type MySumType interface {
	sealed()
}

type VariantA struct{}

func (*VariantA) sealed() {}

type VariantB struct{}

func (*VariantB) sealed() {}

func main() {
	switch MySumType(nil).(type) {
	case *VariantA:
	}
}

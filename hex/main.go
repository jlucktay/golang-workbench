package main

import "fmt"

func main() {
	fmt.Println("Hex math")

	a := 0xffff

	fmt.Printf("a (X): %X\n", a)
	fmt.Printf("a (+v): %+v\n", a)

	b := a + a + a + a

	fmt.Printf("b (X): %X\n", b)
	fmt.Printf("b (+v): %+v\n", b)

	c := b / 4

	fmt.Printf("c (X): %X\n", c)
	fmt.Printf("c (+v): %+v\n", c)

	max8 := ^uint8(0)

	fmt.Printf("max8 (X): %X\n", max8)
	fmt.Printf("max8 (+v): %+v\n", max8)

	max8squared := (256 * 256) - 1

	fmt.Printf("max8squared (X): %X\n", max8squared)
	fmt.Printf("max8squared (+v): %+v\n", max8squared)

	max16 := ^uint16(0)

	fmt.Printf("max16 (X): %X\n", max16)
	fmt.Printf("max16 (+v): %+v\n", max16)

	max16over256 := ^uint16(0) / 256

	fmt.Printf("max16over256 (X): %X\n", max16over256)
	fmt.Printf("max16over256 (+v): %+v\n", max16over256)

	max32 := ^uint32(0)

	fmt.Printf("max32 (X): %X\n", max32)
	fmt.Printf("max32 (+v): %+v\n", max32)

	uint32over8 := max32 / uint32(max8)

	fmt.Printf("uint32over8 (X): %X\n", uint32over8)
	fmt.Printf("uint32over8 (+v): %+v\n", uint32over8)

	uint32over1024 := max32 / 1024

	fmt.Printf("uint32over1024 (X): %X\n", uint32over1024)
	fmt.Printf("uint32over1024 (+v): %+v\n", uint32over1024)
}

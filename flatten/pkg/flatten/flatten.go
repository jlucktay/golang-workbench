package flatten

import "fmt"

// Flatten will take a slice of arbitrarily nested (slices of) interfaces and/or ints, and return a flat slice of ints.
func Flatten(input interface{}) []int {
	result := make([]int, 0)

	switch typeSwitch := input.(type) {

	// Straightforward integer append
	case int:
		result = append(result, typeSwitch)

	// Also fairly straightforward
	case []int:
		result = append(result, typeSwitch...)

	// Nested with the same type; make recursive call(s)
	case []interface{}:
		for _, data := range typeSwitch {
			result = append(result, Flatten(data)...)
		}

	// Not handling other types, beyond these three
	default:
		fmt.Printf("Unknown type '%[1]T': %#[1]v\n", typeSwitch)
	}

	return result
}

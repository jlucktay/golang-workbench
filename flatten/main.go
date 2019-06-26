package main

import "fmt"

func main() {
	input := []interface{}{
		[]interface{}{
			1,
			2,
			[]int{3},
		},
		4,
	}
	fmt.Printf("%#v\n", flattenSlice(input))
}

// [[1,2,[3]],4] -> [1,2,3,4]

func flattenSlice(input []interface{}) []int {
	result := make([]int, 0)

	for _, element := range input {
		switch element.(type) {

		// Straight integer append
		case int:
			result = append(result, element.(int))

		// Need to massage this somewhat, as the data behind the two different types (int and interface) of slices has
		// a different size
		case []int:
			interfaceSlice := make([]interface{}, len(element.([]int)))

			for index, data := range element.([]int) {
				interfaceSlice[index] = data
			}

			result = append(result, flattenSlice(interfaceSlice)...)

		// Nested with the same type, so just make the recursive call
		case []interface{}:
			result = append(result, flattenSlice(element.([]interface{}))...)

		default:
			fmt.Printf("Unknown type '%[1]T': %#[1]v\n", element)
		}
	}

	return result
}

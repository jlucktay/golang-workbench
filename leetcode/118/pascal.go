package pascal

func generate(numRows int) [][]int {
	result := make([][]int, numRows)

	result[0] = append(result[0], 1)

	for i := 1; i < numRows; i++ {
		result[i] = make([]int, 0)
		result[i] = append(result[i], 1)

		for j := 1; j < i; j++ {
			result[i] = append(result[i], result[i-1][j-1]+result[i-1][j])
		}

		result[i] = append(result[i], 1)
	}

	return result
}

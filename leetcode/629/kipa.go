package kipa

var v = [1001][1001]int{}

func KInversePairs(n int, k int) int {
	if n == 0 || k < 0 {
		return 0
	}

	if k == 0 {
		return 1
	}

	if v[n][k] == 0 {
		for i := 0; i <= min(n-1, k); i++ {
			v[n][k] = (v[n][k] + KInversePairs(n-1, k-i)) % 1000000007
		}
	}

	return v[n][k]
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

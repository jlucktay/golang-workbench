package kipa

var memo = [1001][1001]int{}

func init() {
	for a := range memo {
		for b := range memo[a] {
			memo[a][b] = -1
		}
	}
}

func KInversePairs(n int, k int) int {
	if n == 0 {
		return 0
	}

	if k == 0 {
		return 1
	}

	if memo[n][k] != -1 {
		return memo[n][k]
	}

	inv := 0

	for i := 0; i <= min(k, n-1); i++ {
		inv = (inv + KInversePairs(n-1, k-i)) % 1000000007
	}

	memo[n][k] = inv

	return inv
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

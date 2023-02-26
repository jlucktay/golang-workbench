package kipa

func KInversePairs(n, k int) int {
	dp := make([]int, k+1)
	const M = 1000000007

	for i := 1; i <= n; i++ {
		temp := make([]int, k+1)
		temp[0] = 1

		for j := 1; j <= k; j++ {
			x := 0
			if j-i >= 0 {
				x = dp[j-i]
			}

			val := (dp[j] + M - x) % M
			temp[j] = (temp[j-1] + val) % M
		}

		dp = temp
	}

	y := 0
	if k > 0 {
		y = dp[k-1]
	}

	return (dp[k] + M - y) % M
}

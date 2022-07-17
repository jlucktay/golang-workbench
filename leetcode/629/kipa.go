package kipa

func KInversePairs(n int, k int) int {
	dp := [1001][1001]int{}

	for i := 1; i <= n; i++ {
		for j := 0; j <= k; j++ {
			if j == 0 {
				dp[i][j] = 1
			} else {
				for p := 0; p <= min(j, i-1); p++ {
					dp[i][j] = (dp[i][j] + dp[i-1][j-p]) % 1000000007
				}
			}
		}
	}

	return dp[n][k]
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

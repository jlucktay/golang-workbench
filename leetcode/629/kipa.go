package kipa

func KInversePairs(n int, k int) int {
	dp := [1001][1001]int{}
	const M = 1000000007

	for i := 1; i <= n; i++ {
		for j := 0; j <= k; j++ {
			if j == 0 {
				dp[i][j] = 1
			} else {
				// (j - i) >= 0 ? dp[i - 1][j - i] : 0
				x := 0
				if j-i >= 0 {
					x = dp[i-1][j-i]
				}

				val := (dp[i-1][j] + M - x) % M
				dp[i][j] = (dp[i][j-1] + val) % M
			}
		}
	}

	// k > 0 ? dp[n][k - 1] : 0
	a := 0
	if k > 0 {
		a = dp[n][k-1]
	}

	return (dp[n][k] + M - a) % M
}

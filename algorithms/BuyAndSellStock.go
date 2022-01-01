package algorithms

// https://leetcode.com/problems/best-time-to-buy-and-sell-stock/
func MaxProfit(prices []float64) (float64, int) {
	// DP problem, see maxSubArray.
	// we always try find a low buy point, then sell high and record
	// this max. If we find another even low buy point, we change
	// the buy day to this day, then try to sell high again. But,
	// we will compare this new max with older one, and select the
	// bigger one.
	// loop the prices, at any index, we have a older max, which
	// starts as Integer.MIN_VALUE. We move the maxIndex if we find
	// lower price. Also record max if we found sell is maxed.
	// [7,1,5,3,6,4]
	var maxIndex int = 0
	var max float64 = -1000.0 // if buy - sell < 0.0, sell right away
	for i := 0; i < len(prices); i++ {
		tmp := prices[i] - prices[maxIndex]
		if tmp < 0.0 {
			// we found even lower price to buy
			maxIndex = i
		} else {
			if max < tmp {
				// we found even higher price to sell
				max = tmp
			}
		}
	}
	// if max < 0 {
	// 	max = 0
	// }
	return max, maxIndex
}

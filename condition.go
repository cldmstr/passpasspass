package passpasspass

type Condition func(hand *Hand) bool

func NewHCPIntervalCondition(minHCP, maxHCP int) Condition {
	return func(hand *Hand) bool {
		return minHCP <= hand.HCP() && hand.HCP() <= maxHCP
	}
}

func NewMinMaxDistributionCondition(minDist, maxDist Distribution) Condition {
	return func(hand *Hand) bool {
		for i, d := range hand.Distribution() {
			if minDist[i] > d || maxDist[i] < d {
				return false
			}
		}

		return true
	}
}

func NewBalancedDistributionCondition() Condition {
	return NewMinMaxDistributionCondition(
		NewDistribution(3, 3, 3, 3),
		NewDistribution(5, 5, 5, 5),
	)
}

package passpasspass

type Condition func(hand *Hand) bool

func NewMultiCondition(conditions ...Condition) Condition {
	return func(hand *Hand) bool {
		for _, c := range conditions {
			if c(hand) {
				return true
			}
		}

		return false
	}
}

func NewAllCondition(conditions ...Condition) Condition {
	return func(hand *Hand) bool {
		for _, c := range conditions {
			if !c(hand) {
				return false
			}
		}

		return true
	}
}

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
	return func(hand *Hand) bool {
		doubletons := 0
		for _, n := range hand.Distribution() {
			if n <= 1 || n >= 6 {
				return false
			}
			if n == 2 {
				doubletons++
			}
		}
		return doubletons <= 1
	}
}

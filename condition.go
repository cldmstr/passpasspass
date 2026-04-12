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

type SuitRange struct{ Min, Max int }
type DistributionRange [4]SuitRange // [clubs, diamonds, hearts, spades]

func NewMinMaxDistributionCondition(r DistributionRange) Condition {
	return func(hand *Hand) bool {
		for i, d := range hand.Distribution() {
			if r[i].Min > d || r[i].Max < d {
				return false
			}
		}

		return true
	}
}

func NewRuleOfTwentyCondition() Condition {
	return func(hand *Hand) bool {
		d := hand.Distribution()
		longest, second := 0, 0
		for _, n := range d {
			if n >= longest {
				longest, second = n, longest
			} else if n > second {
				second = n
			}
		}
		return hand.HCP()+longest+second >= 20
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

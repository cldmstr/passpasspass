package passpasspass

var Twooverone *ConventionSheet

func init() {
	Twooverone = NewConventionSheet()
	c1 := NewConvention(
		NewBid(1, Clubs),
		NewMinMaxDistributionCondition(NewDistribution(2, 1, 1, 1), NewDistribution(6, 5, 5, 5)),
		NewHCPIntervalCondition(11, 17),
	)
	d1 := NewConvention(
		NewBid(1, Diamonds),
		NewMinMaxDistributionCondition(NewDistribution(1, 4, 1, 1), NewDistribution(5, 6, 5, 5)),
		NewHCPIntervalCondition(11, 17),
	)
	h1 := NewConvention(
		NewBid(1, Hearts),
		NewMinMaxDistributionCondition(NewDistribution(1, 1, 5, 1), NewDistribution(5, 5, 6, 5)),
		NewHCPIntervalCondition(11, 17),
	)
	s1 := NewConvention(
		NewBid(1, Spades),
		NewMinMaxDistributionCondition(NewDistribution(1, 1, 1, 5), NewDistribution(5, 5, 5, 6)),
		NewHCPIntervalCondition(11, 17),
	)
	nt1 := NewConvention(
		NewBid(1, NoTrump),
		NewBalancedDistributionCondition(),
		NewHCPIntervalCondition(15, 17),
	)
	c2 := NewConvention(
		NewBid(2, Clubs),
		NewHCPIntervalCondition(20, 22),
	)
	d2 := NewConvention(
		NewBid(2, Diamonds),
		NewMultiCondition(
			NewAllCondition(NewHCPIntervalCondition(23, 24), NewBalancedDistributionCondition()),
			NewAllCondition(NewHCPIntervalCondition(6, 9), NewMinMaxDistributionCondition(NewDistribution(1, 1, 6, 1), NewDistribution(5, 5, 13, 5))),
			NewAllCondition(NewHCPIntervalCondition(6, 9), NewMinMaxDistributionCondition(NewDistribution(1, 1, 1, 6), NewDistribution(5, 5, 5, 13))),
			NewHCPIntervalCondition(27, 40),
		),
	)
	Twooverone.Add(d2)
	Twooverone.Add(c2)
	Twooverone.Add(nt1)
	Twooverone.Add(s1)
	Twooverone.Add(h1)
	Twooverone.Add(d1)
	Twooverone.Add(c1)
}

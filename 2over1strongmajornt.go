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
	Twooverone.sheet = append(Twooverone.sheet, c2)
	Twooverone.sheet = append(Twooverone.sheet, nt1)
	Twooverone.sheet = append(Twooverone.sheet, s1)
	Twooverone.sheet = append(Twooverone.sheet, h1)
	Twooverone.sheet = append(Twooverone.sheet, d1)
	Twooverone.sheet = append(Twooverone.sheet, c1)
}

package passpasspass

type Condition func(ctx *AuctionContext) bool

func NewMultiCondition(conditions ...Condition) Condition {
	return func(ctx *AuctionContext) bool {
		for _, c := range conditions {
			if c(ctx) {
				return true
			}
		}
		return false
	}
}

func NewAllCondition(conditions ...Condition) Condition {
	return func(ctx *AuctionContext) bool {
		for _, c := range conditions {
			if !c(ctx) {
				return false
			}
		}
		return true
	}
}

func NewHCPIntervalCondition(minHCP, maxHCP int) Condition {
	return func(ctx *AuctionContext) bool {
		return minHCP <= ctx.SelfHand.HCP() && ctx.SelfHand.HCP() <= maxHCP
	}
}

type SuitRange struct{ Min, Max int }
type DistributionRange [4]SuitRange // [clubs, diamonds, hearts, spades]

func NewMinMaxDistributionCondition(r DistributionRange) Condition {
	return func(ctx *AuctionContext) bool {
		for i, d := range ctx.SelfHand.Distribution() {
			if r[i].Min > d || r[i].Max < d {
				return false
			}
		}
		return true
	}
}

func NewRuleOfTwentyCondition() Condition {
	return func(ctx *AuctionContext) bool {
		d := ctx.SelfHand.Distribution()
		longest, second := 0, 0
		for _, n := range d {
			if n >= longest {
				longest, second = n, longest
			} else if n > second {
				second = n
			}
		}
		return ctx.SelfHand.HCP()+longest+second >= 20
	}
}

func NewBalancedDistributionCondition() Condition {
	return func(ctx *AuctionContext) bool {
		doubletons := 0
		for _, n := range ctx.SelfHand.Distribution() {
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

// NewNoPriorBidsCondition returns true when no real (non-pass) bid has been made yet.
func NewNoPriorBidsCondition() Condition {
	return func(ctx *AuctionContext) bool {
		for _, action := range ctx.History {
			if action.Bid != Pass {
				return false
			}
		}
		return true
	}
}

// NewPartnerBidCondition returns true when partner's last bid matches bidStr.
func NewPartnerBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return func(ctx *AuctionContext) bool {
		if target == nil {
			return false
		}
		return lastBidByDir(ctx.History, partnerOf(ctx.Position)).Value() == target.Value()
	}
}

// NewRHOBidCondition returns true when RHO's last bid matches bidStr.
func NewRHOBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return func(ctx *AuctionContext) bool {
		if target == nil {
			return false
		}
		return lastBidByDir(ctx.History, rhoOf(ctx.Position)).Value() == target.Value()
	}
}

// NewLHOBidCondition returns true when LHO's last bid matches bidStr.
func NewLHOBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return func(ctx *AuctionContext) bool {
		if target == nil {
			return false
		}
		return lastBidByDir(ctx.History, lhoOf(ctx.Position)).Value() == target.Value()
	}
}

// NewPartnershipHCPCondition returns true when the inferred combined HCP is within [min, max].
func NewPartnershipHCPCondition(min, max int) Condition {
	return func(ctx *AuctionContext) bool {
		return ctx.Partnership.MinCombinedHCP >= min && ctx.Partnership.MaxCombinedHCP <= max
	}
}

// NewMajorFitCondition returns true when the best major suit fit is within [min, max] cards.
func NewMajorFitCondition(min, max int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.BestMajorFit()
		return fit.Min >= min && fit.Max <= max
	}
}

// NewMinorFitCondition returns true when the best minor suit fit is within [min, max] cards.
func NewMinorFitCondition(min, max int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.BestMinorFit()
		return fit.Min >= min && fit.Max <= max
	}
}

// NewHeartFitCondition returns true when the combined heart holding is within [min, max] cards.
func NewHeartFitCondition(min, max int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.SuitFit[Hearts]
		return fit.Min >= min && fit.Max <= max
	}
}

// NewSpadeFitCondition returns true when the combined spade holding is within [min, max] cards.
func NewSpadeFitCondition(min, max int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.SuitFit[Spades]
		return fit.Min >= min && fit.Max <= max
	}
}

// NewLHOMinHCPCondition returns true when LHO's inferred minimum HCP meets the threshold.
func NewLHOMinHCPCondition(min int) Condition {
	return func(ctx *AuctionContext) bool { return ctx.LHO.MinHCP >= min }
}

// NewRHOMinHCPCondition returns true when RHO's inferred minimum HCP meets the threshold.
func NewRHOMinHCPCondition(min int) Condition {
	return func(ctx *AuctionContext) bool { return ctx.RHO.MinHCP >= min }
}

func partnerOf(d Direction) Direction { return (d + 2) % 4 }
func rhoOf(d Direction) Direction     { return (d + 3) % 4 }
func lhoOf(d Direction) Direction     { return (d + 1) % 4 }

func lastBidByDir(history []BidAction, d Direction) *Bid {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Direction == d {
			return history[i].Bid
		}
	}
	return Pass
}

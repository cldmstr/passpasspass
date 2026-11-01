package passpasspass

// Condition decides whether a convention applies in the given auction context.
type Condition func(ctx *AuctionContext) bool

// NewMultiCondition passes when at least one of conditions passes.
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

// NewAllCondition passes when every one of conditions passes.
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

// NewHCPIntervalCondition passes when the bidder's HCP lies within [minHCP, maxHCP].
func NewHCPIntervalCondition(minHCP, maxHCP int) Condition {
	return func(ctx *AuctionContext) bool {
		return minHCP <= ctx.SelfHand.HCP() && ctx.SelfHand.HCP() <= maxHCP
	}
}

// SuitRange is an inclusive range of card counts.
type SuitRange struct{ Min, Max int }

// DistributionRange holds a SuitRange per suit.
type DistributionRange [4]SuitRange // [clubs, diamonds, hearts, spades]

// NewMinMaxDistributionCondition passes when every suit length lies within its range in r.
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

// NewRuleOfTwentyCondition passes when HCP plus the lengths of the two longest suits is at least 20.
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

// NewBalancedDistributionCondition passes for a hand with no void or singleton, at most one
// doubleton and no suit of six or more cards.
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
		_, opened := ctx.History.Opening()
		return !opened
	}
}

// NewPartnerBidCondition returns true when partner's last bid matches bidStr.
func NewPartnerBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return func(ctx *AuctionContext) bool {
		if target == nil {
			return false
		}
		return ctx.History.LastBidBy(partnerOf(ctx.Position)).Value() == target.Value()
	}
}

// NewRHOBidCondition returns true when RHO's last bid matches bidStr.
func NewRHOBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return func(ctx *AuctionContext) bool {
		if target == nil {
			return false
		}
		return ctx.History.LastBidBy(rhoOf(ctx.Position)).Value() == target.Value()
	}
}

// NewLHOBidCondition returns true when LHO's last bid matches bidStr.
func NewLHOBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return func(ctx *AuctionContext) bool {
		if target == nil {
			return false
		}
		return ctx.History.LastBidBy(lhoOf(ctx.Position)).Value() == target.Value()
	}
}

// NewPartnershipHCPCondition returns true when the inferred combined HCP is within [lo, hi].
func NewPartnershipHCPCondition(lo, hi int) Condition {
	return func(ctx *AuctionContext) bool {
		return ctx.Partnership.MinCombinedHCP >= lo && ctx.Partnership.MaxCombinedHCP <= hi
	}
}

// NewMajorFitCondition returns true when the best major suit fit is within [lo, hi] cards.
func NewMajorFitCondition(lo, hi int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.BestMajorFit()
		return fit.Min >= lo && fit.Max <= hi
	}
}

// NewMinorFitCondition returns true when the best minor suit fit is within [lo, hi] cards.
func NewMinorFitCondition(lo, hi int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.BestMinorFit()
		return fit.Min >= lo && fit.Max <= hi
	}
}

// NewHeartFitCondition returns true when the combined heart holding is within [lo, hi] cards.
func NewHeartFitCondition(lo, hi int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.SuitFit[Hearts]
		return fit.Min >= lo && fit.Max <= hi
	}
}

// NewSpadeFitCondition returns true when the combined spade holding is within [lo, hi] cards.
func NewSpadeFitCondition(lo, hi int) Condition {
	return func(ctx *AuctionContext) bool {
		fit := ctx.Partnership.SuitFit[Spades]
		return fit.Min >= lo && fit.Max <= hi
	}
}

// NewLHOMinHCPCondition returns true when LHO's inferred minimum HCP meets the threshold.
func NewLHOMinHCPCondition(threshold int) Condition {
	return func(ctx *AuctionContext) bool { return ctx.LHO.MinHCP >= threshold }
}

// NewRHOMinHCPCondition returns true when RHO's inferred minimum HCP meets the threshold.
func NewRHOMinHCPCondition(threshold int) Condition {
	return func(ctx *AuctionContext) bool { return ctx.RHO.MinHCP >= threshold }
}

func partnerOf(d Direction) Direction { return (d + 2) % 4 }
func rhoOf(d Direction) Direction     { return (d + 3) % 4 }
func lhoOf(d Direction) Direction     { return (d + 1) % 4 }

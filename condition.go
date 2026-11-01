package passpasspass

// Condition decides whether a convention applies in the given auction context.
// Each condition key in a convention file has its own Condition type.
type Condition interface {
	Eval(ctx *AuctionContext) bool
}

// anyCondition passes when at least one child passes.
type anyCondition struct{ children []Condition }

// NewMultiCondition passes when at least one of conditions passes.
func NewMultiCondition(conditions ...Condition) Condition {
	return anyCondition{children: conditions}
}

func (c anyCondition) Eval(ctx *AuctionContext) bool {
	for _, child := range c.children {
		if child.Eval(ctx) {
			return true
		}
	}
	return false
}

// allCondition passes when every child passes.
type allCondition struct{ children []Condition }

// NewAllCondition passes when every one of conditions passes.
func NewAllCondition(conditions ...Condition) Condition {
	return allCondition{children: conditions}
}

func (c allCondition) Eval(ctx *AuctionContext) bool {
	for _, child := range c.children {
		if !child.Eval(ctx) {
			return false
		}
	}
	return true
}

// hcpCondition passes when the bidder's HCP lies within [lo, hi].
type hcpCondition struct{ lo, hi int }

// NewHCPIntervalCondition passes when the bidder's HCP lies within [minHCP, maxHCP].
func NewHCPIntervalCondition(minHCP, maxHCP int) Condition {
	return hcpCondition{lo: minHCP, hi: maxHCP}
}

func (c hcpCondition) Eval(ctx *AuctionContext) bool {
	hcp := ctx.SelfHand.HCP()
	return c.lo <= hcp && hcp <= c.hi
}

// SuitRange is an inclusive range of card counts.
type SuitRange struct{ Min, Max int }

// DistributionRange holds a SuitRange per suit.
type DistributionRange [4]SuitRange // [clubs, diamonds, hearts, spades]

// suitLengthCondition passes when every suit length lies within its range.
type suitLengthCondition struct{ r DistributionRange }

// NewMinMaxDistributionCondition passes when every suit length lies within its range in r.
func NewMinMaxDistributionCondition(r DistributionRange) Condition {
	return suitLengthCondition{r: r}
}

func (c suitLengthCondition) Eval(ctx *AuctionContext) bool {
	for i, d := range ctx.SelfHand.Distribution() {
		if c.r[i].Min > d || c.r[i].Max < d {
			return false
		}
	}
	return true
}

// ruleOfTwentyCondition passes when HCP plus the two longest suits is at least 20.
type ruleOfTwentyCondition struct{}

// NewRuleOfTwentyCondition passes when HCP plus the lengths of the two longest suits is at least 20.
func NewRuleOfTwentyCondition() Condition {
	return ruleOfTwentyCondition{}
}

func (ruleOfTwentyCondition) Eval(ctx *AuctionContext) bool {
	longest, second := 0, 0
	for _, n := range ctx.SelfHand.Distribution() {
		if n >= longest {
			longest, second = n, longest
		} else if n > second {
			second = n
		}
	}
	return ctx.SelfHand.HCP()+longest+second >= 20
}

// balancedCondition passes for a balanced hand.
type balancedCondition struct{}

// NewBalancedDistributionCondition passes for a hand with no void or singleton, at most one
// doubleton and no suit of six or more cards.
func NewBalancedDistributionCondition() Condition {
	return balancedCondition{}
}

func (balancedCondition) Eval(ctx *AuctionContext) bool {
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

// noPriorBidsCondition passes while nobody has opened.
type noPriorBidsCondition struct{}

// NewNoPriorBidsCondition returns true when no real (non-pass) bid has been made yet.
func NewNoPriorBidsCondition() Condition {
	return noPriorBidsCondition{}
}

func (noPriorBidsCondition) Eval(ctx *AuctionContext) bool {
	_, opened := ctx.History.Opening()
	return !opened
}

// seatBidCondition passes when the last call by the seat that seat(bidder) returns equals bid.
// A nil bid never matches.
type seatBidCondition struct {
	seat func(Direction) Direction
	bid  *Bid
}

func (c seatBidCondition) Eval(ctx *AuctionContext) bool {
	if c.bid == nil {
		return false
	}
	return ctx.History.LastBidBy(c.seat(ctx.Position)).Value() == c.bid.Value()
}

// NewPartnerBidCondition returns true when partner's last bid matches bidStr.
func NewPartnerBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return seatBidCondition{seat: partnerOf, bid: target}
}

// NewRHOBidCondition returns true when RHO's last bid matches bidStr.
func NewRHOBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return seatBidCondition{seat: rhoOf, bid: target}
}

// NewLHOBidCondition returns true when LHO's last bid matches bidStr.
func NewLHOBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return seatBidCondition{seat: lhoOf, bid: target}
}

// partnershipHCPCondition passes when the combined HCP range lies within [lo, hi].
type partnershipHCPCondition struct{ lo, hi int }

// NewPartnershipHCPCondition returns true when the inferred combined HCP is within [lo, hi].
func NewPartnershipHCPCondition(lo, hi int) Condition {
	return partnershipHCPCondition{lo: lo, hi: hi}
}

func (c partnershipHCPCondition) Eval(ctx *AuctionContext) bool {
	return ctx.Partnership.MinCombinedHCP >= c.lo && ctx.Partnership.MaxCombinedHCP <= c.hi
}

// suitFitCondition passes when the combined holding in suit lies within [lo, hi] cards.
type suitFitCondition struct {
	suit   Suit
	lo, hi int
}

// NewHeartFitCondition returns true when the combined heart holding is within [lo, hi] cards.
func NewHeartFitCondition(lo, hi int) Condition {
	return suitFitCondition{suit: Hearts, lo: lo, hi: hi}
}

// NewSpadeFitCondition returns true when the combined spade holding is within [lo, hi] cards.
func NewSpadeFitCondition(lo, hi int) Condition {
	return suitFitCondition{suit: Spades, lo: lo, hi: hi}
}

func (c suitFitCondition) Eval(ctx *AuctionContext) bool {
	return within(ctx.Partnership.SuitFit[c.suit], c.lo, c.hi)
}

// majorFitCondition passes when the better major fit lies within [lo, hi] cards.
type majorFitCondition struct{ lo, hi int }

// NewMajorFitCondition returns true when the best major suit fit is within [lo, hi] cards.
func NewMajorFitCondition(lo, hi int) Condition {
	return majorFitCondition{lo: lo, hi: hi}
}

func (c majorFitCondition) Eval(ctx *AuctionContext) bool {
	return within(ctx.Partnership.BestMajorFit(), c.lo, c.hi)
}

// minorFitCondition passes when the better minor fit lies within [lo, hi] cards.
type minorFitCondition struct{ lo, hi int }

// NewMinorFitCondition returns true when the best minor suit fit is within [lo, hi] cards.
func NewMinorFitCondition(lo, hi int) Condition {
	return minorFitCondition{lo: lo, hi: hi}
}

func (c minorFitCondition) Eval(ctx *AuctionContext) bool {
	return within(ctx.Partnership.BestMinorFit(), c.lo, c.hi)
}

// opponentMinHCPCondition passes when an opponent's inferred minimum HCP is at least threshold.
type opponentMinHCPCondition struct {
	left      bool // LHO when true, RHO otherwise
	threshold int
}

// NewLHOMinHCPCondition returns true when LHO's inferred minimum HCP meets the threshold.
func NewLHOMinHCPCondition(threshold int) Condition {
	return opponentMinHCPCondition{left: true, threshold: threshold}
}

// NewRHOMinHCPCondition returns true when RHO's inferred minimum HCP meets the threshold.
func NewRHOMinHCPCondition(threshold int) Condition {
	return opponentMinHCPCondition{left: false, threshold: threshold}
}

func (c opponentMinHCPCondition) Eval(ctx *AuctionContext) bool {
	opponent := ctx.RHO
	if c.left {
		opponent = ctx.LHO
	}
	return opponent.MinHCP >= c.threshold
}

// within reports whether r lies entirely inside [lo, hi].
func within(r SuitRange, lo, hi int) bool {
	return r.Min >= lo && r.Max <= hi
}

func partnerOf(d Direction) Direction { return (d + 2) % 4 }
func rhoOf(d Direction) Direction     { return (d + 3) % 4 }
func lhoOf(d Direction) Direction     { return (d + 1) % 4 }

package passpasspass

// Condition decides whether a convention applies in the given auction context, and what a call
// made under it shows. Each condition key in a convention file has its own Condition type.
type Condition interface {
	// Eval reports whether the condition holds for the bidder in ctx.
	Eval(ctx *AuctionContext) bool

	// Shown returns the model of the hands for which the condition comes out as holds, for a call
	// made in ctx. Conditions that read the auction rather than the hand give every hand or none;
	// conditions a range cannot express exactly give a wider model, never a narrower one.
	Shown(ctx *CallContext, holds bool) HandModel

	// maxCases bounds the number of alternatives Shown returns for holds, in any auction.
	maxCases(holds bool) int
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

// Shown keeps each child's alternatives side by side; failing means every child fails.
func (c anyCondition) Shown(ctx *CallContext, holds bool) HandModel {
	if holds {
		return unionShown(ctx, true, c.children)
	}
	return intersectShown(ctx, false, c.children)
}

func (c anyCondition) maxCases(holds bool) int {
	if holds {
		return sumCases(true, c.children)
	}
	return productCases(false, c.children)
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

// Shown intersects the children; failing means at least one child fails.
func (c allCondition) Shown(ctx *CallContext, holds bool) HandModel {
	if holds {
		return intersectShown(ctx, true, c.children)
	}
	return unionShown(ctx, false, c.children)
}

func (c allCondition) maxCases(holds bool) int {
	if holds {
		return productCases(true, c.children)
	}
	return sumCases(false, c.children)
}

// intersectShown returns the hands for which every child comes out as holds.
func intersectShown(ctx *CallContext, holds bool, children []Condition) HandModel {
	m := DefaultHandModel()
	for _, child := range children {
		m = m.Intersect(child.Shown(ctx, holds))
		if m.IsEmpty() {
			break
		}
	}
	return m
}

// unionShown returns the hands for which at least one child comes out as holds.
func unionShown(ctx *CallContext, holds bool, children []Condition) HandModel {
	models := make([]HandModel, len(children))
	for i, child := range children {
		models[i] = child.Shown(ctx, holds)
	}
	return unionModels(models...)
}

// sumCases and productCases bound the alternatives of unionShown and intersectShown. Both stop
// counting just above maxConventionCases, so that deep nesting cannot overflow.
func sumCases(holds bool, children []Condition) int {
	n := 0
	for _, child := range children {
		n = min(n+child.maxCases(holds), maxConventionCases+1)
	}
	return n
}

func productCases(holds bool, children []Condition) int {
	n := 1
	for _, child := range children {
		n = min(n*child.maxCases(holds), maxConventionCases+1)
	}
	return n
}

// notCondition passes when its child fails.
type notCondition struct{ child Condition }

// NewNotCondition passes when condition fails.
func NewNotCondition(condition Condition) Condition {
	return notCondition{child: condition}
}

func (c notCondition) Eval(ctx *AuctionContext) bool {
	return !c.child.Eval(ctx)
}

// Shown asks the child for the opposite outcome.
func (c notCondition) Shown(ctx *CallContext, holds bool) HandModel {
	return c.child.Shown(ctx, !holds)
}

func (c notCondition) maxCases(holds bool) int {
	return c.child.maxCases(!holds)
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

// Shown is the HCP range, or one alternative below it and one above it.
func (c hcpCondition) Shown(_ *CallContext, holds bool) HandModel {
	if holds {
		return modelOf(hcpRange(c.lo, c.hi))
	}
	var m HandModel
	if c.lo > 0 {
		m.Cases = append(m.Cases, hcpRange(0, c.lo-1))
	}
	if c.hi < MaxHandHCP {
		m.Cases = append(m.Cases, hcpRange(c.hi+1, MaxHandHCP))
	}
	return m.Tighten()
}

func (c hcpCondition) maxCases(holds bool) int {
	return len(c.Shown(nil, holds).Cases)
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

// Shown is the suit ranges, or one alternative for each bound that a suit can break: shorter
// than its minimum or longer than its maximum.
func (c suitLengthCondition) Shown(_ *CallContext, holds bool) HandModel {
	if holds {
		r := DefaultHandRange()
		r.SuitLength = c.r
		return modelOf(r)
	}
	var m HandModel
	for s, sr := range c.r {
		if sr.Min > 0 {
			m.Cases = append(m.Cases, suitRange(s, 0, sr.Min-1))
		}
		if sr.Max < HandSize {
			m.Cases = append(m.Cases, suitRange(s, sr.Max+1, HandSize))
		}
	}
	return m.Tighten()
}

func (c suitLengthCondition) maxCases(holds bool) int {
	return len(c.Shown(nil, holds).Cases)
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

// Shown is every hand: a range cannot express the Rule of 20.
func (ruleOfTwentyCondition) Shown(*CallContext, bool) HandModel {
	return DefaultHandModel()
}

func (ruleOfTwentyCondition) maxCases(bool) int { return 1 }

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

// Shown is 2–5 cards in every suit, which also allows some unbalanced hands such as 5-4-2-2.
// The complement of that approximation is not known, so failing shows every hand.
func (balancedCondition) Shown(_ *CallContext, holds bool) HandModel {
	if !holds {
		return DefaultHandModel()
	}
	r := DefaultHandRange()
	for s := range r.SuitLength {
		r.SuitLength[s] = SuitRange{Min: 2, Max: 5}
	}
	return modelOf(r)
}

func (balancedCondition) maxCases(bool) int { return 1 }

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

func (c noPriorBidsCondition) Shown(ctx *CallContext, holds bool) HandModel {
	return guardModel(c.Eval(ctx.auctionContext()) == holds)
}

func (noPriorBidsCondition) maxCases(bool) int { return 1 }

// auctionPatternCondition passes when its pattern matches the whole auction so far.
type auctionPatternCondition struct{ pattern AuctionPattern }

// NewAuctionPatternCondition passes when pattern matches every call made so far, from the
// dealer's first call, including passes before the opening.
func NewAuctionPatternCondition(pattern AuctionPattern) Condition {
	return auctionPatternCondition{pattern: pattern}
}

func (c auctionPatternCondition) Eval(ctx *AuctionContext) bool {
	return c.pattern.Matches(ctx.History)
}

func (c auctionPatternCondition) Shown(ctx *CallContext, holds bool) HandModel {
	return guardModel(c.Eval(ctx.auctionContext()) == holds)
}

func (auctionPatternCondition) maxCases(bool) int { return 1 }

// partnerBidCondition passes when partner's last call equals bid. A nil bid never matches.
type partnerBidCondition struct{ bid *Bid }

// NewPartnerBidCondition returns true when partner's last bid matches bidStr.
func NewPartnerBidCondition(bidStr string) Condition {
	target, _ := NewBidFromString(bidStr)
	return partnerBidCondition{bid: target}
}

func (c partnerBidCondition) Eval(ctx *AuctionContext) bool {
	if c.bid == nil {
		return false
	}
	return ctx.History.LastBidBy(partnerOf(ctx.Position)).Value() == c.bid.Value()
}

func (c partnerBidCondition) Shown(ctx *CallContext, holds bool) HandModel {
	return guardModel(c.Eval(ctx.auctionContext()) == holds)
}

func (partnerBidCondition) maxCases(bool) int { return 1 }

// partnershipHCPCondition passes when the combined HCP range lies within [lo, hi].
type partnershipHCPCondition struct{ lo, hi int }

// NewPartnershipHCPCondition returns true when the inferred combined HCP is within [lo, hi].
func NewPartnershipHCPCondition(lo, hi int) Condition {
	return partnershipHCPCondition{lo: lo, hi: hi}
}

func (c partnershipHCPCondition) Eval(ctx *AuctionContext) bool {
	return ctx.Partnership.MinCombinedHCP >= c.lo && ctx.Partnership.MaxCombinedHCP <= c.hi
}

// Shown turns the combined range into a range for the caller's own HCP, given what partner had
// shown: own HCP plus partner's minimum must reach lo, and own HCP plus partner's maximum must
// stay within hi. Failing shows every hand, as the combined range is itself an estimate.
func (c partnershipHCPCondition) Shown(ctx *CallContext, holds bool) HandModel {
	partner, ok := ctx.Shown[partnerOf(ctx.Caller)].Bounds()
	if !holds || !ok {
		return DefaultHandModel()
	}
	return modelOf(hcpRange(max(0, c.lo-partner.MinHCP), c.hi-partner.MaxHCP))
}

func (partnershipHCPCondition) maxCases(bool) int { return 1 }

func partnerOf(d Direction) Direction { return (d + 2) % 4 }
func rhoOf(d Direction) Direction     { return (d + 3) % 4 }
func lhoOf(d Direction) Direction     { return (d + 1) % 4 }

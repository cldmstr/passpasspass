package passpasspass

import (
	"fmt"
	"strings"
)

// BidAction is one step in the auction.
type BidAction struct {
	Direction Direction
	Bid       *Bid
}

// History is the ordered record of calls in an auction, starting with the dealer's call.
// The auction as it stood when call i was made is h[:i].
type History []BidAction

// BidsBy returns every call made by d, Passes included, in order.
func (h History) BidsBy(d Direction) []*Bid {
	var bids []*Bid
	for _, action := range h {
		if action.Direction == d {
			bids = append(bids, action.Bid)
		}
	}
	return bids
}

// LastBidBy returns the most recent call made by d, or Pass if d has not called yet.
func (h History) LastBidBy(d Direction) *Bid {
	for i := len(h) - 1; i >= 0; i-- {
		if h[i].Direction == d {
			return h[i].Bid
		}
	}
	return Pass
}

// Opening returns the first non-pass call and true, or false when nobody has bid yet.
func (h History) Opening() (BidAction, bool) {
	for _, action := range h {
		if action.Bid != Pass {
			return action, true
		}
	}
	return BidAction{}, false
}

// Contract returns the most recent non-pass call and true, or false when nobody has bid yet.
func (h History) Contract() (BidAction, bool) {
	for i := len(h) - 1; i >= 0; i-- {
		if h[i].Bid != Pass {
			return h[i], true
		}
	}
	return BidAction{}, false
}

// HandRange holds inferred HCP and suit-length bounds for one player.
type HandRange struct {
	MinHCP, MaxHCP int
	SuitLength     [4]SuitRange // indexed clubs(0)…spades(3)
}

// DefaultHandRange returns an unconstrained HandRange.
func DefaultHandRange() HandRange {
	return HandRange{
		MinHCP: 0, MaxHCP: MaxHandHCP,
		SuitLength: [4]SuitRange{
			{0, HandSize}, {0, HandSize}, {0, HandSize}, {0, HandSize},
		},
	}
}

// PartnershipEval captures the combined picture of self + partner.
type PartnershipEval struct {
	MinCombinedHCP int
	MaxCombinedHCP int
	SuitFit        [4]SuitRange
}

// NewPartnershipEval computes the combined picture from own hand and partner's inferred range.
func NewPartnershipEval(self *Hand, partner HandRange) PartnershipEval {
	selfDist := self.Distribution()
	var fit [4]SuitRange
	for i := range 4 {
		fit[i] = SuitRange{
			Min: selfDist[i] + partner.SuitLength[i].Min,
			Max: selfDist[i] + partner.SuitLength[i].Max,
		}
	}
	return PartnershipEval{
		MinCombinedHCP: self.HCP() + partner.MinHCP,
		MaxCombinedHCP: self.HCP() + partner.MaxHCP,
		SuitFit:        fit,
	}
}

// AuctionContext is passed to every Condition during a live auction.
type AuctionContext struct {
	SelfHand    *Hand
	Position    Direction
	History     History
	Partner     HandRange
	LHO         HandRange
	RHO         HandRange
	Partnership PartnershipEval
}

// Auction records a complete bridge auction: the deal, all bids, and the final contract.
type Auction struct {
	Dealer   Direction   // who opened the bidding
	Hands    [4]*Hand    // hands indexed by Direction (North=0…West=3)
	Bids     []BidAction // every bid in order, populated after Run
	Contract *Bid        // nil = passed out, populated after Run
	Declarer Direction   // who made the final non-pass bid, populated after Run
	cs       *ConventionSheet
}

// NewAuction creates an Auction with the given dealer direction, convention sheet, and hands.
func NewAuction(dealer Direction, cs *ConventionSheet, hands []*Hand) *Auction {
	a := &Auction{
		Dealer: dealer,
		Bids:   make([]BidAction, 0, 40),
		cs:     cs,
	}
	copy(a.Hands[:], hands)
	return a
}

// Run executes the full bidding sequence and populates Bids, Contract, and Declarer.
func (a *Auction) Run() {
	var handRanges [4]HandRange
	for i := range 4 {
		handRanges[i] = DefaultHandRange()
	}

	consecutivePasses := 0
	current := a.Dealer

	for {
		pos := current
		partner := partnerOf(pos)
		h := a.Hands[int(pos)]

		ctx := &AuctionContext{
			SelfHand:    h,
			Position:    pos,
			History:     a.Bids,
			Partner:     handRanges[partner],
			LHO:         handRanges[lhoOf(pos)],
			RHO:         handRanges[rhoOf(pos)],
			Partnership: NewPartnershipEval(h, handRanges[partner]),
		}

		bid := a.cs.Bid(ctx)

		// Bid legality: must be strictly above the current contract
		if bid != Pass && a.Contract != nil && bid.Value() <= a.Contract.Value() {
			bid = Pass
		}

		a.Bids = append(a.Bids, BidAction{Direction: pos, Bid: bid})

		if bid == Pass {
			consecutivePasses++
		} else {
			consecutivePasses = 0
			a.Contract = bid
			a.Declarer = pos
		}

		inferAndUpdate(&handRanges, a.cs, pos, bid)
		current = (current + 1) % 4

		if (a.Contract == nil && consecutivePasses == 4) ||
			(a.Contract != nil && consecutivePasses == 3) {
			break
		}
	}
}

// String shows all four hands followed by the bidding sequence and contract.
func (a *Auction) String() string {
	var sb strings.Builder

	for i, dir := range allDirections {
		h := a.Hands[i]
		fmt.Fprintf(&sb, "%s (%d HCP %s):\n%s\n", dir, h.HCP(), h.Distribution(), h)
	}

	fmt.Fprintf(&sb, "Dealer: %s\n", allDirections[a.Dealer])
	sb.WriteString("North\tEast\tSouth\tWest\n")

	col := int(a.Dealer)
	for range col {
		sb.WriteString("\t")
	}
	for _, ba := range a.Bids {
		sb.WriteString(ba.Bid.String())
		col++
		if col == 4 {
			sb.WriteString("\n")
			col = 0
		} else {
			sb.WriteString("\t")
		}
	}
	if col != 0 {
		sb.WriteString("\n")
	}

	if a.Contract == nil {
		sb.WriteString("\nPassed out\n")
	} else {
		fmt.Fprintf(&sb, "\nContract: %s by %s\n", a.Contract, allDirections[a.Declarer])
	}

	return sb.String()
}

func inferAndUpdate(handRanges *[4]HandRange, cs *ConventionSheet, pos Direction, bid *Bid) {
	if bid == Pass {
		return
	}
	for _, conv := range cs.sheet {
		if conv.Bid != nil && conv.Bid.Value() == bid.Value() {
			handRanges[pos] = intersectRanges(handRanges[pos], extractBounds(conv.conditions...))
			return
		}
	}
}

func extractBounds(conditions ...Condition) HandRange {
	r := DefaultHandRange()
	for _, c := range conditions {
		applyBounds(c, &r)
	}
	return r
}

// applyBounds narrows r by the HCP and suit-length ranges that c requires. An `any` narrows r
// by the union of its children's ranges; conditions without a range leave r unchanged.
func applyBounds(c Condition, r *HandRange) {
	switch c := c.(type) {
	case hcpCondition:
		b := DefaultHandRange()
		b.MinHCP, b.MaxHCP = c.lo, c.hi
		*r = intersectRanges(*r, b)

	case suitLengthCondition:
		b := DefaultHandRange()
		b.SuitLength = c.r
		*r = intersectRanges(*r, b)

	case allCondition:
		for _, child := range c.children {
			applyBounds(child, r)
		}

	case anyCondition:
		if len(c.children) == 0 {
			return
		}
		u := extractBounds(c.children[0])
		for _, child := range c.children[1:] {
			u = unionRanges(u, extractBounds(child))
		}
		*r = intersectRanges(*r, u)
	}
}

func intersectRanges(a, b HandRange) HandRange {
	r := a
	if b.MinHCP > r.MinHCP {
		r.MinHCP = b.MinHCP
	}
	if b.MaxHCP < r.MaxHCP {
		r.MaxHCP = b.MaxHCP
	}
	for i := range 4 {
		if b.SuitLength[i].Min > r.SuitLength[i].Min {
			r.SuitLength[i].Min = b.SuitLength[i].Min
		}
		if b.SuitLength[i].Max < r.SuitLength[i].Max {
			r.SuitLength[i].Max = b.SuitLength[i].Max
		}
	}
	return r
}

// unionRanges returns the smallest range that holds both a and b.
func unionRanges(a, b HandRange) HandRange {
	r := a
	r.MinHCP = min(r.MinHCP, b.MinHCP)
	r.MaxHCP = max(r.MaxHCP, b.MaxHCP)
	for i := range 4 {
		r.SuitLength[i].Min = min(r.SuitLength[i].Min, b.SuitLength[i].Min)
		r.SuitLength[i].Max = max(r.SuitLength[i].Max, b.SuitLength[i].Max)
	}
	return r
}

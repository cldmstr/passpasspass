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

// BestMajorFit returns the combined card range for the better major suit fit.
func (p PartnershipEval) BestMajorFit() SuitRange {
	h, s := p.SuitFit[Hearts], p.SuitFit[Spades]
	if s.Max >= h.Max {
		return s
	}
	return h
}

// BestMinorFit returns the combined card range for the better minor suit fit.
func (p PartnershipEval) BestMinorFit() SuitRange {
	c, d := p.SuitFit[Clubs], p.SuitFit[Diamonds]
	if d.Max >= c.Max {
		return d
	}
	return c
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
	for i, h := range hands {
		a.Hands[i] = h
	}
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
		sb.WriteString(fmt.Sprintf("%s (%d HCP %s):\n%s\n", dir, h.HCP(), h.Distribution(), h))
	}

	sb.WriteString(fmt.Sprintf("Dealer: %s\n", allDirections[a.Dealer]))
	sb.WriteString("North\tEast\tSouth\tWest\n")

	col := int(a.Dealer)
	for i := 0; i < col; i++ {
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
		sb.WriteString(fmt.Sprintf("\nContract: %s by %s\n", a.Contract, allDirections[a.Declarer]))
	}

	return sb.String()
}

func inferAndUpdate(handRanges *[4]HandRange, cs *ConventionSheet, pos Direction, bid *Bid) {
	if bid == Pass {
		return
	}
	for _, conv := range cs.sheet {
		if conv.Bid != nil && conv.Bid.Value() == bid.Value() {
			handRanges[pos] = intersectRanges(handRanges[pos], extractBounds(conv.node))
			return
		}
	}
}

func extractBounds(node conditionNode) HandRange {
	r := DefaultHandRange()
	applyBounds(node, &r)
	return r
}

func applyBounds(node conditionNode, r *HandRange) {
	switch {
	case node.HCP != nil && len(node.HCP) == 2:
		if node.HCP[0] > r.MinHCP {
			r.MinHCP = node.HCP[0]
		}
		if node.HCP[1] < r.MaxHCP {
			r.MaxHCP = node.HCP[1]
		}

	case node.SuitLength != nil:
		for suitStr, rng := range node.SuitLength {
			if len(rng) != 2 {
				continue
			}
			s, err := NewSuitFromString(suitStr)
			if err != nil || int(s) >= 4 {
				continue
			}
			if rng[0] > r.SuitLength[int(s)].Min {
				r.SuitLength[int(s)].Min = rng[0]
			}
			if rng[1] < r.SuitLength[int(s)].Max {
				r.SuitLength[int(s)].Max = rng[1]
			}
		}

	case node.All != nil:
		for _, child := range node.All {
			applyBounds(child, r)
		}

	case node.Any != nil:
		if len(node.Any) == 0 {
			return
		}
		unionMin, unionMax := MaxHandHCP, 0
		unionSuitMin := [4]int{HandSize, HandSize, HandSize, HandSize}
		unionSuitMax := [4]int{0, 0, 0, 0}
		for _, child := range node.Any {
			cr := DefaultHandRange()
			applyBounds(child, &cr)
			if cr.MinHCP < unionMin {
				unionMin = cr.MinHCP
			}
			if cr.MaxHCP > unionMax {
				unionMax = cr.MaxHCP
			}
			for i := range 4 {
				if cr.SuitLength[i].Min < unionSuitMin[i] {
					unionSuitMin[i] = cr.SuitLength[i].Min
				}
				if cr.SuitLength[i].Max > unionSuitMax[i] {
					unionSuitMax[i] = cr.SuitLength[i].Max
				}
			}
		}
		if unionMin > r.MinHCP {
			r.MinHCP = unionMin
		}
		if unionMax < r.MaxHCP {
			r.MaxHCP = unionMax
		}
		for i := range 4 {
			if unionSuitMin[i] > r.SuitLength[i].Min {
				r.SuitLength[i].Min = unionSuitMin[i]
			}
			if unionSuitMax[i] < r.SuitLength[i].Max {
				r.SuitLength[i].Max = unionSuitMax[i]
			}
		}
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

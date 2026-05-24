package passpasspass

import (
	"fmt"
	"strings"
)

const (
	North Direction = iota
	East
	South
	West
)

type Direction int

var allDirections = []string{"North", "East", "South", "West"}

type Table struct {
	cs         *ConventionSheet
	hands      []*Hand
	handRanges [4]HandRange
}

func NewTable(cs *ConventionSheet, hands []*Hand) *Table {
	var ranges [4]HandRange
	for i := range 4 {
		ranges[i] = DefaultHandRange()
	}
	return &Table{cs: cs, hands: hands, handRanges: ranges}
}

func (t Table) String() string {
	directionFormat := "%s (%d : %s):\nBid: %s\n%s\n"
	sb := strings.Builder{}
	for i, d := range allDirections {
		pos := Direction(i)
		h := t.hands[i]
		ctx := &AuctionContext{
			SelfHand:    h,
			Position:    pos,
			Partner:     t.handRanges[partnerOf(pos)],
			LHO:         t.handRanges[lhoOf(pos)],
			RHO:         t.handRanges[rhoOf(pos)],
			Partnership: NewPartnershipEval(h, t.handRanges[partnerOf(pos)]),
		}
		sb.WriteString(fmt.Sprintf(directionFormat, d, h.HCP(), h.Distribution(), t.cs.Bid(ctx), h))
	}
	return sb.String()
}

// Auction records the full sequence of bids and the final contract.
type Auction struct {
	Bids     []BidAction
	Contract *Bid
	Declarer Direction
}

// RunAuction runs a sequential auction starting from dealer.
func (t *Table) RunAuction(dealer Direction) *Auction {
	auction := &Auction{Bids: make([]BidAction, 0, 40)}
	consecutivePasses := 0
	current := dealer

	for {
		pos := current
		partner := partnerOf(pos)
		h := t.hands[int(pos)]

		ctx := &AuctionContext{
			SelfHand:    h,
			Position:    pos,
			History:     auction.Bids,
			Partner:     t.handRanges[partner],
			LHO:         t.handRanges[lhoOf(pos)],
			RHO:         t.handRanges[rhoOf(pos)],
			Partnership: NewPartnershipEval(h, t.handRanges[partner]),
		}

		bid := t.cs.Bid(ctx)

		// Bid legality: must be strictly above the current contract
		if bid != Pass && auction.Contract != nil && bid.Value() <= auction.Contract.Value() {
			bid = Pass
		}

		auction.Bids = append(auction.Bids, BidAction{Direction: pos, Bid: bid})

		if bid == Pass {
			consecutivePasses++
		} else {
			consecutivePasses = 0
			auction.Contract = bid
			auction.Declarer = pos
		}

		t.inferAndUpdate(pos, bid)
		current = (current + 1) % 4

		if (auction.Contract == nil && consecutivePasses == 4) ||
			(auction.Contract != nil && consecutivePasses == 3) {
			break
		}
	}

	return auction
}

func (t *Table) inferAndUpdate(pos Direction, bid *Bid) {
	if bid == Pass {
		return
	}
	for _, conv := range t.cs.sheet {
		if conv.Bid != nil && conv.Bid.Value() == bid.Value() {
			t.handRanges[pos] = intersectRanges(t.handRanges[pos], extractBounds(conv.node))
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

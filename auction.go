package passpasspass

import (
	"fmt"
	"slices"
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

// PartnershipEval captures the combined picture of self + partner.
type PartnershipEval struct {
	MinCombinedHCP int
	MaxCombinedHCP int
	SuitFit        [4]SuitRange
}

// NewPartnershipEval computes the combined picture from own hand and partner's inferred range.
// A partnership never holds more than the deck's MaxDeckHCP points or HandSize cards of a suit,
// so the combined ranges are capped there.
func NewPartnershipEval(self *Hand, partner HandRange) PartnershipEval {
	selfDist := self.Distribution()
	var fit [4]SuitRange
	for i := range 4 {
		fit[i] = SuitRange{
			Min: min(selfDist[i]+partner.SuitLength[i].Min, HandSize),
			Max: min(selfDist[i]+partner.SuitLength[i].Max, HandSize),
		}
	}
	return PartnershipEval{
		MinCombinedHCP: min(self.HCP()+partner.MinHCP, MaxDeckHCP),
		MaxCombinedHCP: min(self.HCP()+partner.MaxHCP, MaxDeckHCP),
		SuitFit:        fit,
	}
}

// AuctionContext is passed to every Condition during a live auction. It is one player's view:
// partner's model is seen through the player's own cards, the opponents' models are as shown.
type AuctionContext struct {
	SelfHand    *Hand
	Position    Direction
	History     History
	Contract    *Bid // the highest bid so far, nil before the opening
	Partner     HandModel
	LHO         HandModel
	RHO         HandModel
	Partnership PartnershipEval
}

// AuctionStep is one call of an auction together with the context the caller's sheet saw.
type AuctionStep struct {
	Context *AuctionContext
	Call    *Bid
}

// Auction records a complete bridge auction: the deal, all bids, and the final contract.
type Auction struct {
	Dealer   Direction     // who opened the bidding
	Hands    [4]*Hand      // hands indexed by Direction (North=0…West=3)
	Bids     []BidAction   // every bid in order, populated after Run
	Steps    []AuctionStep // the context of every call, in order, populated after Run
	Contract *Bid          // nil = passed out, populated after Run
	Declarer Direction     // who made the final non-pass bid, populated after Run
	state    *AuctionState
}

// NewAuction creates an Auction with the given dealer, hands, and the state its calls are
// recorded in.
func NewAuction(dealer Direction, hands []*Hand, state *AuctionState) *Auction {
	a := &Auction{
		Dealer: dealer,
		state:  state,
	}
	copy(a.Hands[:], hands)
	return a
}

// Run executes the full bidding sequence and populates Bids, Steps, Contract, and Declarer.
func (a *Auction) Run() {
	consecutivePasses := 0
	for pos := a.Dealer; !a.over(consecutivePasses); pos = (pos + 1) % 4 {
		ctx := a.state.Perspective(pos, a.Hands[pos]).Context()
		bid := a.state.Sheet(pos).Bid(ctx)

		// Bid legality: must be strictly above the current contract
		if bid != Pass && a.Contract != nil && bid.Value() <= a.Contract.Value() {
			bid = Pass
		}

		a.state.Record(pos, bid)
		a.Bids = a.state.History()
		a.Steps = append(a.Steps, AuctionStep{Context: ctx, Call: bid})

		if bid == Pass {
			consecutivePasses++
		} else {
			consecutivePasses = 0
			a.Contract = bid
			a.Declarer = pos
		}
	}
}

// over reports whether the auction has ended: four passes before an opening, or three after a
// bid.
func (a *Auction) over(consecutivePasses int) bool {
	if a.Contract == nil {
		return consecutivePasses == 4
	}
	return consecutivePasses == 3
}

// String shows all four hands followed by the bidding sequence and contract.
func (a *Auction) String() string {
	return a.HandsString() + a.BiddingString()
}

// HandsString shows all four hands with their HCP and distribution.
func (a *Auction) HandsString() string {
	var sb strings.Builder

	for i, dir := range allDirections {
		h := a.Hands[i]
		fmt.Fprintf(&sb, "%s (%d HCP %s):\n%s\n", dir, h.HCP(), h.Distribution(), h)
	}

	return sb.String()
}

// BiddingString shows the dealer, the bidding grid and the contract.
func (a *Auction) BiddingString() string {
	var sb strings.Builder

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

// boundsOf returns the bounds of m, or the unconstrained range when m is empty.
func boundsOf(m HandModel) HandRange {
	if b, ok := m.Bounds(); ok {
		return b
	}
	return DefaultHandRange()
}

func (d Direction) String() string {
	return allDirections[d]
}

func (p PartnershipEval) String() string {
	return fmt.Sprintf("%2d–%-2d HCP", p.MinCombinedHCP, p.MaxCombinedHCP) + suitRangesString(p.SuitFit)
}

// String shows what the caller's sheet saw, one labelled line each: own hand, the calls so far,
// what partner, LHO and RHO had shown, and the partnership evaluation.
func (c *AuctionContext) String() string {
	auction := "—"
	if len(c.History) > 0 {
		calls := make([]string, len(c.History))
		for i, ba := range c.History {
			calls[i] = ba.Bid.String()
		}
		auction = strings.Join(calls, " ")
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Hand         %d HCP %s\n", c.SelfHand.HCP(), c.SelfHand.Distribution())
	fmt.Fprintf(&sb, "Auction      %s\n", auction)
	if sameModel(c.Partner, seenWith(DefaultHandModel(), c.SelfHand)) {
		// Partner has shown nothing beyond what the player's own cards rule out.
		writeModel(&sb, "Partner", DefaultHandModel())
	} else {
		writeModel(&sb, "Partner", c.Partner)
	}
	writeModel(&sb, "LHO", c.LHO)
	writeModel(&sb, "RHO", c.RHO)
	fmt.Fprintf(&sb, "Partnership  %s\n", c.Partnership)
	return sb.String()
}

// sameModel reports whether m and n hold the same alternatives in the same order and the same
// conflict flag.
func sameModel(m, n HandModel) bool {
	return m.Conflict == n.Conflict && slices.Equal(m.Cases, n.Cases)
}

// writeModel writes m after label, one alternative per line aligned under the first, and marks a
// model that conflicts with the viewer's own hand.
func writeModel(sb *strings.Builder, label string, m HandModel) {
	if m.IsEmpty() {
		fmt.Fprintf(sb, "%-13sno hand fits\n", label)
		return
	}
	for i, r := range m.Cases {
		if i > 0 {
			label = ""
		}
		fmt.Fprintf(sb, "%-13s%s", label, r)
		if i == 0 && m.Conflict {
			sb.WriteString("  (conflicts with own hand)")
		}
		sb.WriteString("\n")
	}
}

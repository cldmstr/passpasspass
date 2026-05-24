package passpasspass

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// BidAction is one step in the auction.
type BidAction struct {
	Direction Direction
	Bid       *Bid
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
	History     []BidAction
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
	table    *Table      // internal; holds the convention sheet and inference state
}

// NewAuction randomly deals cards and picks a dealer, returning an Auction ready
// to run. Call Run to execute the bidding sequence.
func NewAuction(cs *ConventionSheet) *Auction {
	now := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(now))
	dealer := Direction(rng.Intn(4))
	hands := NewRandomDealer(now + 1).Deal()
	t := NewTable(cs, hands)
	a := &Auction{
		Dealer: dealer,
		Bids:   make([]BidAction, 0, 40),
		table:  t,
	}
	for i, h := range hands {
		a.Hands[i] = h
	}
	return a
}

// Run executes the full bidding sequence and populates Bids, Contract, and Declarer.
func (a *Auction) Run() {
	result := a.table.RunAuction(a.Dealer)
	a.Bids = result.Bids
	a.Contract = result.Contract
	a.Declarer = result.Declarer
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

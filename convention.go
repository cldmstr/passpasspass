package passpasspass

import (
	"fmt"
	"strconv"
)

type Bid struct {
	level int
	suit  Suit
}

func NewBid(level int, suit Suit) *Bid {
	return &Bid{level: level, suit: suit}
}

var Pass = NewBid(0, 0)

func NewBidFromString(s string) (*Bid, error) {
	if s == "Pass" || s == "pass" {
		return Pass, nil
	}
	if len(s) < 2 {
		return nil, fmt.Errorf("invalid bid: %q", s)
	}
	level, err := strconv.Atoi(s[:1])
	if err != nil || level < 1 || level > 7 {
		return nil, fmt.Errorf("invalid bid level in %q", s)
	}
	suit, err := NewSuitFromString(s[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid bid suit in %q", s)
	}
	return NewBid(level, suit), nil
}

func (b *Bid) Value() int {
	return 10*b.level + int(b.suit)
}

func (b *Bid) String() string {
	if b == Pass {
		return "Pass"
	}
	return fmt.Sprintf("%d%s", b.level, b.suit)
}

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

type Convention struct {
	Bid        *Bid
	conditions []Condition
	node       conditionNode
}

func NewConvention(bid *Bid, conditions ...Condition) *Convention {
	return &Convention{Bid: bid, conditions: conditions}
}

func (c *Convention) NextBid(ctx *AuctionContext) (bool, *Bid) {
	for _, cond := range c.conditions {
		if !cond(ctx) {
			return false, Pass
		}
	}
	return true, c.Bid
}

type ConventionSheet struct {
	sheet []*Convention
}

func NewConventionSheet() *ConventionSheet {
	return &ConventionSheet{sheet: make([]*Convention, 0, 35)}
}

func (cs *ConventionSheet) Add(convention *Convention) {
	cs.sheet = append(cs.sheet, convention)
}

func (cs *ConventionSheet) Bid(ctx *AuctionContext) *Bid {
	for _, c := range cs.sheet {
		if ok, bid := c.NextBid(ctx); ok {
			return bid
		}
	}
	return Pass
}

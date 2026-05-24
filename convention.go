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

package passpasspass

import (
	"fmt"
	"strconv"
)

// Bid is a call: a level from 1 to 7 and a strain, or Pass.
type Bid struct {
	level int
	suit  Suit
}

// NewBid returns the bid of the given level and strain.
func NewBid(level int, suit Suit) *Bid {
	return &Bid{level: level, suit: suit}
}

// Pass is the single Pass call; compare against it by identity.
var Pass = NewBid(0, 0)

// NewBidFromString parses a bid such as "1NT" or "2H", or "Pass".
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

// Value ranks bids by level, then by strain; a higher bid has a higher value.
func (b *Bid) Value() int {
	return 10*b.level + int(b.suit)
}

func (b *Bid) String() string {
	if b == Pass {
		return "Pass"
	}
	return fmt.Sprintf("%d%s", b.level, b.suit)
}

// Convention is a bid together with the conditions under which it is made.
type Convention struct {
	Bid        *Bid
	conditions []Condition
}

// NewConvention returns a convention that bids bid when all conditions pass.
func NewConvention(bid *Bid, conditions ...Condition) *Convention {
	return &Convention{Bid: bid, conditions: conditions}
}

// maxConventionCases caps the alternatives one convention can show, so that `all` over several
// `any` nodes cannot multiply them without limit. The loader rejects conventions above it.
const maxConventionCases = 64

// Shown returns the model of the hands with which the convention's call could be made in ctx:
// the hands for which every condition holds.
func (c *Convention) Shown(ctx *CallContext) HandModel {
	return intersectShown(ctx, true, c.conditions)
}

// NextBid reports whether every condition passes in ctx and, if so, returns the bid.
func (c *Convention) NextBid(ctx *AuctionContext) (bool, *Bid) {
	for _, cond := range c.conditions {
		if !cond.Eval(ctx) {
			return false, Pass
		}
	}
	return true, c.Bid
}

// ConventionSheet is an ordered list of conventions evaluated first-match.
type ConventionSheet struct {
	sheet []*Convention
}

// NewConventionSheet returns an empty sheet.
func NewConventionSheet() *ConventionSheet {
	return &ConventionSheet{sheet: make([]*Convention, 0, 35)}
}

// Add appends a convention to the end of the sheet.
func (cs *ConventionSheet) Add(convention *Convention) {
	cs.sheet = append(cs.sheet, convention)
}

// Bid returns the bid of the first convention that matches ctx, or Pass if none does.
func (cs *ConventionSheet) Bid(ctx *AuctionContext) *Bid {
	for _, c := range cs.sheet {
		if ok, bid := c.NextBid(ctx); ok {
			return bid
		}
	}
	return Pass
}

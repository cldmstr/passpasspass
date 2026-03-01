package passpasspass

import (
	"fmt"
)

type Bid struct {
	level int
	suit  Suit
}

func NewBid(level int, suit Suit) *Bid {
	return &Bid{
		level: level,
		suit:  suit,
	}
}

var Pass = NewBid(0, 0)

func (b Bid) Value() int {
	return 10*b.level + int(b.suit)
}

func (b Bid) String() string {
	return fmt.Sprintf("%d%s", b.level, b.suit)
}

type Convention struct {
	Bid        *Bid
	conditions []Condition
}

func NewConvention(bid *Bid, conditions ...Condition) *Convention {
	return &Convention{
		Bid:        bid,
		conditions: conditions,
	}
}

func (c *Convention) NextBid(hand *Hand) (bool, *Bid) {
	for _, c := range c.conditions {
		if !c(hand) {
			return false, nil
		}
	}

	return true, c.Bid
}

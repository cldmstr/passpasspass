package passpasspass

import (
	"fmt"
	"strconv"
)

const (
	Clubs     Suit = iota
	Diamonds
	Hearts
	Spades
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	default:
		panic("Invalid Suit")
	}
}

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Rank int

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case King:
		return "K"
	case Queen:
		return "Q"
	case Jack:
		return "J"
	default:
		// Check if it's a valid numeric rank (2-10)
		if r >= 2 && r <= 10 {
			return strconv.Itoa(int(r))
		}
		panic("Invalid Rank")
	}
}

type Card struct {
	suit Suit
	rank Rank 
}

func NewCard(suit Suit, rank Rank) *Card {
	return &Card {
		rank: rank,
		suit: suit,
	}
}

func (c *Card) HCP() int {
	switch c.rank {
	case Jack:
		return 1
	case Queen:
		return 2
	case King:
		return 3
	case Ace:
		return 4
	default:
		return 0
	}
}

func (c *Card) Value() int {
	return int(c.suit) * int(HandSize) + int(c.rank)
}

func (c *Card) HasSuit(s Suit) bool {
	return c.suit == s
}

func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.rank, c.suit)
}

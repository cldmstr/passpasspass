package passpasspass

import (
	"fmt"
	"slices"
	"strconv"
)

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
	NoTrump
)

type Suit int

var allSuits = []string{"♣️", "♦️", "❤️", "♠️", "NT"}

var suitLetters = []string{"C", "D", "H", "S", "NT"}

func NewSuitFromString(suit string) (Suit, error) {
	if i := slices.Index(allSuits, suit); i >= 0 {
		return Suit(i), nil
	}
	if i := slices.Index(suitLetters, suit); i >= 0 {
		return Suit(i), nil
	}
	return Suit(-1), fmt.Errorf("invalid suit: %s", suit)
}

func (s Suit) String() string {
	if int(s) > int(NoTrump) || int(s) < int(Clubs) {
		return "Invalid Suit"
	}

	return allSuits[int(s)]
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

var allRanks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Rank int

func NewRankFromString(rank string) (Rank, error) {
	i := slices.Index(allRanks, rank)
	if i < 0 {
		return Rank(0), fmt.Errorf("invalid rank: %s", rank)
	}

	return Rank(i + 2), nil
}

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
	return &Card{
		rank: rank,
		suit: suit,
	}
}

func NewCardFromString(card string) (*Card, error) {
	offset := len(card) - 6 // 6 is the length of the unicode suit markers
	r := card[:offset]
	s := card[offset:]

	rank, err := NewRankFromString(r)
	if err != nil {
		return nil, fmt.Errorf("invalid format in %q: %w", card, err)
	}
	suit, err := NewSuitFromString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid format in %q: %w", card, err)
	}

	return &Card{suit: suit, rank: rank}, nil
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
	return int(c.suit)*int(HandSize) + int(c.rank)
}

func (c *Card) HasSuit(s Suit) bool {
	return c.suit == s
}

func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.rank, c.suit)
}

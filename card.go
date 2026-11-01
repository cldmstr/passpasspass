package passpasspass

import (
	"fmt"
	"slices"
	"strconv"
)

// Suits in bidding order; NoTrump ranks above the four suits.
const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
	NoTrump
)

// Suit is a card suit or, in a bid, a strain including notrump.
type Suit int

var allSuits = []string{"♣️", "♦️", "❤️", "♠️", "NT"}

var suitLetters = []string{"C", "D", "H", "S", "NT"}
var suitNames = []string{"clubs", "diamonds", "hearts", "spades"}

// NewSuitFromString parses a suit from its emoji, letter (C D H S NT) or lower-case name.
func NewSuitFromString(suit string) (Suit, error) {
	if i := slices.Index(allSuits, suit); i >= 0 {
		return Suit(i), nil
	}
	if i := slices.Index(suitLetters, suit); i >= 0 {
		return Suit(i), nil
	}
	if i := slices.Index(suitNames, suit); i >= 0 {
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

// Ranks from two to ace; a rank's value is its number, with J=11 up to A=14.
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

// Rank is the rank of a card.
type Rank int

// NewRankFromString parses a rank from 2–10, J, Q, K or A.
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

// Card is a single playing card.
type Card struct {
	suit Suit
	rank Rank
}

// NewCard returns the card of the given suit and rank.
func NewCard(suit Suit, rank Rank) *Card {
	return &Card{
		rank: rank,
		suit: suit,
	}
}

// NewCardFromString parses a card written as rank followed by suit emoji, e.g. "A♣️" or "10❤️".
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

// HCP returns the card's high-card points: A=4, K=3, Q=2, J=1, otherwise 0.
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

// Value orders cards by suit, then by rank.
func (c *Card) Value() int {
	return int(c.suit)*HandSize + int(c.rank)
}

// HasSuit reports whether the card belongs to suit s.
func (c *Card) HasSuit(s Suit) bool {
	return c.suit == s
}

func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.rank, c.suit)
}

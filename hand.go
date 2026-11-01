package passpasspass

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Sizes and limits of a bridge hand.
const (
	DistributionLength int = 4
	HandSize           int = 13
	MaxHandHCP         int = 37
)

// Distribution holds a hand's suit lengths in suit order: clubs, diamonds, hearts, spades.
type Distribution []int

// NewDistribution returns a Distribution from the four suit lengths.
func NewDistribution(clubs, diamonds, hearts, spades int) Distribution {
	return []int{clubs, diamonds, hearts, spades}
}

func (d Distribution) String() string {
	if len(d) != DistributionLength {
		return "[invalid distribution]"
	}
	distStrings := make([]string, DistributionLength)
	for i, s := range d {
		distStrings[i] = strconv.Itoa(s)
	}

	return fmt.Sprintf("[%s]", strings.Join(distStrings, ","))
}

// Description summarises a hand by its HCP and distribution.
type Description struct {
	HCP          int
	Distribution Distribution
}

// Hand is a player's cards, kept sorted by suit and rank.
type Hand struct {
	cards []*Card
}

// NewHand returns an empty hand.
func NewHand() *Hand {
	return &Hand{
		cards: make([]*Card, 0, HandSize),
	}
}

// NewHandFromStrings parses 13 card strings into a hand and validates it.
func NewHandFromStrings(cards []string) (*Hand, error) {
	hand := NewHand()
	for _, card := range cards {
		c, err := NewCardFromString(card)
		if err != nil {
			return nil, fmt.Errorf("invalid hand description: %w", err)
		}
		hand.cards = append(hand.cards, c)
	}
	hand.sort()

	err := hand.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid hand: %w", err)
	}

	return hand, nil
}

// AddCard adds a card to the hand and keeps it sorted.
func (h *Hand) AddCard(card *Card) {
	h.cards = append(h.cards, card)
	h.sort()
}

// HCP returns the hand's total high-card points.
func (h *Hand) HCP() int {
	hcp := 0
	for _, c := range h.cards {
		hcp += c.HCP()
	}

	return hcp
}

// Description returns the hand's HCP and distribution.
func (h *Hand) Description() *Description {
	return &Description{
		HCP:          h.HCP(),
		Distribution: h.Distribution(),
	}
}

// Validate returns an error unless the hand holds exactly 13 distinct cards.
func (h *Hand) Validate() error {
	if len(h.cards) != HandSize {
		return fmt.Errorf("invalid hand size should be %d: %d", HandSize, len(h.cards))
	}

	encountered := map[Card]bool{}
	duplicates := []string{}

	for _, v := range h.cards {
		if encountered[*v] {
			duplicates = append(duplicates, v.String())
		} else {
			encountered[*v] = true
		}
	}
	if len(duplicates) > 0 {
		return fmt.Errorf("invalid hand contains duplicates: %v", duplicates)
	}

	return nil
}

func (h *Hand) String() string {
	result := ""
	lastOffset := 0
	for _, dist := range h.Distribution() {
		offset := lastOffset + dist
		for i, card := range h.cards[lastOffset:offset] {
			if i < offset {
				result += card.String()
			}
			if i < dist-1 {
				result += ","
			}
		}
		lastOffset = offset
		result += "\n"
	}

	return result
}

// Distribution returns the hand's suit lengths in suit order.
func (h *Hand) Distribution() Distribution {
	suitCount := make([]int, 4)
	suit := 0
	numSuit := 0
	for _, card := range h.cards {
		for !card.HasSuit(Suit(suit)) {
			suitCount[suit] = numSuit
			suit++
			numSuit = 0
			continue
		}
		numSuit++
	}
	suitCount[suit] = numSuit

	return suitCount
}

func (h *Hand) sort() {
	slices.SortFunc(h.cards, func(i, j *Card) int {
		return i.Value() - j.Value()
	})
}

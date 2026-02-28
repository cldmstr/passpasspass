package passpasspass

import (
	"fmt"
	"slices"
	"strings"
)

const (
	HandSize int = 13
)

type Distribution []int

func NewDistribution(clubs, diamonds, hearts, spades int) Distribution {
	return []int{clubs, diamonds, hearts, spades}
}

type Description struct {
	HCP          int
	Distribution Distribution
}

type Hand struct {
	cards []*Card
}

func NewHand() *Hand {
	return &Hand{
		cards: make([]*Card, 0, HandSize),
	}
}

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

func (h *Hand) AddCard(card *Card) {
	h.cards = append(h.cards, card)
	h.sort()
}

func (h *Hand) HCP() int {
	hcp := 0
	for _, c := range h.cards {
		hcp += c.HCP()
	}

	return hcp
}

func (h *Hand) Description() *Description {
	return &Description{
		HCP:          h.HCP(),
		Distribution: h.distribution(),
	}
}

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
	cards := make([]string, len(h.cards))
	suitOffsets := make([]int, 4)
	suitOffsets[3] = int(HandSize)
	suit := 0
	for i, card := range h.cards {
		if !card.HasSuit(Suit(suit)) {
			suitOffsets[suit] = i
			suit++
		}
		cards[i] = card.String()
	}

	result := ""
	lastOffset := 0
	for i, suitOffset := range suitOffsets {
		result += strings.Join(cards[lastOffset:suitOffset], ",")
		if i < len(suitOffsets)-1 {
			result += "\n"
		}
		lastOffset = suitOffset
	}

	return result
}

func (h *Hand) distribution() Distribution {
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

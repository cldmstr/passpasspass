package passpasspass

import (
	"slices"
	"strings"
)

const (
	HandSize int64 = 13
)

type Hand struct {
	cards []*Card
}

func NewHand() *Hand {
	return &Hand {
		cards: make([]*Card, 0, HandSize),
	}
}

func (h *Hand) AddCard(card *Card) {
	h.cards = append(h.cards, card)
}

func (h *Hand) String() string {
	h.sort()
	cards := make([]string, len(h.cards))
	suitOffsets := make([]int,4)
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

func (h *Hand) HCP() int {
	hcp := 0
	for _, c := range h.cards {
		hcp += c.HCP()
	}

	return hcp
}

func (h *Hand) sort() {
	slices.SortFunc(h.cards, func(i,j *Card) int {
		return i.Value() - j.Value()
	} )
}

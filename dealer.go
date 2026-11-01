package passpasspass

import (
	"math/rand"
)

// Deck maps positions 0–51 to cards.
type Deck map[int]*Card

// StandardDeck returns the 52 cards in suit and rank order.
func StandardDeck() Deck {
	d := map[int]*Card{}
	for s := range 4 {
		for r := range HandSize {
			card := NewCard(Suit(s), Rank(r+2))
			d[s*HandSize+r] = card
		}
	}

	return d
}

// Dealer deals four hands, indexed by Direction.
type Dealer interface {
	Deal() []*Hand
}

type randomDealer struct {
	r *rand.Rand
}

// NewRandomDealer returns a Dealer that shuffles with a random source seeded by seed.
func NewRandomDealer(seed int64) Dealer {
	return &randomDealer{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (d *randomDealer) Deal() []*Hand {
	deck := StandardDeck()
	d.r.Shuffle(HandSize*4, func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	suits := []int{0, HandSize, HandSize * 2, HandSize * 3}
	hands := make([]*Hand, 4)

	for i, offset := range suits {
		hand := NewHand()
		for j := range HandSize {
			hand.AddCard(deck[offset+j])
		}
		hands[i] = hand
	}

	return hands
}

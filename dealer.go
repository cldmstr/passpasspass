package passpasspass

import (
	"math/rand"
)

type Deck map[int]*Card

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

type Dealer interface {
	Deal() []*Hand
}

type randomDealer struct {
	r *rand.Rand
}

func NewRandomDealer(seed int64) Dealer {
	return &randomDealer{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (d *randomDealer) Deal() []*Hand {
	deck := StandardDeck()
	d.r.Shuffle(int(HandSize*4), func(i, j int) {
		hold := deck[i]
		deck[i] = deck[j]
		deck[j] = hold
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

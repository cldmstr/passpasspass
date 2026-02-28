package passpasspass

import (
	"fmt"
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

type Table struct {
	North *Hand
	South *Hand
	East  *Hand
	West  *Hand
}

func (t Table) String() string {
	return fmt.Sprintf("North (%d): \n%s\nEast (%d): \n%s\nSouth (%d):  \n%s\nWest (%d):  \n%s\n", t.North.HCP(), t.North, t.East.HCP(), t.East, t.South.HCP(), t.South, t.West.HCP(), t.West)
}

type Dealer interface {
	Deal() *Table
}

type randomDealer struct {
	r *rand.Rand
}

func NewRandomDealer(seed int64) Dealer {
	return &randomDealer{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (d *randomDealer) Deal() *Table {
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

	return &Table{
		North: hands[0],
		South: hands[1],
		East:  hands[2],
		West:  hands[3],
	}
}

package passpasspass

import (
	"math/rand"
	"time"
)

const (
	North Direction = iota
	East
	South
	West
)

type Direction int

var allDirections = []string{"North", "East", "South", "West"}

type Table struct {
	cs     *ConventionSheet
	dealer Dealer
}

func NewTable(cs *ConventionSheet, dealer Dealer) *Table {
	return &Table{cs: cs, dealer: dealer}
}

func (t *Table) RunAuction() *Auction {
	hands := t.dealer.Deal()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dir := Direction(rng.Intn(4))
	auction := NewAuction(dir, t.cs, hands)
	auction.Run()
	return auction
}

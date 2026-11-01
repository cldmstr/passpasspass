package passpasspass

import (
	"math/rand"
	"time"
)

// Seats in bidding order.
const (
	North Direction = iota
	East
	South
	West
)

// Direction is a seat at the table.
type Direction int

var allDirections = []string{"North", "East", "South", "West"}

// Table wires a dealer and a convention sheet together to run auctions.
type Table struct {
	cs     *ConventionSheet
	dealer Dealer
}

// NewTable returns a Table that deals with dealer and bids by cs.
func NewTable(cs *ConventionSheet, dealer Dealer) *Table {
	return &Table{cs: cs, dealer: dealer}
}

// RunAuction deals a new board, picks a random dealer, and runs the auction to its end.
func (t *Table) RunAuction() *Auction {
	hands := t.dealer.Deal()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dir := Direction(rng.Intn(4))
	auction := NewAuction(dir, t.cs, hands)
	auction.Run()
	return auction
}

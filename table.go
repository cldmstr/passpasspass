package passpasspass

import (
	"fmt"
	"strings"
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
	cs    *ConventionSheet
	hands []*Hand
}

func NewTable(cs *ConventionSheet, hands []*Hand) *Table {
	return &Table{
		cs:    cs,
		hands: hands,
	}
}

func (t Table) String() string {
	directionFormat := "%s (%d):\nBid: %s\n%s\n"
	sb := strings.Builder{}
	for i, d := range allDirections {
		hand := t.hands[int(i)]
		sb.WriteString(fmt.Sprintf(directionFormat, d, hand.HCP(), t.cs.Bid(hand), hand))
	}
	return sb.String()
}

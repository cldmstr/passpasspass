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
	directionFormat := "%s (%d : %s):\nBid: %s\n%s\n"
	sb := strings.Builder{}
	for i, d := range allDirections {
		h := t.hands[int(i)]
		sb.WriteString(fmt.Sprintf(directionFormat, d, h.HCP(), h.Distribution(), t.cs.Bid(h), h))
	}
	return sb.String()
}

package passpasspass

import (
	"fmt"
	"strings"
)

// HandRange holds inferred HCP and suit-length bounds for one player.
type HandRange struct {
	MinHCP, MaxHCP int
	SuitLength     [4]SuitRange // indexed clubs(0)…spades(3)
}

// String shows the HCP range and the length range of every suit, or "nothing shown" for the
// unconstrained range.
func (r HandRange) String() string {
	if r == DefaultHandRange() {
		return "nothing shown"
	}
	return fmt.Sprintf("%2d–%-2d HCP", r.MinHCP, r.MaxHCP) + suitRangesString(r.SuitLength)
}

// suitRangesString shows a length range per suit, from spades down to clubs.
func suitRangesString(suits [4]SuitRange) string {
	var sb strings.Builder
	for s := Spades; s >= Clubs; s-- {
		fmt.Fprintf(&sb, "  %s %d–%d", allSuits[s], suits[s].Min, suits[s].Max)
	}
	return sb.String()
}

// DefaultHandRange returns an unconstrained HandRange.
func DefaultHandRange() HandRange {
	return HandRange{
		MinHCP: 0, MaxHCP: MaxHandHCP,
		SuitLength: [4]SuitRange{
			{0, HandSize}, {0, HandSize}, {0, HandSize}, {0, HandSize},
		},
	}
}

// HandModel is what a player has shown: a list of alternative ranges, any one of which their
// hand may fit. A model without alternatives is empty: no hand fits it.
type HandModel struct {
	Cases []HandRange

	// Conflict marks a partner's model that no hand can fit as seen from the viewer's own cards.
	// It is kept as partner showed it, since partner has misbid or is lying.
	Conflict bool
}

// DefaultHandModel returns the model of a player who has shown nothing: one unconstrained range.
func DefaultHandModel() HandModel {
	return HandModel{Cases: []HandRange{DefaultHandRange()}}
}

// modelOf returns the model with the single alternative r, tightened.
func modelOf(r HandRange) HandModel {
	return HandModel{Cases: []HandRange{r}}.Tighten()
}

// guardModel returns the model of a condition that does not depend on the hand: every hand when
// it comes out as wanted, no hand otherwise.
func guardModel(asWanted bool) HandModel {
	if asWanted {
		return DefaultHandModel()
	}
	return HandModel{}
}

// unionModels returns a model holding the alternatives of every model in models.
func unionModels(models ...HandModel) HandModel {
	var cases []HandRange
	for _, m := range models {
		cases = append(cases, m.Cases...)
	}
	return HandModel{Cases: cases}
}

// hcpRange returns the unconstrained range narrowed to lo–hi HCP.
func hcpRange(lo, hi int) HandRange {
	r := DefaultHandRange()
	r.MinHCP, r.MaxHCP = lo, hi
	return r
}

// suitRange returns the unconstrained range narrowed to lo–hi cards in suit s.
func suitRange(s, lo, hi int) HandRange {
	r := DefaultHandRange()
	r.SuitLength[s] = SuitRange{Min: lo, Max: hi}
	return r
}

// IsEmpty reports whether no hand fits the model.
func (m HandModel) IsEmpty() bool {
	return len(m.Cases) == 0
}

// Intersect returns the model of hands that fit both m and o: the intersection of every
// alternative of m with every alternative of o, tightened, without those no hand can fit.
func (m HandModel) Intersect(o HandModel) HandModel {
	var cases []HandRange
	for _, a := range m.Cases {
		for _, b := range o.Cases {
			if r, ok := tightenRange(intersectRanges(a, b)); ok {
				cases = append(cases, r)
			}
		}
	}
	return HandModel{Cases: cases}
}

// Tighten applies the 13-card rule and the HCP limit of one hand to every alternative, and drops
// the alternatives no hand can fit.
func (m HandModel) Tighten() HandModel {
	var cases []HandRange
	for _, c := range m.Cases {
		if r, ok := tightenRange(c); ok {
			cases = append(cases, r)
		}
	}
	return HandModel{Cases: cases}
}

// Bounds returns the smallest single range that covers every alternative, and false for an
// empty model.
func (m HandModel) Bounds() (HandRange, bool) {
	if m.IsEmpty() {
		return HandRange{}, false
	}
	b := m.Cases[0]
	for _, c := range m.Cases[1:] {
		b = unionRanges(b, c)
	}
	return b, true
}

// tightenRange narrows r to what a 13-card hand can hold: no suit longer than 13 minus the other
// suits' minimums, no suit shorter than 13 minus the other suits' maximums, and at most
// MaxHandHCP. Narrowing one suit can narrow another, so it repeats until nothing changes. It
// reports false when no hand fits r.
func tightenRange(r HandRange) (HandRange, bool) {
	r.MaxHCP = min(r.MaxHCP, MaxHandHCP)
	if r.MinHCP > r.MaxHCP {
		return r, false
	}
	for changed := true; changed; {
		changed = false
		sumMin, sumMax := 0, 0
		for _, s := range r.SuitLength {
			sumMin += s.Min
			sumMax += s.Max
		}
		for i, s := range r.SuitLength {
			tight := SuitRange{
				Min: max(s.Min, HandSize-(sumMax-s.Max)),
				Max: min(s.Max, HandSize-(sumMin-s.Min)),
			}
			if tight.Min > tight.Max {
				return r, false
			}
			changed = changed || tight != s
			r.SuitLength[i] = tight
		}
	}
	return r, true
}

func intersectRanges(a, b HandRange) HandRange {
	r := a
	if b.MinHCP > r.MinHCP {
		r.MinHCP = b.MinHCP
	}
	if b.MaxHCP < r.MaxHCP {
		r.MaxHCP = b.MaxHCP
	}
	for i := range 4 {
		if b.SuitLength[i].Min > r.SuitLength[i].Min {
			r.SuitLength[i].Min = b.SuitLength[i].Min
		}
		if b.SuitLength[i].Max < r.SuitLength[i].Max {
			r.SuitLength[i].Max = b.SuitLength[i].Max
		}
	}
	return r
}

// unionRanges returns the smallest range that holds both a and b.
func unionRanges(a, b HandRange) HandRange {
	r := a
	r.MinHCP = min(r.MinHCP, b.MinHCP)
	r.MaxHCP = max(r.MaxHCP, b.MaxHCP)
	for i := range 4 {
		r.SuitLength[i].Min = min(r.SuitLength[i].Min, b.SuitLength[i].Min)
		r.SuitLength[i].Max = max(r.SuitLength[i].Max, b.SuitLength[i].Max)
	}
	return r
}

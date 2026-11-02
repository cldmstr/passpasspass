package passpasspass

import "fmt"

// patternKind tells a call from the wildcards in an auction pattern.
type patternKind int

const (
	patternCall     patternKind = iota // exactly the given call
	patternAny                         // exactly one call of any kind
	patternAnyStar                     // zero or more calls of any kind
	patternPassStar                    // zero or more passes
)

// patternWildcards maps each wildcard's spelling in a convention file to its kind.
var patternWildcards = map[string]patternKind{
	"Any":   patternAny,
	"Any*":  patternAnyStar,
	"Pass*": patternPassStar,
}

// patternElement is one element of an AuctionPattern: a call, or a wildcard.
type patternElement struct {
	kind patternKind
	call *Bid // set for patternCall only
}

// AuctionPattern matches a whole auction, from the dealer's first call to the most recent one,
// like a glob over calls anchored at both ends.
type AuctionPattern struct {
	elements []patternElement
}

// NewAuctionPattern parses a pattern from its elements: calls such as "1H" or "Pass", and the
// wildcards "Any" (one call), "Any*" (zero or more calls) and "Pass*" (zero or more passes).
func NewAuctionPattern(elements ...string) (AuctionPattern, error) {
	p := AuctionPattern{elements: make([]patternElement, 0, len(elements))}
	for _, s := range elements {
		e, err := parsePatternElement(s)
		if err != nil {
			return AuctionPattern{}, err
		}
		p.elements = append(p.elements, e)
	}
	return p, nil
}

// parsePatternElement parses one element: a wildcard, or else a call.
func parsePatternElement(s string) (patternElement, error) {
	if kind, ok := patternWildcards[s]; ok {
		return patternElement{kind: kind}, nil
	}
	call, err := NewBidFromString(s)
	if err != nil {
		return patternElement{}, fmt.Errorf("%q is neither a call nor a wildcard (Any, Any*, Pass*)", s)
	}
	return patternElement{kind: patternCall, call: call}, nil
}

// Matches reports whether the pattern covers every call in h.
func (p AuctionPattern) Matches(h History) bool {
	return matchElements(p.elements, h)
}

func matchElements(elements []patternElement, h History) bool {
	if len(elements) == 0 {
		return len(h) == 0
	}
	e, rest := elements[0], elements[1:]
	if e.kind == patternAnyStar || e.kind == patternPassStar {
		return matchStar(e, rest, h)
	}
	return len(h) > 0 && e.matchesCall(h[0].Bid) && matchElements(rest, h[1:])
}

// matchStar lets the star wildcard e cover as few calls as possible, then as many as it has to
// for the rest of the pattern to match.
func matchStar(e patternElement, rest []patternElement, h History) bool {
	for i := range len(h) + 1 {
		if matchElements(rest, h[i:]) {
			return true
		}
		if i < len(h) && !e.matchesCall(h[i].Bid) {
			return false
		}
	}
	return false
}

// matchesCall reports whether e accepts call as one of the calls it covers.
func (e patternElement) matchesCall(call *Bid) bool {
	switch e.kind {
	case patternCall:
		return call.Value() == e.call.Value()
	case patternPassStar:
		return call == Pass
	default:
		return true
	}
}

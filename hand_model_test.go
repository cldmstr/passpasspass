package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// handRange returns the unconstrained range narrowed to minHCP–maxHCP and to the given suits.
func handRange(minHCP, maxHCP int, suits map[ppp.Suit]ppp.SuitRange) ppp.HandRange {
	r := ppp.DefaultHandRange()
	r.MinHCP, r.MaxHCP = minHCP, maxHCP
	for s, sr := range suits {
		r.SuitLength[s] = sr
	}
	return r
}

// model returns a model with the given alternatives.
func model(cases ...ppp.HandRange) ppp.HandModel {
	return ppp.HandModel{Cases: cases}
}

func TestHandModel_Intersect(t *testing.T) {
	weakHearts := handRange(6, 10, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {6, 13}})
	weakSpades := handRange(6, 10, map[ppp.Suit]ppp.SuitRange{ppp.Spades: {6, 13}})
	strongBalanced := handRange(23, 24, map[ppp.Suit]ppp.SuitRange{
		ppp.Clubs: {2, 5}, ppp.Diamonds: {2, 5}, ppp.Hearts: {2, 5}, ppp.Spades: {2, 5},
	})
	multi := model(strongBalanced, weakHearts, weakSpades, handRange(27, 37, nil))

	t.Run("disjoint alternatives give an empty model", func(t *testing.T) {
		is := is.New(t)
		got := model(handRange(0, 9, nil)).Intersect(model(handRange(15, 17, nil)))
		is.True(got.IsEmpty())
	})

	t.Run("alternatives no 13-card hand can fit are dropped", func(t *testing.T) {
		is := is.New(t)
		sevenHearts := handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {7, 13}})
		sevenSpades := handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Spades: {7, 13}})
		is.True(model(sevenHearts).Intersect(model(sevenSpades)).IsEmpty())
	})

	t.Run("every pairwise intersection that a hand can fit is kept", func(t *testing.T) {
		is := is.New(t)
		got := multi.Intersect(model(
			handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Spades: {6, 13}}),
			handRange(20, 37, nil),
		))
		// Pairs in order: strong balanced keeps only 20+ HCP (it has 2–5 spades); weak hearts and
		// weak spades keep only 6+ spades (they have 6–10 HCP); 27+ keeps both.
		is.Equal(len(got.Cases), 5)
		is.Equal(got.Cases[0].MinHCP, 23)                                                   // strong balanced, 20+ HCP
		is.Equal(got.Cases[1].SuitLength, [4]ppp.SuitRange{{0, 1}, {0, 1}, {6, 7}, {6, 7}}) // weak hearts, 6+ spades
		is.Equal(got.Cases[2].SuitLength[ppp.Spades], ppp.SuitRange{Min: 6, Max: 13})       // weak spades
		is.Equal(got.Cases[2].MaxHCP, 10)
		is.Equal(got.Cases[3].SuitLength[ppp.Spades], ppp.SuitRange{Min: 6, Max: 13}) // strong, 6+ spades
		is.Equal(got.Cases[3].MinHCP, 27)
		is.Equal(got.Cases[4].MinHCP, 27) // strong, 20+ HCP
	})

	t.Run("the default model keeps the alternatives", func(t *testing.T) {
		is := is.New(t)
		got := ppp.DefaultHandModel().Intersect(model(weakHearts))
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].MinHCP, 6)
		is.Equal(got.Cases[0].MaxHCP, 10)
		is.Equal(got.Cases[0].SuitLength[ppp.Hearts], ppp.SuitRange{Min: 6, Max: 13})
	})

	t.Run("an empty model stays empty", func(t *testing.T) {
		is := is.New(t)
		is.True(model().Intersect(multi).IsEmpty())
		is.True(multi.Intersect(model()).IsEmpty())
	})
}

func TestHandModel_Tighten(t *testing.T) {
	t.Run("6+ hearts and 6+ spades limit both minors to 0–1", func(t *testing.T) {
		is := is.New(t)
		got := model(handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {6, 13}, ppp.Spades: {6, 13}})).Tighten()
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].SuitLength, [4]ppp.SuitRange{{0, 1}, {0, 1}, {6, 7}, {6, 7}})
	})

	t.Run("short suits force a long one", func(t *testing.T) {
		is := is.New(t)
		// At most 2 clubs, 2 diamonds and 3 hearts leaves at least 6 spades.
		got := model(handRange(0, 37, map[ppp.Suit]ppp.SuitRange{
			ppp.Clubs: {0, 2}, ppp.Diamonds: {0, 2}, ppp.Hearts: {0, 3},
		})).Tighten()
		is.Equal(got.Cases[0].SuitLength[ppp.Spades], ppp.SuitRange{Min: 6, Max: 13})
	})

	t.Run("narrowing one suit narrows another", func(t *testing.T) {
		is := is.New(t)
		// Exactly 4 clubs and exactly 4 diamonds leave 5 cards for the majors; at least 3 hearts
		// then leaves at most 2 spades, and at most 2 spades leaves at least 3 hearts.
		got := model(handRange(0, 37, map[ppp.Suit]ppp.SuitRange{
			ppp.Clubs: {4, 4}, ppp.Diamonds: {4, 4}, ppp.Hearts: {3, 13},
		})).Tighten()
		is.Equal(got.Cases[0].SuitLength, [4]ppp.SuitRange{{4, 4}, {4, 4}, {3, 5}, {0, 2}})
	})

	t.Run("HCP is capped at 37", func(t *testing.T) {
		is := is.New(t)
		got := model(handRange(30, 40, nil)).Tighten()
		is.Equal(got.Cases[0].MaxHCP, 37)
	})

	t.Run("alternatives no hand can fit are dropped", func(t *testing.T) {
		is := is.New(t)
		got := model(
			handRange(38, 40, nil), // more HCP than one hand holds
			handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {7, 13}, ppp.Spades: {7, 13}}), // 14+ cards
			handRange(0, 37, map[ppp.Suit]ppp.SuitRange{
				ppp.Clubs: {0, 3}, ppp.Diamonds: {0, 3}, ppp.Hearts: {0, 3}, ppp.Spades: {0, 3},
			}), // at most 12 cards
			handRange(6, 10, nil),
		).Tighten()
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].MinHCP, 6)
	})

	t.Run("the default model is unchanged", func(t *testing.T) {
		is := is.New(t)
		is.Equal(ppp.DefaultHandModel().Tighten(), ppp.DefaultHandModel())
	})

	t.Run("tightening twice changes nothing", func(t *testing.T) {
		is := is.New(t)
		once := model(handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {6, 13}, ppp.Spades: {5, 13}})).Tighten()
		is.Equal(once.Tighten(), once)
	})
}

func TestHandModel_Bounds(t *testing.T) {
	t.Run("covers all alternatives", func(t *testing.T) {
		is := is.New(t)
		m := model(
			handRange(6, 10, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {6, 7}, ppp.Spades: {0, 3}}),
			handRange(20, 22, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {2, 4}, ppp.Spades: {6, 7}}),
		)
		got, ok := m.Bounds()
		is.True(ok)
		is.Equal(got.MinHCP, 6)
		is.Equal(got.MaxHCP, 22)
		is.Equal(got.SuitLength[ppp.Hearts], ppp.SuitRange{Min: 2, Max: 7})
		is.Equal(got.SuitLength[ppp.Spades], ppp.SuitRange{Min: 0, Max: 7})
		is.Equal(got.SuitLength[ppp.Clubs], ppp.SuitRange{Min: 0, Max: 13})
	})

	t.Run("one alternative is its own bounds", func(t *testing.T) {
		is := is.New(t)
		r := handRange(15, 17, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {2, 5}})
		got, ok := model(r).Bounds()
		is.True(ok)
		is.Equal(got, r)
	})

	t.Run("an empty model has no bounds", func(t *testing.T) {
		is := is.New(t)
		_, ok := model().Bounds()
		is.True(!ok)
	})
}

func TestHandModel_IsEmpty(t *testing.T) {
	is := is.New(t)
	is.True(model().IsEmpty())
	is.True(!ppp.DefaultHandModel().IsEmpty())
}

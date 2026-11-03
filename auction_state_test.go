package passpasspass_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// sixSpades holds A♠️(4)+K♠️(3)+Q♠️(2) = 9 HCP and 6 spades.
var sixSpades = []string{"A♠️", "K♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "2❤️", "3❤️"}

// eightSpades holds A♠️(4)+K♠️(3)+Q♠️(2)+J♠️(1) = 10 HCP and 8 spades.
var eightSpades = []string{"A♠️", "K♠️", "Q♠️", "J♠️", "2♠️", "3♠️", "4♠️", "5♠️", "2♣️", "3♣️", "2♦️", "2❤️", "3❤️"}

// stateWith returns a state in which every seat bids by cs and each player has shown the model
// given for them, or nothing.
func stateWith(cs *ppp.ConventionSheet, shown map[ppp.Direction]ppp.HandModel) *ppp.AuctionState {
	models := nothingShown()
	for d, m := range shown {
		models[d] = m
	}
	return ppp.NewAuctionState(ppp.NewSheetInterpreter(), [4]*ppp.ConventionSheet{cs, cs, cs, cs}, models)
}

func TestAuctionState_Perspective(t *testing.T) {
	spadesFrom := func(lo int) ppp.HandModel {
		return model(handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Spades: {lo, 13}}))
	}

	t.Run("own cards limit partner's suit", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(ppp.NewConventionSheet(), map[ppp.Direction]ppp.HandModel{ppp.North: spadesFrom(5)})
		partner := state.Perspective(ppp.South, mustHand(sixSpades)).Partner
		is.Equal(len(partner.Cases), 1)
		is.Equal(partner.Cases[0].SuitLength[ppp.Spades], ppp.SuitRange{Min: 5, Max: 7}) // 13 − 6
		is.True(!partner.Conflict)
	})

	t.Run("own HCP limit partner's HCP", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(ppp.NewConventionSheet(), nil)
		// A♣️(4)+K♣️(3)+Q♣️(2)+A♦️(4)+K♦️(3) = 16 HCP
		strong := mustHand([]string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"})
		is.Equal(state.Perspective(ppp.South, strong).Partner.Cases[0].MaxHCP, 24) // 40 − 16

		// J♣️(1)+J♦️(1) = 2 HCP leaves 38, more than one hand holds.
		weak := mustHand([]string{"J♣️", "2♣️", "3♣️", "4♣️", "J♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"})
		is.Equal(state.Perspective(ppp.South, weak).Partner.Cases[0].MaxHCP, 37)
	})

	t.Run("own cards that leave no hand for partner flag a conflict", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(ppp.NewConventionSheet(), map[ppp.Direction]ppp.HandModel{ppp.North: spadesFrom(6)})
		p := state.Perspective(ppp.South, mustHand(eightSpades))
		is.True(p.Partner.Conflict)
		is.Equal(p.Partner.Cases, spadesFrom(6).Cases) // kept as partner showed it
		is.True(strings.Contains(p.Context().String(), "(conflicts with own hand)"))
	})

	t.Run("opponents are not limited by own cards", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(ppp.NewConventionSheet(), map[ppp.Direction]ppp.HandModel{
			ppp.East: spadesFrom(5), ppp.West: spadesFrom(6),
		})
		p := state.Perspective(ppp.South, mustHand(eightSpades))
		is.Equal(p.RHO, spadesFrom(5)) // East
		is.Equal(p.LHO, spadesFrom(6)) // West, impossible from South's hand but not partner
	})

	t.Run("the current contract reaches the context", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(ppp.NewConventionSheet(), nil)
		for i, call := range []string{"1H", "2C", "Pass"} {
			bid, err := ppp.NewBidFromString(call)
			is.NoErr(err)
			state.Record(ppp.Direction(i), bid)
		}
		ctx := state.Perspective(ppp.West, mustHand(sixSpades)).Context()
		is.Equal(ctx.Contract.String(), "2♣️")
		is.Equal(len(ctx.History), 3)
	})

	t.Run("no contract before the opening", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(ppp.NewConventionSheet(), nil)
		state.Record(ppp.North, ppp.Pass)
		is.True(state.Perspective(ppp.East, mustHand(sixSpades)).Context().Contract == nil)
	})
}

func TestAuctionState_Record(t *testing.T) {
	record := func(t *testing.T, state *ppp.AuctionState, d ppp.Direction, call string) {
		t.Helper()
		bid, err := ppp.NewBidFromString(call)
		if err != nil {
			t.Fatal(err)
		}
		state.Record(d, bid)
	}
	sheet := sheetOf(t,
		[2]string{"1C", `hcp: [12, 21]`},
		[2]string{"2C", `hcp: [15, 30]`},
		[2]string{"3C", `hcp: [0, 5]`},
	)

	t.Run("a call narrows the caller's model", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(sheet, nil)
		record(t, state, ppp.North, "1C")
		record(t, state, ppp.East, "Pass")
		record(t, state, ppp.South, "Pass")
		record(t, state, ppp.West, "Pass")
		record(t, state, ppp.North, "2C")
		is.Equal(state.Shown(ppp.North).Cases, []ppp.HandRange{handRange(15, 21, nil)})
	})

	t.Run("a pass leaves the model unchanged", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(sheet, nil)
		record(t, state, ppp.North, "1C")
		record(t, state, ppp.East, "Pass")
		is.Equal(state.Shown(ppp.North).Cases, []ppp.HandRange{handRange(12, 21, nil)})
		is.Equal(state.Shown(ppp.East), ppp.DefaultHandModel())
	})

	t.Run("a contradicting call replaces the caller's model", func(t *testing.T) {
		is := is.New(t)
		state := stateWith(sheet, nil)
		record(t, state, ppp.North, "1C")
		record(t, state, ppp.East, "Pass")
		record(t, state, ppp.South, "Pass")
		record(t, state, ppp.West, "Pass")
		record(t, state, ppp.North, "3C")
		is.Equal(state.Shown(ppp.North).Cases, []ppp.HandRange{handRange(0, 5, nil)})
	})
}

// TestAuctionContext_String_Alternatives follows the cli spec's "Several alternatives" scenario.
func TestAuctionContext_String_Alternatives(t *testing.T) {
	is := is.New(t)
	state := stateWith(colinSheet(t), nil)
	state.Record(ppp.North, ppp.NewBid(2, ppp.Diamonds)) // Multi
	state.Record(ppp.East, ppp.Pass)

	text := state.Perspective(ppp.South, mustHand(sixSpades)).Context().String()
	lines := strings.Split(text, "\n")
	start := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "Partner ") {
			start = i
		}
	}
	is.True(strings.HasPrefix(lines[start], "Partner      23–24 HCP"))
	is.True(strings.HasPrefix(lines[start+1], "              6–10 HCP  ♠️ 0–7  ❤️ 6–11"))
	is.True(strings.HasPrefix(lines[start+2], "              6–10 HCP  ♠️ 6–7"))
	is.True(strings.HasPrefix(lines[start+3], "             27–31 HCP"))
	is.True(strings.HasPrefix(lines[start+4], "LHO"))
}

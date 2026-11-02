package passpasspass_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// mustHand builds a Hand from card strings, panicking on error.
func mustHand(cards []string) *ppp.Hand {
	h, err := ppp.NewHandFromStrings(cards)
	if err != nil {
		panic(err)
	}
	return h
}

// oneSuitedDeal returns four non-overlapping 13-card hands, one suit each.
// North=clubs, East=diamonds, South=hearts, West=spades.
func oneSuitedDeal() []*ppp.Hand {
	return []*ppp.Hand{
		mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"}),
		mustHand([]string{"2♦️", "3♦️", "4♦️", "5♦️", "6♦️", "7♦️", "8♦️", "9♦️", "10♦️", "J♦️", "Q♦️", "K♦️", "A♦️"}),
		mustHand([]string{"2❤️", "3❤️", "4❤️", "5❤️", "6❤️", "7❤️", "8❤️", "9❤️", "10❤️", "J❤️", "Q❤️", "K❤️", "A❤️"}),
		mustHand([]string{"2♠️", "3♠️", "4♠️", "5♠️", "6♠️", "7♠️", "8♠️", "9♠️", "10♠️", "J♠️", "Q♠️", "K♠️", "A♠️"}),
	}
}

func TestAuction_Run_PassedOut(t *testing.T) {
	is := is.New(t)

	cs := ppp.NewConventionSheet() // empty — every player passes
	a := ppp.NewAuction(ppp.North, cs, oneSuitedDeal())
	a.Run()

	is.True(a.Contract == nil)
	is.Equal(4, len(a.Bids))
	for _, ba := range a.Bids {
		is.Equal(ppp.Pass, ba.Bid)
	}
}

func TestAuction_Run_SingleOpener(t *testing.T) {
	is := is.New(t)

	// North: balanced 4-3-3-3, 16 HCP (A♣K♣Q♣ + A♦K♦).
	// East/South: 0–2 HCP. West: 22 HCP (above the 1NT ceiling).
	// Deck is fully partitioned across the four hands.
	hands := []*ppp.Hand{
		mustHand([]string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"}),
		mustHand([]string{"3♣️", "4♣️", "5♣️", "3♦️", "4♦️", "5♦️", "5❤️", "6❤️", "7❤️", "8❤️", "5♠️", "6♠️", "7♠️"}),
		mustHand([]string{"6♣️", "7♣️", "8♣️", "6♦️", "7♦️", "8♦️", "9❤️", "10❤️", "J❤️", "8♠️", "9♠️", "10♠️", "J♠️"}),
		mustHand([]string{"9♣️", "10♣️", "J♣️", "9♦️", "10♦️", "J♦️", "Q♦️", "Q❤️", "K❤️", "A❤️", "Q♠️", "K♠️", "A♠️"}),
	}

	cs := ppp.NewConventionSheet()
	cs.Add(ppp.NewConvention(
		ppp.NewBid(1, ppp.NoTrump),
		ppp.NewBalancedDistributionCondition(),
		ppp.NewHCPIntervalCondition(15, 17),
	))

	a := ppp.NewAuction(ppp.North, cs, hands)
	a.Run()

	is.Equal("1NT", a.Contract.String())
	is.Equal(ppp.North, a.Declarer)
	assertCalls(t, a.Bids, []call{
		{ppp.North, "1NT"},  // North opens
		{ppp.East, "Pass"},  // East: 0 HCP
		{ppp.South, "Pass"}, // South: 2 HCP
		{ppp.West, "Pass"},  // West: 22 HCP > ceiling
	})

	// Every call keeps the context it was made in.
	is.Equal(len(a.Steps), 4)
	for i, step := range a.Steps {
		is.Equal(step.Call, a.Bids[i].Bid)
		is.Equal(step.Context.Position, a.Bids[i].Direction)
		is.Equal(len(step.Context.History), i) // the calls before this one
	}
	is.Equal(a.Steps[0].Context.Partner.String(), "nothing shown")
	is.Equal(a.Steps[2].Context.Partner.String(), "15–17 HCP  ♠️ 2–5  ❤️ 2–5  ♦️ 2–5  ♣️ 2–5") // South sees North's 1NT
	is.Equal(a.Steps[1].Context.RHO.String(), "15–17 HCP  ♠️ 2–5  ❤️ 2–5  ♦️ 2–5  ♣️ 2–5")     // and so does East
}

func TestAuctionContext_String(t *testing.T) {
	is := is.New(t)
	hand := mustHand([]string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"})
	partner := ppp.DefaultHandRange()
	partner.MinHCP, partner.MaxHCP = 6, 10
	partner.SuitLength[ppp.Hearts] = ppp.SuitRange{Min: 6, Max: 13}
	ctx := &ppp.AuctionContext{
		SelfHand:    hand,
		Position:    ppp.South,
		History:     calls(t, "2H", "Pass"),
		Partner:     partner,
		LHO:         ppp.DefaultHandRange(),
		RHO:         ppp.DefaultHandRange(),
		Partnership: ppp.NewPartnershipEval(hand, partner),
	}
	is.Equal(ctx.String(), `Hand         16 HCP [4,3,3,3]
Auction      2❤️ Pass
Partner       6–10 HCP  ♠️ 0–13  ❤️ 6–13  ♦️ 0–13  ♣️ 0–13
LHO          nothing shown
RHO          nothing shown
Partnership  22–26 HCP  ♠️ 3–13  ❤️ 9–13  ♦️ 3–13  ♣️ 4–13
`)
}

func TestAuctionContext_String_BeforeTheFirstCall(t *testing.T) {
	is := is.New(t)
	hand := mustHand([]string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"})
	ctx := &ppp.AuctionContext{SelfHand: hand, Partnership: ppp.NewPartnershipEval(hand, ppp.DefaultHandRange())}
	is.True(strings.Contains(ctx.String(), "Auction      —\n"))
}

func TestNewPartnershipEval(t *testing.T) {
	is := is.New(t)

	// Self: clubs void, 4-3-3-3, 10 HCP (A♣K♣Q♣J♣).
	self := mustHand([]string{
		"A♣️", "K♣️", "Q♣️", "J♣️",
		"2♦️", "3♦️", "4♦️",
		"2❤️", "3❤️", "4❤️",
		"2♠️", "3♠️", "4♠️",
	})

	partner := ppp.HandRange{
		MinHCP: 6, MaxHCP: 10,
		SuitLength: [4]ppp.SuitRange{
			{Min: 2, Max: 5}, // clubs
			{Min: 2, Max: 4}, // diamonds
			{Min: 2, Max: 4}, // hearts
			{Min: 2, Max: 4}, // spades
		},
	}

	eval := ppp.NewPartnershipEval(self, partner)

	is.Equal(16, eval.MinCombinedHCP)                                   // 10 + 6
	is.Equal(20, eval.MaxCombinedHCP)                                   // 10 + 10
	is.Equal(ppp.SuitRange{Min: 6, Max: 9}, eval.SuitFit[ppp.Clubs])    // 4 + [2,5]
	is.Equal(ppp.SuitRange{Min: 5, Max: 7}, eval.SuitFit[ppp.Diamonds]) // 3 + [2,4]
	is.Equal(ppp.SuitRange{Min: 5, Max: 7}, eval.SuitFit[ppp.Hearts])   // 3 + [2,4]
	is.Equal(ppp.SuitRange{Min: 5, Max: 7}, eval.SuitFit[ppp.Spades])   // 3 + [2,4]
}

func TestNewPartnershipEval_Caps(t *testing.T) {
	// Self: 7 hearts, A♣️(4)+K♣️(3)+Q♣️(2)+A♦️(4)+K♦️(3)+A❤️(4)+K❤️(3)+Q❤️(2) = 25 HCP.
	self := mustHand([]string{
		"A♣️", "K♣️", "Q♣️",
		"A♦️", "K♦️",
		"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "5❤️",
		"2♠️",
	})

	t.Run("opposite a player who has shown nothing", func(t *testing.T) {
		is := is.New(t)
		eval := ppp.NewPartnershipEval(self, ppp.DefaultHandRange())
		is.Equal(25, eval.MinCombinedHCP)                                    // 25 + 0
		is.Equal(40, eval.MaxCombinedHCP)                                    // 25 + 37, capped at the deck's 40
		is.Equal(ppp.SuitRange{Min: 7, Max: 13}, eval.SuitFit[ppp.Hearts])   // 7 + [0,13], capped at 13
		is.Equal(ppp.SuitRange{Min: 1, Max: 13}, eval.SuitFit[ppp.Spades])   // 1 + [0,13], capped at 13
		is.Equal(ppp.SuitRange{Min: 3, Max: 13}, eval.SuitFit[ppp.Clubs])    // 3 + [0,13], capped at 13
		is.Equal(ppp.SuitRange{Min: 2, Max: 13}, eval.SuitFit[ppp.Diamonds]) // 2 + [0,13], capped at 13
	})

	t.Run("minimums are capped too", func(t *testing.T) {
		is := is.New(t)
		partner := ppp.DefaultHandRange()
		partner.MinHCP, partner.MaxHCP = 20, 22 // more than the 15 left in the deck
		partner.SuitLength[ppp.Hearts] = ppp.SuitRange{Min: 7, Max: 8}
		eval := ppp.NewPartnershipEval(self, partner)
		is.Equal(40, eval.MinCombinedHCP)
		is.Equal(40, eval.MaxCombinedHCP)
		is.Equal(ppp.SuitRange{Min: 13, Max: 13}, eval.SuitFit[ppp.Hearts])
	})
}

func TestAuction_Run_BidLegalityEnforcement(t *testing.T) {
	is := is.New(t)

	// All four players have 10 HCP and qualify for 1♥.
	// After North opens, East/South/West try the same bid — it must be coerced to Pass.
	cs := ppp.NewConventionSheet()
	cs.Add(ppp.NewConvention(
		ppp.NewBid(1, ppp.Hearts),
		ppp.NewHCPIntervalCondition(10, 37),
	))

	a := ppp.NewAuction(ppp.North, cs, oneSuitedDeal())
	a.Run()

	is.Equal("1❤️", a.Contract.String())
	is.Equal(ppp.North, a.Declarer)
	assertCalls(t, a.Bids, []call{
		{ppp.North, "1❤️"},  // North opens
		{ppp.East, "Pass"},  // East: 1♥ ≤ contract → coerced
		{ppp.South, "Pass"}, // South: same
		{ppp.West, "Pass"},  // West: same
	})
}

// call is an expected call in an auction: the seat that made it and the bid.
type call struct {
	dir ppp.Direction
	bid string
}

// assertCalls checks that bids holds exactly the calls in want, in order.
func assertCalls(t *testing.T, bids []ppp.BidAction, want []call) {
	t.Helper()
	is := is.NewRelaxed(t)
	is.Equal(len(want), len(bids))
	for i := range min(len(want), len(bids)) {
		is.Equal(want[i].dir, bids[i].Direction)    // seat of call i
		is.Equal(want[i].bid, bids[i].Bid.String()) // bid of call i
	}
}

// bidStrings renders bids for comparison; NewBid returns a fresh pointer per call.
func bidStrings(bids []*ppp.Bid) []string {
	s := make([]string, len(bids))
	for i, b := range bids {
		s[i] = b.String()
	}
	return s
}

func historyCases() []struct {
	name    string
	history ppp.History
} {
	pass := ppp.Pass
	return []struct {
		name    string
		history ppp.History
	}{
		{name: "empty", history: nil},
		{name: "passed out", history: ppp.History{
			{Direction: ppp.North, Bid: pass},
			{Direction: ppp.East, Bid: pass},
			{Direction: ppp.South, Bid: pass},
			{Direction: ppp.West, Bid: pass},
		}},
		{name: "leading pass", history: ppp.History{
			{Direction: ppp.North, Bid: pass},
			{Direction: ppp.East, Bid: ppp.NewBid(1, ppp.Hearts)},
			{Direction: ppp.South, Bid: pass},
			{Direction: ppp.West, Bid: ppp.NewBid(2, ppp.Hearts)},
		}},
	}
}

func TestHistory_BidsBy(t *testing.T) {
	want := map[string][]string{
		"empty":        {},
		"passed out":   {"Pass"},
		"leading pass": {"1❤️"},
	}
	for _, tc := range historyCases() {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(want[tc.name], bidStrings(tc.history.BidsBy(ppp.East)))
		})
	}
}

func TestHistory_LastBidBy(t *testing.T) {
	want := map[string]struct{ north, east string }{
		"empty":        {"Pass", "Pass"}, // a seat that has not called counts as Pass
		"passed out":   {"Pass", "Pass"},
		"leading pass": {"Pass", "1❤️"},
	}
	for _, tc := range historyCases() {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(want[tc.name].north, tc.history.LastBidBy(ppp.North).String())
			is.Equal(want[tc.name].east, tc.history.LastBidBy(ppp.East).String())
		})
	}
}

func TestHistory_Opening(t *testing.T) {
	type result struct {
		ok  bool
		dir ppp.Direction
		bid string
	}
	want := map[string]result{
		"empty":        {ok: false},
		"passed out":   {ok: false},
		"leading pass": {ok: true, dir: ppp.East, bid: "1❤️"},
	}
	for _, tc := range historyCases() {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			opening, ok := tc.history.Opening()
			is.Equal(want[tc.name].ok, ok)
			if ok {
				is.Equal(want[tc.name].dir, opening.Direction)
				is.Equal(want[tc.name].bid, opening.Bid.String())
			}
		})
	}
}

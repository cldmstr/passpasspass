package passpasspass_test

import (
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
	is.Equal(4, len(a.Bids))
	is.Equal("1NT", a.Bids[0].Bid.String()) // North opens
	is.Equal(ppp.Pass, a.Bids[1].Bid)       // East: 0 HCP
	is.Equal(ppp.Pass, a.Bids[2].Bid)       // South: 2 HCP
	is.Equal(ppp.Pass, a.Bids[3].Bid)       // West: 22 HCP > ceiling
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

	is.Equal(16, eval.MinCombinedHCP)                                     // 10 + 6
	is.Equal(20, eval.MaxCombinedHCP)                                     // 10 + 10
	is.Equal(ppp.SuitRange{Min: 6, Max: 9}, eval.SuitFit[ppp.Clubs])     // 4 + [2,5]
	is.Equal(ppp.SuitRange{Min: 5, Max: 7}, eval.SuitFit[ppp.Diamonds])  // 3 + [2,4]
	is.Equal(ppp.SuitRange{Min: 5, Max: 7}, eval.SuitFit[ppp.Hearts])    // 3 + [2,4]
	is.Equal(ppp.SuitRange{Min: 5, Max: 7}, eval.SuitFit[ppp.Spades])    // 3 + [2,4]
}

func TestPartnershipEval_BestMajorFit(t *testing.T) {
	tt := []struct {
		name   string
		hearts ppp.SuitRange
		spades ppp.SuitRange
		want   ppp.SuitRange
	}{
		{
			name:   "spades have higher max",
			hearts: ppp.SuitRange{Min: 7, Max: 9},
			spades: ppp.SuitRange{Min: 8, Max: 11},
			want:   ppp.SuitRange{Min: 8, Max: 11},
		},
		{
			name:   "hearts have higher max",
			hearts: ppp.SuitRange{Min: 9, Max: 11},
			spades: ppp.SuitRange{Min: 7, Max: 9},
			want:   ppp.SuitRange{Min: 9, Max: 11},
		},
		{
			name:   "equal max favors spades",
			hearts: ppp.SuitRange{Min: 7, Max: 9},
			spades: ppp.SuitRange{Min: 6, Max: 9},
			want:   ppp.SuitRange{Min: 6, Max: 9},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			eval := ppp.PartnershipEval{
				SuitFit: [4]ppp.SuitRange{{}, {}, tc.hearts, tc.spades},
			}
			is.Equal(tc.want, eval.BestMajorFit())
		})
	}
}

func TestPartnershipEval_BestMinorFit(t *testing.T) {
	tt := []struct {
		name     string
		clubs    ppp.SuitRange
		diamonds ppp.SuitRange
		want     ppp.SuitRange
	}{
		{
			name:     "diamonds have higher max",
			clubs:    ppp.SuitRange{Min: 6, Max: 8},
			diamonds: ppp.SuitRange{Min: 7, Max: 10},
			want:     ppp.SuitRange{Min: 7, Max: 10},
		},
		{
			name:     "clubs have higher max",
			clubs:    ppp.SuitRange{Min: 8, Max: 10},
			diamonds: ppp.SuitRange{Min: 6, Max: 8},
			want:     ppp.SuitRange{Min: 8, Max: 10},
		},
		{
			name:     "equal max favors diamonds",
			clubs:    ppp.SuitRange{Min: 7, Max: 9},
			diamonds: ppp.SuitRange{Min: 6, Max: 9},
			want:     ppp.SuitRange{Min: 6, Max: 9},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			eval := ppp.PartnershipEval{
				SuitFit: [4]ppp.SuitRange{tc.clubs, tc.diamonds, {}, {}},
			}
			is.Equal(tc.want, eval.BestMinorFit())
		})
	}
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
	is.Equal(4, len(a.Bids))
	is.Equal("1❤️", a.Bids[0].Bid.String()) // North opens
	is.Equal(ppp.Pass, a.Bids[1].Bid)        // East: 1♥ ≤ contract → coerced
	is.Equal(ppp.Pass, a.Bids[2].Bid)        // South: same
	is.Equal(ppp.Pass, a.Bids[3].Bid)        // West: same
}

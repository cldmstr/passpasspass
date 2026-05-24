package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

func TestRuleOfTwentyCondition(t *testing.T) {
	sut := ppp.NewConvention(
		ppp.NewBid(1, ppp.Clubs),
		ppp.NewRuleOfTwentyCondition(),
	)

	tt := []struct {
		name      string
		hand      []string
		shouldBid bool
	}{
		{
			// 10 HCP (A K Q J clubs) + 6 clubs + 4 diamonds = 20
			name:      "exactly 20",
			hand:      []string{"A♣️", "K♣️", "Q♣️", "J♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "2♠️", "3♠️"},
			shouldBid: true,
		},
		{
			// 10 HCP (A K Q J clubs) + 5 clubs + 4 diamonds = 19
			name:      "one short",
			hand:      []string{"A♣️", "K♣️", "Q♣️", "J♣️", "2♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "2♠️", "3♠️"},
			shouldBid: false,
		},
		{
			// 16 HCP + 7 spades + 4 hearts = 27
			name:      "well above 20",
			hand:      []string{"2♣️", "2♦️", "A❤️", "K❤️", "2❤️", "3❤️", "A♠️", "K♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "5♠️"},
			shouldBid: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			hand, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)
			willBid, _ := sut.NextBid(&ppp.AuctionContext{SelfHand: hand})
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestBalancedDistributionCondition(t *testing.T) {
	sut := ppp.NewConvention(
		ppp.NewBid(1, ppp.NoTrump),
		ppp.NewBalancedDistributionCondition(),
	)

	tt := []struct {
		name string
		hand []string

		shouldBid bool
		wantBid   string
	}{
		{
			name: "4-3-3-3 balanced",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"},

			shouldBid: true,
			wantBid:   "1NT",
		},
		{
			name: "4-4-3-2 balanced",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},

			shouldBid: true,
			wantBid:   "1NT",
		},
		{
			name: "5-3-3-2 balanced",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},

			shouldBid: true,
			wantBid:   "1NT",
		},
		{
			name: "4-4-4-1 singleton unbalanced",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "4❤️", "5❤️", "2♠️"},

			shouldBid: false,
			wantBid:   "Pass",
		},
		{
			name: "6-3-2-2 long suit unbalanced",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️", "2♠️", "3♠️"},

			shouldBid: false,
			wantBid:   "Pass",
		},
		{
			name: "5-4-2-2 two doubletons unbalanced",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "2♠️", "3♠️"},

			shouldBid: false,
			wantBid:   "Pass",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			hand, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)

			willBid, bid := sut.NextBid(&ppp.AuctionContext{SelfHand: hand})
			is.Equal(tc.shouldBid, willBid)
			is.Equal(tc.wantBid, bid.String())
		})
	}
}

func partnershipCtx(hand *ppp.Hand, eval ppp.PartnershipEval) *ppp.AuctionContext {
	return &ppp.AuctionContext{SelfHand: hand, Partnership: eval}
}

func TestNewPartnershipHCPCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	conv := ppp.NewConvention(ppp.NewBid(3, ppp.NoTrump), ppp.NewPartnershipHCPCondition(25, 28))

	tt := []struct {
		name        string
		minCombined int
		maxCombined int
		shouldBid   bool
	}{
		{name: "combined HCP in range", minCombined: 25, maxCombined: 28, shouldBid: true},
		{name: "min too low", minCombined: 23, maxCombined: 28, shouldBid: false},
		{name: "max too high", minCombined: 25, maxCombined: 30, shouldBid: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			eval := ppp.PartnershipEval{MinCombinedHCP: tc.minCombined, MaxCombinedHCP: tc.maxCombined}
			willBid, _ := conv.NextBid(partnershipCtx(hand, eval))
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestNewMajorFitCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	conv := ppp.NewConvention(ppp.NewBid(4, ppp.Spades), ppp.NewMajorFitCondition(8, 10))

	tt := []struct {
		name      string
		hearts    ppp.SuitRange
		spades    ppp.SuitRange
		shouldBid bool
	}{
		{
			name:      "spade fit in range",
			hearts:    ppp.SuitRange{Min: 6, Max: 8},
			spades:    ppp.SuitRange{Min: 8, Max: 10},
			shouldBid: true,
		},
		{
			name:      "fit max too high",
			hearts:    ppp.SuitRange{Min: 6, Max: 8},
			spades:    ppp.SuitRange{Min: 8, Max: 11},
			shouldBid: false,
		},
		{
			name:      "fit min too low",
			hearts:    ppp.SuitRange{Min: 6, Max: 8},
			spades:    ppp.SuitRange{Min: 6, Max: 9},
			shouldBid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			eval := ppp.PartnershipEval{
				SuitFit: [4]ppp.SuitRange{{}, {}, tc.hearts, tc.spades},
			}
			willBid, _ := conv.NextBid(partnershipCtx(hand, eval))
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestNewMinorFitCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	conv := ppp.NewConvention(ppp.NewBid(5, ppp.Diamonds), ppp.NewMinorFitCondition(9, 11))

	tt := []struct {
		name      string
		clubs     ppp.SuitRange
		diamonds  ppp.SuitRange
		shouldBid bool
	}{
		{
			name:      "diamond fit in range",
			clubs:     ppp.SuitRange{Min: 7, Max: 9},
			diamonds:  ppp.SuitRange{Min: 9, Max: 11},
			shouldBid: true,
		},
		{
			name:      "fit max too high",
			clubs:     ppp.SuitRange{Min: 7, Max: 9},
			diamonds:  ppp.SuitRange{Min: 9, Max: 12},
			shouldBid: false,
		},
		{
			name:      "fit min too low",
			clubs:     ppp.SuitRange{Min: 7, Max: 9},
			diamonds:  ppp.SuitRange{Min: 7, Max: 10},
			shouldBid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			eval := ppp.PartnershipEval{
				SuitFit: [4]ppp.SuitRange{tc.clubs, tc.diamonds, {}, {}},
			}
			willBid, _ := conv.NextBid(partnershipCtx(hand, eval))
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestNewHeartFitCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	conv := ppp.NewConvention(ppp.NewBid(4, ppp.Hearts), ppp.NewHeartFitCondition(8, 10))

	is := is.New(t)

	good := ppp.PartnershipEval{SuitFit: [4]ppp.SuitRange{{}, {}, {Min: 8, Max: 10}, {}}}
	willBid, _ := conv.NextBid(partnershipCtx(hand, good))
	is.True(willBid)

	tooHigh := ppp.PartnershipEval{SuitFit: [4]ppp.SuitRange{{}, {}, {Min: 8, Max: 11}, {}}}
	willBid, _ = conv.NextBid(partnershipCtx(hand, tooHigh))
	is.True(!willBid)
}

func TestNewSpadeFitCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	conv := ppp.NewConvention(ppp.NewBid(4, ppp.Spades), ppp.NewSpadeFitCondition(8, 10))

	is := is.New(t)

	good := ppp.PartnershipEval{SuitFit: [4]ppp.SuitRange{{}, {}, {}, {Min: 8, Max: 10}}}
	willBid, _ := conv.NextBid(partnershipCtx(hand, good))
	is.True(willBid)

	tooLow := ppp.PartnershipEval{SuitFit: [4]ppp.SuitRange{{}, {}, {}, {Min: 6, Max: 9}}}
	willBid, _ = conv.NextBid(partnershipCtx(hand, tooLow))
	is.True(!willBid)
}

func TestNewNoPriorBidsCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	conv := ppp.NewConvention(ppp.NewBid(1, ppp.NoTrump), ppp.NewNoPriorBidsCondition())

	tt := []struct {
		name      string
		history   []ppp.BidAction
		shouldBid bool
	}{
		{
			name:      "empty history",
			history:   nil,
			shouldBid: true,
		},
		{
			name:      "only passes before",
			history:   []ppp.BidAction{{Direction: ppp.West, Bid: ppp.Pass}},
			shouldBid: true,
		},
		{
			name:      "real bid in history",
			history:   []ppp.BidAction{{Direction: ppp.West, Bid: ppp.NewBid(1, ppp.Clubs)}},
			shouldBid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			willBid, _ := conv.NextBid(&ppp.AuctionContext{
				SelfHand: hand,
				Position: ppp.North,
				History:  tc.history,
			})
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestNewPartnerBidCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	// Jacoby transfer — fires only after partner opens 1NT. South's partner is North.
	conv := ppp.NewConvention(ppp.NewBid(2, ppp.Hearts), ppp.NewPartnerBidCondition("1NT"))

	tt := []struct {
		name      string
		history   []ppp.BidAction
		shouldBid bool
	}{
		{
			name:      "partner bid 1NT",
			history:   []ppp.BidAction{{Direction: ppp.North, Bid: ppp.NewBid(1, ppp.NoTrump)}},
			shouldBid: true,
		},
		{
			name:      "partner bid something else",
			history:   []ppp.BidAction{{Direction: ppp.North, Bid: ppp.NewBid(1, ppp.Clubs)}},
			shouldBid: false,
		},
		{
			name:      "partner has not bid",
			history:   nil,
			shouldBid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			willBid, _ := conv.NextBid(&ppp.AuctionContext{
				SelfHand: hand,
				Position: ppp.South, // North is partner
				History:  tc.history,
			})
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestNewRHOBidCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	// South's RHO is East: rhoOf(South=2) = (2+3)%4 = 1 = East.
	conv := ppp.NewConvention(ppp.NewBid(1, ppp.NoTrump), ppp.NewRHOBidCondition("1C"))

	tt := []struct {
		name      string
		history   []ppp.BidAction
		shouldBid bool
	}{
		{
			name:      "RHO bid 1C",
			history:   []ppp.BidAction{{Direction: ppp.East, Bid: ppp.NewBid(1, ppp.Clubs)}},
			shouldBid: true,
		},
		{
			name:      "RHO bid something else",
			history:   []ppp.BidAction{{Direction: ppp.East, Bid: ppp.NewBid(1, ppp.Hearts)}},
			shouldBid: false,
		},
		{
			name:      "RHO has not bid",
			history:   nil,
			shouldBid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			willBid, _ := conv.NextBid(&ppp.AuctionContext{
				SelfHand: hand,
				Position: ppp.South,
				History:  tc.history,
			})
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

func TestNewLHOBidCondition(t *testing.T) {
	hand := mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"})
	// South's LHO is West: lhoOf(South=2) = (2+1)%4 = 3 = West.
	conv := ppp.NewConvention(ppp.NewBid(1, ppp.NoTrump), ppp.NewLHOBidCondition("1C"))

	tt := []struct {
		name      string
		history   []ppp.BidAction
		shouldBid bool
	}{
		{
			name:      "LHO bid 1C",
			history:   []ppp.BidAction{{Direction: ppp.West, Bid: ppp.NewBid(1, ppp.Clubs)}},
			shouldBid: true,
		},
		{
			name:      "LHO bid something else",
			history:   []ppp.BidAction{{Direction: ppp.West, Bid: ppp.NewBid(1, ppp.Hearts)}},
			shouldBid: false,
		},
		{
			name:      "LHO has not bid",
			history:   nil,
			shouldBid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			willBid, _ := conv.NextBid(&ppp.AuctionContext{
				SelfHand: hand,
				Position: ppp.South,
				History:  tc.history,
			})
			is.Equal(tc.shouldBid, willBid)
		})
	}
}

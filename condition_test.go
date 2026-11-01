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

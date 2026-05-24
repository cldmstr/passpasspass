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

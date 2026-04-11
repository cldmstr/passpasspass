package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

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

			willBid, bid := sut.NextBid(hand)
			is.Equal(tc.shouldBid, willBid)
			is.Equal(tc.wantBid, bid.String())
		})
	}
}

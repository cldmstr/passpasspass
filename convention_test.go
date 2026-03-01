package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

func TestConvention_NextBid(t *testing.T) {
	sut := ppp.NewConvention(
		ppp.NewBid(1, ppp.Clubs),
		ppp.NewHCPIntervalCondition(11, 17),
		ppp.NewMinMaxDistributionCondition(ppp.NewDistribution(2, 1, 1, 1), ppp.NewDistribution(7, 4, 4, 4)),
	)

	tt := []struct {
		name string
		hand []string

		shouldBid bool
		wantBid   string
	}{
		{
			name: "bid 2 clubs",
			hand: []string{"2♣️", "K♣️", "A♣️", "3❤️", "Q❤️", "K❤️", "6♠️", "7♠️", "9♠️", "J♦️", "4♦️", "5♦️", "7♦️"},

			shouldBid: true,
			wantBid:   "1♣️",
		},
		{
			name: "low points",
			hand: []string{"2♣️", "5♣️", "A♣️", "3❤️", "9❤️", "4❤️", "6♠️", "7♠️", "9♠️", "J♦️", "4♦️", "5♦️", "7♦️"},

			shouldBid: false,
			wantBid:   "",
		},
		{
			name: "high points",
			hand: []string{"K♣️", "Q♣️", "A♣️", "Q❤️", "K❤️", "A❤️", "6♠️", "7♠️", "9♠️", "J♦️", "4♦️", "5♦️", "7♦️"},

			shouldBid: false,
			wantBid:   "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			hand, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)

			willBid, bid := sut.NextBid(hand)
			is.Equal(tc.shouldBid, willBid)
			if tc.shouldBid {
				is.Equal(tc.wantBid, bid.String())
			}
		})
	}
}

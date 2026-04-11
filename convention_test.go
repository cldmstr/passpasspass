package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

func TestBidOutput(t *testing.T) {
	tt := []struct {
		name string
		bid  *ppp.Bid

		want string
	}{
		{
			name: "clubs",
			bid:  ppp.NewBid(3, ppp.Clubs),

			want: "3♣️",
		},
		{
			name: "diamonds",
			bid:  ppp.NewBid(1, ppp.Diamonds),

			want: "1♦️",
		},
		{
			name: "hearts",
			bid:  ppp.NewBid(3, ppp.Hearts),

			want: "3❤️",
		},
		{
			name: "spades",
			bid:  ppp.NewBid(7, ppp.Spades),

			want: "7♠️",
		},
		{
			name: "pass",
			bid:  ppp.Pass,

			want: "Pass",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			is.Equal(tc.want, tc.bid.String())
		})
	}
}

func TestConvention_NextBid(t *testing.T) {
	sut := ppp.NewConvention(
		ppp.NewBid(1, ppp.Clubs),
		ppp.NewHCPIntervalCondition(11, 17),
		ppp.NewMinMaxDistributionCondition(ppp.DistributionRange{{2, 7}, {1, 4}, {1, 4}, {1, 4}}),
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
			wantBid:   "Pass",
		},
		{
			name: "high points",
			hand: []string{"K♣️", "Q♣️", "A♣️", "Q❤️", "K❤️", "A❤️", "6♠️", "7♠️", "9♠️", "J♦️", "4♦️", "5♦️", "7♦️"},

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

package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

func TestNewBidFromString(t *testing.T) {
	tt := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "1C", input: "1C", want: "1♣️"},
		{name: "1D", input: "1D", want: "1♦️"},
		{name: "1H", input: "1H", want: "1❤️"},
		{name: "1S", input: "1S", want: "1♠️"},
		{name: "1NT", input: "1NT", want: "1NT"},
		{name: "2C", input: "2C", want: "2♣️"},
		{name: "7NT", input: "7NT", want: "7NT"},
		{name: "Pass", input: "Pass", want: "Pass"},
		{name: "pass", input: "pass", want: "Pass"},
		{name: "invalid suit", input: "1X", wantErr: true},
		{name: "invalid level 0", input: "0C", wantErr: true},
		{name: "invalid level 8", input: "8H", wantErr: true},
		{name: "too short", input: "1", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			bid, err := ppp.NewBidFromString(tc.input)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(tc.want, bid.String())
		})
	}
}

func TestNewBidFromString_PassIsSingleton(t *testing.T) {
	is := is.New(t)
	bid, err := ppp.NewBidFromString("Pass")
	is.NoErr(err)
	is.True(bid == ppp.Pass)
}

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
			hand: []string{"2♣️", "K♣️", "A♣️", "J♦️", "4♦️", "5♦️", "7♦️", "3❤️", "Q❤️", "K❤️", "6♠️", "7♠️", "9♠️"},

			shouldBid: true,
			wantBid:   "1♣️",
		},
		{
			name: "low points",
			hand: []string{"2♣️", "5♣️", "A♣️", "J♦️", "4♦️", "5♦️", "7♦️", "3❤️", "9❤️", "4❤️", "6♠️", "7♠️", "9♠️"},

			shouldBid: false,
			wantBid:   "Pass",
		},
		{
			name: "high points",
			hand: []string{"K♣️", "Q♣️", "A♣️", "J♦️", "4♦️", "5♦️", "7♦️", "Q❤️", "K❤️", "A❤️", "6♠️", "7♠️", "9♠️"},

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

package passpasspass_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

func TestHandDistributionHappy(t *testing.T) {
	tt := []struct {
		name string
		hand []string

		wantDist ppp.Distribution
		wantHCP  int
	}{
		{
			name: "all hearts",
			hand: []string{"2❤️", "3❤️", "4❤️", "5❤️", "6❤️", "7❤️", "8❤️", "9❤️", "10❤️", "J❤️", "Q❤️", "K❤️", "A❤️"},

			wantDist: ppp.Distribution([]int{0, 0, 13, 0}),
			wantHCP:  10,
		},
		{
			name: "all spades",
			hand: []string{"2♠️", "3♠️", "4♠️", "5♠️", "6♠️", "7♠️", "8♠️", "9♠️", "10♠️", "J♠️", "Q♠️", "K♠️", "A♠️"},

			wantDist: ppp.Distribution([]int{0, 0, 0, 13}),
			wantHCP:  10,
		},
		{
			name: "all diamonds",
			hand: []string{"2♦️", "3♦️", "4♦️", "5♦️", "6♦️", "7♦️", "8♦️", "9♦️", "10♦️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantDist: ppp.Distribution([]int{0, 13, 0, 0}),
			wantHCP:  10,
		},
		{
			name: "all clubs",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"},

			wantDist: ppp.Distribution([]int{13, 0, 0, 0}),
			wantHCP:  10,
		},
		{
			name: "normal flat",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "6❤️", "7❤️", "8❤️", "9♠️", "10♠️", "J♠️", "Q♦️", "K♦️", "A♦️"},

			wantDist: ppp.Distribution([]int{4, 3, 3, 3}),
			wantHCP:  10,
		},
		{
			name: "max HCP",
			hand: []string{"Q♣️", "K♣️", "A♣️", "Q❤️", "K❤️", "A❤️", "Q♠️", "K♠️", "A♠️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantDist: ppp.Distribution([]int{3, 4, 3, 3}),
			wantHCP:  37,
		},
		{
			name: "missing clubs",
			hand: []string{"2♦️", "3♦️", "4♦️", "5♦️", "6♦️", "7♠️", "8♠️", "9♠️", "10❤️", "J❤️", "Q❤️", "K❤️", "A❤️"},

			wantDist: ppp.Distribution([]int{0, 5, 5, 3}),
			wantHCP:  10,
		},
		{
			name: "missing diamonds",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "6❤️", "7♠️", "8♠️", "9♠️", "10❤️", "J❤️", "Q❤️", "K❤️", "A❤️"},

			wantDist: ppp.Distribution([]int{4, 0, 6, 3}),
			wantHCP:  10,
		},
		{
			name: "missing hearts",
			hand: []string{"2♣️", "3♣️", "4♣️", "Q♣️", "K♣️", "A♣️", "Q♠️", "K♠️", "A♠️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantDist: ppp.Distribution([]int{6, 4, 0, 3}),
			wantHCP:  28,
		},
		{
			name: "missing spades",
			hand: []string{"Q♣️", "K♣️", "A♣️", "Q❤️", "6❤️", "7❤️", "8❤️", "K❤️", "A❤️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantDist: ppp.Distribution([]int{3, 4, 6, 0}),
			wantHCP:  28,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			sut, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)

			result := sut.Description()

			is.Equal(result.Distribution, tc.wantDist)
			is.Equal(result.HCP, tc.wantHCP)
		})
	}
}

func TestHandDistributionUnhappy(t *testing.T) {
	tt := []struct {
		name string
		hand []string

		wantError error
	}{
		{
			name: "invalid hand duplicate card",
			hand: []string{"Q♣️", "K♣️", "A♣️", "Q❤️", "6❤️", "7❤️", "8❤️", "K❤️", "A♣️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantError: fmt.Errorf("invalid hand: %w", fmt.Errorf("invalid hand contains duplicates: [A♣️]")),
		},
		{
			name: "invalid hand too many cards",
			hand: []string{"Q♣️", "K♣️", "A♣️", "3♣️", "4♣️", "Q❤️", "6❤️", "7❤️", "8❤️", "K❤️", "A♣️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantError: fmt.Errorf("invalid hand: %w", fmt.Errorf("invalid hand size should be 13: 15")),
		},
		{
			name: "invalid hand too few cards",
			hand: []string{"Q♣️", "K♣️", "K❤️", "A♣️", "J♦️", "Q♦️", "K♦️", "A♦️"},

			wantError: fmt.Errorf("invalid hand: %w", fmt.Errorf("invalid hand size should be 13: 8")),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			_, err := ppp.NewHandFromStrings(tc.hand)

			is.Equal(err, tc.wantError)
		})
	}
}

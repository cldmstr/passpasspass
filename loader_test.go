package passpasspass_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// writeTempFile writes content to a temp file and returns its path.
func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoader_LoadCard_HCPCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 1NT
    conditions:
      - hcp: [15, 17]
`)
	cs, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.NoErr(err)

	inRange, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("1NT", cs.Bid(inRange).String())

	tooLow, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("Pass", cs.Bid(tooLow).String())
}

func TestLoader_LoadCard_BalancedCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 1NT
    conditions:
      all:
        - balanced: true
        - hcp: [15, 17]
`)
	cs, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.NoErr(err)

	balanced, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("1NT", cs.Bid(balanced).String())

	unbalanced, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "J♣️", "10♣️", "9♣️", "A♦️", "K♦️", "Q♦️", "2❤️", "3❤️", "2♠️", "3♠️",
	})
	is.Equal("Pass", cs.Bid(unbalanced).String())
}

func TestLoader_LoadCard_SuitLengthCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 1H
    conditions:
      all:
        - hcp: [11, 17]
        - suit_length:
            hearts: [5, 13]
`)
	cs, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.NoErr(err)

	fiveHearts, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "A❤️", "K❤️", "Q❤️", "J❤️", "10❤️", "2♠️",
	})
	is.Equal("1❤️", cs.Bid(fiveHearts).String())

	fourHearts, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "5♦️", "A❤️", "K❤️", "Q❤️", "J❤️", "2♠️",
	})
	is.Equal("Pass", cs.Bid(fourHearts).String())
}

func TestLoader_LoadCard_AnyCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 2D
    conditions:
      any:
        - hcp: [6, 9]
        - hcp: [23, 24]
`)
	cs, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.NoErr(err)

	// 6 HCP — matches hcp: [6, 9]
	weak, _ := ppp.NewHandFromStrings([]string{
		"J♣️", "2♣️", "3♣️", "4♣️", "Q♦️", "2♦️", "3♦️", "4♦️", "K❤️", "2❤️", "3❤️", "2♠️", "3♠️",
	})
	is.Equal("2♦️", cs.Bid(weak).String())

	// 16 HCP — matches neither [6,9] nor [23,24]
	middle, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("Pass", cs.Bid(middle).String())
}

func TestLoader_LoadCard_Includes(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "strong_notrump.yaml", `
name: strong_notrump
conventions:
  - bid: 1NT
    conditions:
      all:
        - balanced: true
        - hcp: [15, 17]
`)
	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
includes:
  - strong_notrump
conventions:
  - bid: 1C
    conditions:
      - hcp: [11, 17]
`)
	cs, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.NoErr(err)

	// Balanced 15 HCP — 1NT (from include) wins over 1C (card convention)
	hand, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("1NT", cs.Bid(hand).String())
}

func TestLoader_LoadCard_ErrorOnMissingInclude(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
includes:
  - does_not_exist
`)
	_, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.True(err != nil)
}

func TestLoader_LoadCard_ErrorOnInvalidBid(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	cardPath := writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 9NT
    conditions:
      - hcp: [11, 17]
`)
	_, err := ppp.NewLoader(dir).LoadCard(cardPath)
	is.True(err != nil)
}

func TestLoader_RoundTrip_TwoOverOne(t *testing.T) {
	is := is.New(t)

	cs, err := ppp.NewLoader("conventions").LoadCard("cards/two_over_one.yaml")
	is.NoErr(err)

	tt := []struct {
		name string
		hand []string
		want string
	}{
		{
			name: "strong 2D - game forcing",
			hand: []string{"A♣️", "K♣️", "Q♣️", "J♣️", "A♦️", "K♦️", "Q♦️", "J♦️", "A❤️", "K❤️", "Q❤️", "J❤️", "A♠️"},
			want: "2♦️",
		},
		{
			// 20 HCP, no earlier convention matches
			name: "strong 2C",
			hand: []string{"A♣️", "K♣️", "Q♣️", "J♣️", "A♦️", "K♦️", "Q♦️", "J♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},
			want: "2♣️",
		},
		{
			// 16 HCP, balanced 4-3-3-3
			name: "1NT balanced 15-17",
			hand: []string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},
			want: "1NT",
		},
		{
			// 14 HCP, 6 spades — triggers 1S before 1H/1D/1C
			name: "1S five card major",
			hand: []string{"A♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2❤️", "A♠️", "K♠️", "Q♠️", "J♠️", "10♠️", "9♠️"},
			want: "1♠️",
		},
		{
			// 11 HCP, 5 diamonds, no 5-card major
			name: "1D four diamonds",
			hand: []string{"J♣️", "2♣️", "3♣️", "A♦️", "K♦️", "Q♦️", "J♦️", "2♦️", "2❤️", "3❤️", "2♠️", "3♠️", "4♠️"},
			want: "1♦️",
		},
		{
			name: "pass - too few points",
			hand: []string{"2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},
			want: "Pass",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			hand, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)
			is.Equal(tc.want, cs.Bid(hand).String())
		})
	}
}

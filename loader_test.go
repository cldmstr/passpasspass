package passpasspass_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// writeTempFile writes content to the file name in dir.
func writeTempFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func ctx(hand *ppp.Hand) *ppp.AuctionContext {
	return &ppp.AuctionContext{SelfHand: hand}
}

func ctxAfterPartnerNT(hand *ppp.Hand) *ppp.AuctionContext {
	return &ppp.AuctionContext{
		SelfHand: hand,
		Position: ppp.South,
		History: []ppp.BidAction{
			{Direction: ppp.North, Bid: ppp.NewBid(1, ppp.NoTrump)},
		},
	}
}

func TestLoader_LoadCard_HCPCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 1NT
    conditions:
      - hcp: [15, 17]
`)
	cs, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.NoErr(err)

	inRange, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("1NT", cs.Bid(ctx(inRange)).String())

	tooLow, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "5♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("Pass", cs.Bid(ctx(tooLow)).String())
}

func TestLoader_LoadCard_BalancedCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 1NT
    conditions:
      all:
        - balanced: true
        - hcp: [15, 17]
`)
	cs, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.NoErr(err)

	balanced, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("1NT", cs.Bid(ctx(balanced)).String())

	unbalanced, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "J♣️", "10♣️", "9♣️", "A♦️", "K♦️", "Q♦️", "2❤️", "3❤️", "2♠️", "3♠️",
	})
	is.Equal("Pass", cs.Bid(ctx(unbalanced)).String())
}

func TestLoader_LoadCard_SuitLengthCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 1H
    conditions:
      all:
        - hcp: [11, 17]
        - suit_length:
            hearts: [5, 13]
`)
	cs, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.NoErr(err)

	fiveHearts, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "A❤️", "K❤️", "Q❤️", "J❤️", "10❤️", "2♠️",
	})
	is.Equal("1❤️", cs.Bid(ctx(fiveHearts)).String())

	fourHearts, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "5♦️", "A❤️", "K❤️", "Q❤️", "J❤️", "2♠️",
	})
	is.Equal("Pass", cs.Bid(ctx(fourHearts)).String())
}

func TestLoader_LoadCard_AnyCondition(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 2D
    conditions:
      any:
        - hcp: [6, 9]
        - hcp: [23, 24]
`)
	cs, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.NoErr(err)

	// 6 HCP — matches hcp: [6, 9]
	weak, _ := ppp.NewHandFromStrings([]string{
		"J♣️", "2♣️", "3♣️", "4♣️", "Q♦️", "2♦️", "3♦️", "4♦️", "K❤️", "2❤️", "3❤️", "2♠️", "3♠️",
	})
	is.Equal("2♦️", cs.Bid(ctx(weak)).String())

	// 16 HCP — matches neither [6,9] nor [23,24]
	middle, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("Pass", cs.Bid(ctx(middle)).String())
}

func TestLoader_LoadCard_Includes(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	convDir := filepath.Join(dir, "conventions")
	if err := os.MkdirAll(convDir, 0o700); err != nil {
		t.Fatal(err)
	}

	writeTempFile(t, convDir, "strong_notrump.yaml", `
name: strong_notrump
conventions:
  - bid: 1NT
    conditions:
      all:
        - balanced: true
        - hcp: [15, 17]
`)
	writeTempFile(t, dir, "card.yaml", `
name: test
includes:
  - strong_notrump
conventions:
  - bid: 1C
    conditions:
      - hcp: [11, 17]
`)
	cs, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.NoErr(err)

	// Balanced 15 HCP — 1NT (from include) wins over 1C (card convention)
	hand, _ := ppp.NewHandFromStrings([]string{
		"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️",
	})
	is.Equal("1NT", cs.Bid(ctx(hand)).String())
}

func TestLoader_LoadCard_ErrorOnMissingInclude(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "card.yaml", `
name: test
includes:
  - does_not_exist
`)
	_, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.True(err != nil)
}

func TestLoader_LoadCard_ErrorOnInvalidBid(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()

	writeTempFile(t, dir, "card.yaml", `
name: test
conventions:
  - bid: 9NT
    conditions:
      - hcp: [11, 17]
`)
	_, err := ppp.NewLoader(os.DirFS(dir)).LoadCard("card.yaml")
	is.True(err != nil)
}

func TestLoader_LoadCard_Errors(t *testing.T) {
	baronCard := "name: test\nincludes:\n  - baron\n"

	tt := []struct {
		name    string
		files   map[string]string
		wantErr string
	}{
		{
			name: "misspelt conditions",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1NT
    condtions:
      - hcp: [15, 17]
`},
			wantErr: `card.yaml:4: bid 1NT: condtions: unknown convention key (known keys: bid, conditions)`,
		},
		{
			name: "bid missing",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - conditions:
      - hcp: [15, 17]
`},
			wantErr: `card.yaml:3: convention has no bid`,
		},
		{
			name: "conditions missing",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1NT
`},
			wantErr: `card.yaml:3: bid 1NT: convention has no conditions`,
		},
		{
			name: "conditions without a value",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1NT
    conditions:
`},
			wantErr: `card.yaml:3: bid 1NT: convention has no conditions`,
		},
		{
			name: "empty conditions",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1NT
    conditions: []
`},
			wantErr: `card.yaml:4: bid 1NT: a condition list must not be empty`,
		},
		{
			name: "misspelt condition key",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 2S
    conditions:
      hcp: [0, 10]
      suitlength: {spades: [6, 13]}
`},
			wantErr: `card.yaml:6: bid 2S: suitlength: unknown condition key (known keys: all, any, auction, balanced, hcp, no_prior_bids, not, partner_bid, partnership_hcp, rule_of_20, suit_length)`,
		},
		{
			name: "invalid auction pattern element",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 2H
    conditions:
      auction: [1H, Foo]
`},
			wantErr: `card.yaml:5: bid 2H: auction: "Foo" is neither a call nor a wildcard (Any, Any*, Pass*)`,
		},
		{
			name: "too many alternatives",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1C
    conditions:
      - any: [{hcp: [0, 5]}, {hcp: [20, 25]}]
      - any: [{hcp: [0, 6]}, {hcp: [20, 26]}]
      - any: [{hcp: [0, 7]}, {hcp: [20, 27]}]
      - any: [{hcp: [0, 8]}, {hcp: [20, 28]}]
      - any: [{hcp: [0, 9]}, {hcp: [20, 29]}]
      - any: [{hcp: [0, 10]}, {hcp: [20, 30]}]
      - any: [{hcp: [0, 11]}, {hcp: [20, 31]}]
`},
			wantErr: `card.yaml:5: bid 1C: conditions can show more than 64 alternatives; split them into several conventions`,
		},
		{
			name: "key repeated in a convention",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1NT
    bid: 2NT
    conditions:
      - hcp: [15, 17]
`},
			wantErr: `card.yaml:4: bid 1NT: bid: key appears more than once`,
		},
		{
			name: "invalid bid",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 9Z
    conditions:
      - hcp: [15, 17]
`},
			wantErr: `card.yaml:3: bid 9Z: invalid bid level in "9Z"`,
		},
		{
			name: "convention is not a mapping",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - 1NT
`},
			wantErr: `card.yaml:3: a convention must be a mapping with bid and conditions`,
		},
		{
			name: "conditions are neither a mapping nor a list",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1NT
    conditions: balanced
`},
			wantErr: `card.yaml:4: bid 1NT: a condition must be a mapping or a list`,
		},
		{
			name: "unknown top-level key in a card",
			files: map[string]string{"card.yaml": `name: test
incldes:
  - baron
`},
			wantErr: `card.yaml:2: incldes: unknown card key (known keys: conventions, includes, name)`,
		},
		{
			name: "YAML syntax error",
			files: map[string]string{"card.yaml": `name: test
conventions:
  - bid: 1C
    conditions: [hcp: [11, 17]
`},
			wantErr: "card.yaml:3: did not find expected ',' or ']'",
		},
		{
			name:    "card is empty",
			files:   map[string]string{"card.yaml": ""},
			wantErr: `card.yaml: file is empty`,
		},
		{
			name:    "card is not a mapping",
			files:   map[string]string{"card.yaml": "- baron\n"},
			wantErr: `card.yaml:1: must be a mapping`,
		},
		{
			name:    "missing card",
			files:   map[string]string{},
			wantErr: `card.yaml: file does not exist`,
		},
		{
			name:    "missing sub-convention",
			files:   map[string]string{"card.yaml": baronCard},
			wantErr: `card.yaml: include baron: conventions/baron.yaml: file does not exist`,
		},
		{
			name: "unknown top-level key in a sub-convention",
			files: map[string]string{
				"card.yaml": baronCard,
				"conventions/baron.yaml": `name: baron
convention:
  - bid: 2S
`,
			},
			wantErr: `conventions/baron.yaml:2: convention: unknown sub-convention key (known keys: conventions, name)`,
		},
		{
			name: "invalid condition in a sub-convention",
			files: map[string]string{
				"card.yaml": baronCard,
				"conventions/baron.yaml": `name: baron
conventions:
  - bid: 2S
    conditions:
      all:
        - partner_bid: "1NT"
        - any:
            - suit_length: {hearts: [0, 14]}
`,
			},
			wantErr: `conventions/baron.yaml:8: bid 2S: all.any.suit_length.hearts: [0, 14] is outside 0 to 13`,
		},
		{
			name: "error inside a sub-convention's convention",
			files: map[string]string{
				"card.yaml": baronCard,
				"conventions/baron.yaml": `name: baron
conventions:
  - bid: 2S
    conditions:
      all:
        - partner_bid: "1NT"
        - hcp: [8, 40]
    forcing: true
`,
			},
			wantErr: `conventions/baron.yaml:8: bid 2S: forcing: unknown convention key (known keys: bid, conditions)`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			fsys := fstest.MapFS{}
			for name, content := range tc.files {
				fsys[name] = &fstest.MapFile{Data: []byte(content)}
			}
			_, err := ppp.NewLoader(fsys).LoadCard("card.yaml")
			is.True(err != nil)
			is.Equal(err.Error(), tc.wantErr)
		})
	}
}

// TestLoader_LoadCard_MaxAlternatives loads a convention at the cap: six `any` of two under one
// `all` can show 2⁶ = 64 alternatives.
func TestLoader_LoadCard_MaxAlternatives(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{"card.yaml": &fstest.MapFile{Data: []byte(`name: test
conventions:
  - bid: 1C
    conditions:
      - any: [{hcp: [0, 5]}, {hcp: [20, 25]}]
      - any: [{hcp: [0, 6]}, {hcp: [20, 26]}]
      - any: [{hcp: [0, 7]}, {hcp: [20, 27]}]
      - any: [{hcp: [0, 8]}, {hcp: [20, 28]}]
      - any: [{hcp: [0, 9]}, {hcp: [20, 29]}]
      - any: [{hcp: [0, 10]}, {hcp: [20, 30]}]
`)}}
	_, err := ppp.NewLoader(fsys).LoadCard("card.yaml")
	is.NoErr(err)
}

func TestLoader_RoundTrip_TwoOverOne(t *testing.T) {
	is := is.New(t)

	dataFS, err := fs.Sub(ppp.DefaultData, "data")
	is.NoErr(err)
	cs, err := ppp.NewLoader(dataFS).LoadCard("cards/two_over_one.yaml")
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
			// 11 HCP, 5 diamonds, 4 spades — rule of 20: 11+5+4=20
			name: "1D five card diamond",
			hand: []string{"J♣️", "2♣️", "K♦️", "Q♦️", "J♦️", "2♦️", "3♦️", "2❤️", "3❤️", "A♠️", "2♠️", "3♠️", "4♠️"},
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
			is.Equal(tc.want, cs.Bid(ctx(hand)).String())
		})
	}
}

func TestLoader_RoundTrip_MultiTwoDiamond_Weak(t *testing.T) {
	is := is.New(t)

	dataFS, err := fs.Sub(ppp.DefaultData, "data")
	is.NoErr(err)
	loader := ppp.NewLoader(dataFS)
	twoOverOne, err := loader.LoadCard("cards/two_over_one.yaml")
	is.NoErr(err)
	colin, err := loader.LoadCard("cards/colin.yaml")
	is.NoErr(err)

	tt := []struct {
		name string
		card *ppp.ConventionSheet
		hand []string
		want string
	}{
		{
			// A❤️(4)+Q❤️(2) = 6 HCP: lower bound
			name: "weak hearts - 6 HCP",
			card: twoOverOne,
			hand: []string{"A❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "5❤️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "2♠️"},
			want: "2♦️",
		},
		{
			// A❤️(4)+K❤️(3)+Q❤️(2)+J♣️(1) = 10 HCP: upper bound
			name: "weak hearts - 10 HCP",
			card: twoOverOne,
			hand: []string{"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "J♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2♠️"},
			want: "2♦️",
		},
		{
			// A♠️(4)+Q♠️(2) = 6 HCP: lower bound
			name: "weak spades - 6 HCP",
			card: twoOverOne,
			hand: []string{"A♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "5♠️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "2❤️"},
			want: "2♦️",
		},
		{
			// A♠️(4)+K♠️(3)+Q♠️(2)+J♣️(1) = 10 HCP: upper bound
			name: "weak spades - 10 HCP",
			card: twoOverOne,
			hand: []string{"A♠️", "K♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "J♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2❤️"},
			want: "2♦️",
		},
		{
			// A❤️(4)+K❤️(3)+Q❤️(2)+J❤️(1) = 10 HCP, 7 hearts
			name: "weak hearts - 7-card suit",
			card: twoOverOne,
			hand: []string{"A❤️", "K❤️", "Q❤️", "J❤️", "2❤️", "3❤️", "4❤️", "2♣️", "3♣️", "2♦️", "3♦️", "2♠️", "3♠️"},
			want: "2♦️",
		},
		{
			// A❤️(4)+J❤️(1) = 5 HCP: below lower bound
			name: "pass - 6 hearts 5 HCP",
			card: twoOverOne,
			hand: []string{"A❤️", "J❤️", "2❤️", "3❤️", "4❤️", "5❤️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "2♠️"},
			want: "Pass",
		},
		{
			// A♠️(4)+J♠️(1) = 5 HCP: below lower bound
			name: "pass - 6 spades 5 HCP",
			card: twoOverOne,
			hand: []string{"A♠️", "J♠️", "2♠️", "3♠️", "4♠️", "5♠️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "2❤️"},
			want: "Pass",
		},
		{
			// 10 HCP but only 5 hearts: suit too short for weak 2D, too weak for 1H
			name: "pass - 5 hearts 10 HCP",
			card: twoOverOne,
			hand: []string{"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "J♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2♠️", "3♠️"},
			want: "Pass",
		},
		{
			// 6 diamonds, 10 HCP: minor suits don't qualify for weak Multi
			name: "pass - 6 diamonds 10 HCP",
			card: twoOverOne,
			hand: []string{"A♦️", "K♦️", "Q♦️", "2♦️", "3♦️", "4♦️", "J♣️", "2♣️", "3♣️", "2❤️", "3❤️", "4❤️", "2♠️"},
			want: "Pass",
		},
		{
			// A❤️(4)+K❤️(3)+Q❤️(2)+Q♣️(2) = 11 HCP: just above weak range
			// rule of 20: 11+6+3 = 20
			name: "1H - 6 hearts 11 HCP",
			card: twoOverOne,
			hand: []string{"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "Q♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2♠️"},
			want: "1❤️",
		},
		{
			// A♠️(4)+K♠️(3)+Q♠️(2)+Q♣️(2) = 11 HCP: just above weak range
			// rule of 20: 11+6+3 = 20
			name: "1S - 6 spades 11 HCP",
			card: twoOverOne,
			hand: []string{"A♠️", "K♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "Q♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2❤️"},
			want: "1♠️",
		},
		{
			// 10 HCP, 6 hearts + 5 clubs: below Tartan's 11 HCP, Multi fires
			name: "colin weak hearts - 10 HCP two-suiter",
			card: colin,
			hand: []string{"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "J♣️", "2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "2♠️"},
			want: "2♦️",
		},
		{
			// 11 HCP, 6 hearts + 5 clubs: above Multi range, Tartan fires
			name: "colin tartan 2H - 11 HCP two-suiter",
			card: colin,
			hand: []string{"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "Q♣️", "2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "2♠️"},
			want: "2❤️",
		},
		{
			// 10 HCP, 6 spades + 5 diamonds: below Tartan's 11 HCP, Multi fires
			name: "colin weak spades - 10 HCP two-suiter",
			card: colin,
			hand: []string{"A♠️", "K♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "J♦️", "2♦️", "3♦️", "4♦️", "5♦️", "2♣️", "2❤️"},
			want: "2♦️",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			hand, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)
			is.Equal(tc.want, tc.card.Bid(ctx(hand)).String())
		})
	}

	t.Run("not an opening - partner already bid", func(t *testing.T) {
		is := is.New(t)
		// 10 HCP, 6 hearts: weak 2D hand, but no_prior_bids blocks it
		hand, err := ppp.NewHandFromStrings([]string{"A❤️", "K❤️", "Q❤️", "2❤️", "3❤️", "4❤️", "J♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2♠️"})
		is.NoErr(err)
		is.Equal("Pass", twoOverOne.Bid(ctxAfterPartnerNT(hand)).String())
	})
}

func TestLoader_RoundTrip_Colin(t *testing.T) {

	is := is.New(t)

	dataFS, err := fs.Sub(ppp.DefaultData, "data")
	is.NoErr(err)
	cs, err := ppp.NewLoader(dataFS).LoadCard("cards/colin.yaml")
	is.NoErr(err)

	tt := []struct {
		name string
		hand []string
		want string
	}{
		{
			name: "tartan 2H - 5+ hearts and 5+ clubs",
			hand: []string{"K❤️", "Q❤️", "J❤️", "2❤️", "3❤️", "A♣️", "K♣️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "2♠️"},
			want: "2❤️",
		},
		{
			name: "tartan 2S - 5+ spades and 5+ diamonds",
			hand: []string{"K♠️", "Q♠️", "J♠️", "2♠️", "3♠️", "A♦️", "K♦️", "2♦️", "3♦️", "4♦️", "2♣️", "3♣️", "2❤️"},
			want: "2♠️",
		},
		{
			name: "1H not tartan - 5 hearts but no 5-card side suit",
			hand: []string{"A❤️", "K❤️", "2❤️", "3❤️", "4❤️", "Q♣️", "J♣️", "2♣️", "3♣️", "Q♦️", "2♦️", "3♦️", "2♠️"},
			want: "1❤️",
		},
		{
			name: "2NT - 5-5 minors 10 HCP",
			hand: []string{"K♣️", "Q♣️", "2♣️", "3♣️", "4♣️", "A♦️", "J♦️", "2♦️", "3♦️", "4♦️", "2❤️", "2♠️", "3♠️"},
			want: "2NT",
		},
		{
			name: "1NT - balanced 16 HCP",
			hand: []string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},
			want: "1NT",
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
			is.Equal(tc.want, cs.Bid(ctx(hand)).String())
		})
	}
}

func TestLoader_RoundTrip_Colin_1NTResponses(t *testing.T) {
	is := is.New(t)

	dataFS, err := fs.Sub(ppp.DefaultData, "data")
	is.NoErr(err)
	cs, err := ppp.NewLoader(dataFS).LoadCard("cards/colin.yaml")
	is.NoErr(err)

	tt := []struct {
		name    string
		hand    []string
		want    string
		opening bool
	}{
		{
			// 5 hearts + 5 clubs, 11 HCP: Jacoby fires in response, not Tartan
			name: "Jacoby transfer to hearts - 5+ hearts",
			hand: []string{"K❤️", "Q❤️", "J❤️", "2❤️", "3❤️", "A♣️", "K♣️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "2♠️"},
			want: "2♦️",
		},
		{
			// exactly 5 spades, 8 HCP: Jacoby fires before Baron
			name: "Jacoby transfer to spades - exactly 5 spades",
			hand: []string{"A♠️", "K♠️", "J♠️", "2♠️", "3♠️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️"},
			want: "2❤️",
		},
		{
			// 10 HCP, 4-4-3-2 balanced, no 4-card major — direct 3NT
			// K♣️(3)+J♣️(1)+Q♦️(2)+J♦️(1)+Q❤️(2)+J♠️(1) = 10 HCP
			name: "3NT - balanced 10 HCP no major",
			hand: []string{"K♣️", "J♣️", "2♣️", "3♣️", "Q♦️", "J♦️", "2♦️", "3♦️", "Q❤️", "2❤️", "3❤️", "J♠️", "2♠️"},
			want: "3NT",
		},
		{
			// 10 HCP, 4-4-4-1: not balanced (singleton), so 3NT fails; 4 hearts → Baron
			// A♣️(4)+K♣️(3)+K❤️(3) = 10 HCP
			name: "Baron - 10 HCP with 4-card major",
			hand: []string{"A♣️", "K♣️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "5♦️", "K❤️", "2❤️", "3❤️", "4❤️", "2♠️"},
			want: "2♠️",
		},
		{
			// 4 HCP, 4 hearts: below Baron threshold (8 HCP), Stayman fires
			// J♣️(1)+J❤️(1)+Q❤️(2) = 4 HCP
			name: "Stayman - 4+ major below Baron threshold",
			hand: []string{"J♣️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "5♦️", "J❤️", "Q❤️", "2❤️", "3❤️", "2♠️"},
			want: "2♣️",
		},
		{
			// 6 HCP, 6 spades: 6 ≠ 5 so Jacoby fails; weakness transfer fires
			// A♠️(4)+Q♠️(2) = 6 HCP
			name: "weakness transfer to spades - 6+ spades 6-9 HCP",
			hand: []string{"A♠️", "Q♠️", "2♠️", "3♠️", "4♠️", "5♠️", "2♣️", "3♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️"},
			want: "2NT",
		},
		{
			// 6 HCP, 6 diamonds: transfer to diamonds fires
			// A♦️(4)+Q♦️(2) = 6 HCP
			name: "transfer to diamonds - 6+ diamonds 6-9 HCP",
			hand: []string{"A♦️", "Q♦️", "2♦️", "3♦️", "4♦️", "5♦️", "2♣️", "3♣️", "4♣️", "2❤️", "3❤️", "2♠️", "3♠️"},
			want: "3♣️",
		},
		{
			// Confirm openings still fire with no history (no_prior_bids check)
			name:    "opening 1NT still fires with no history",
			hand:    []string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "3♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️"},
			want:    "1NT",
			opening: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			hand, err := ppp.NewHandFromStrings(tc.hand)
			is.NoErr(err)
			var bidCtx *ppp.AuctionContext
			if tc.opening {
				bidCtx = ctx(hand)
			} else {
				bidCtx = ctxAfterPartnerNT(hand)
			}
			is.Equal(tc.want, cs.Bid(bidCtx).String())
		})
	}
}

package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"
	"gopkg.in/yaml.v3"

	ppp "github.com/cldmstr/passpasspass"
)

// parseCondition parses src, a YAML condition, with ppp.ParseCondition.
func parseCondition(t *testing.T, src string) (ppp.Condition, error) {
	t.Helper()
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(src), &doc); err != nil {
		t.Fatal(err)
	}
	return ppp.ParseCondition(doc.Content[0])
}

func TestParseCondition(t *testing.T) {
	hcp := ppp.NewHCPIntervalCondition(15, 17)
	balanced := ppp.NewBalancedDistributionCondition()
	sixHearts := ppp.NewMinMaxDistributionCondition(ppp.DistributionRange{{0, 13}, {0, 13}, {6, 13}, {0, 13}})
	sixSpades := ppp.NewMinMaxDistributionCondition(ppp.DistributionRange{{0, 13}, {0, 13}, {0, 13}, {6, 13}})

	tt := []struct {
		name string
		src  string
		want ppp.Condition
	}{
		{name: "hcp", src: `hcp: [15, 17]`, want: hcp},
		{name: "hcp at both limits", src: `hcp: [0, 40]`, want: ppp.NewHCPIntervalCondition(0, 40)},
		{name: "hcp with min equal to max", src: `hcp: [10, 10]`, want: ppp.NewHCPIntervalCondition(10, 10)},
		{name: "partnership_hcp at both limits", src: `partnership_hcp: [0, 40]`, want: ppp.NewPartnershipHCPCondition(0, 40)},
		{
			name: "suit_length at both limits",
			src:  `suit_length: {spades: [0, 13]}`,
			want: ppp.NewMinMaxDistributionCondition(ppp.DistributionRange{{0, 13}, {0, 13}, {0, 13}, {0, 13}}),
		},
		{name: "partnership_hcp", src: `partnership_hcp: [25, 40]`, want: ppp.NewPartnershipHCPCondition(25, 40)},
		{name: "balanced", src: `balanced: true`, want: balanced},
		{name: "rule_of_20", src: `rule_of_20: true`, want: ppp.NewRuleOfTwentyCondition()},
		{name: "no_prior_bids", src: `no_prior_bids: true`, want: ppp.NewNoPriorBidsCondition()},
		{name: "partner_bid", src: `partner_bid: "1NT"`, want: ppp.NewPartnerBidCondition("1NT")},
		{name: "partner_bid Pass", src: `partner_bid: Pass`, want: ppp.NewPartnerBidCondition("Pass")},
		{name: "auction", src: `auction: [Pass*, 1H, Pass, Any, Any*]`, want: auctionCondition(t, "Pass*", "1H", "Pass", "Any", "Any*")},
		{name: "auction before the dealer's call", src: `auction: []`, want: auctionCondition(t)},
		{
			name: "suit_length with several suits in any spelling",
			src:  `suit_length: {hearts: [6, 13], S: [0, 3], ♣️: [1, 2]}`,
			want: ppp.NewMinMaxDistributionCondition(ppp.DistributionRange{{1, 2}, {0, 13}, {6, 13}, {0, 3}}),
		},
		{
			name: "two keys in one mapping",
			src:  `{hcp: [15, 17], balanced: true}`,
			want: ppp.NewAllCondition(hcp, balanced),
		},
		{
			name: "keys in a list",
			src:  `[{hcp: [15, 17]}, {balanced: true}]`,
			want: ppp.NewAllCondition(hcp, balanced),
		},
		{name: "list with one item", src: `[{hcp: [15, 17]}]`, want: hcp},
		{
			name: "all",
			src:  `all: [{hcp: [15, 17]}, {balanced: true}]`,
			want: ppp.NewAllCondition(hcp, balanced),
		},
		{
			name: "any with several keys per item",
			src: `
any:
  - hcp: [15, 17]
    balanced: true
  - suit_length: {hearts: [6, 13]}
  - all:
      - suit_length: {spades: [6, 13]}
      - hcp: [15, 17]
`,
			want: ppp.NewMultiCondition(
				ppp.NewAllCondition(hcp, balanced),
				sixHearts,
				ppp.NewAllCondition(sixSpades, hcp),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			got, err := parseCondition(t, tc.src)
			is.NoErr(err)
			is.Equal(got, tc.want)
		})
	}
}

// TestParseCondition_SeveralKeys follows the "several keys" scenarios of the convention-files spec.
func TestParseCondition_SeveralKeys(t *testing.T) {
	weakThreeSpades, _ := ppp.NewHandFromStrings([]string{
		"A♠️", "2♠️", "3♠️", "2♣️", "3♣️", "4♣️", "5♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️", "4❤️",
	})
	weakSixSpades, _ := ppp.NewHandFromStrings([]string{
		"A♠️", "2♠️", "3♠️", "4♠️", "5♠️", "6♠️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "2❤️", "3❤️",
	})
	eightSixSpades, _ := ppp.NewHandFromStrings([]string{
		"A♠️", "K♠️", "2♠️", "3♠️", "4♠️", "5♠️", "J♣️", "2♣️", "3♣️", "2♦️", "3♦️", "2❤️", "3❤️",
	})
	eightFlat, _ := ppp.NewHandFromStrings([]string{
		"A♠️", "K♠️", "2♠️", "J♣️", "2♣️", "3♣️", "4♣️", "2♦️", "3♦️", "4♦️", "2❤️", "3❤️", "4❤️",
	})

	weakTwo := `{hcp: [0, 10], suit_length: {spades: [6, 13]}}`
	multi := `
any:
  - {hcp: [6, 10], suit_length: {hearts: [6, 13]}}
  - {hcp: [6, 10], suit_length: {spades: [6, 13]}}
`

	tt := []struct {
		name   string
		src    string
		hand   *ppp.Hand
		passes bool
	}{
		{name: "mapping: second key fails", src: weakTwo, hand: weakThreeSpades, passes: false},
		{name: "mapping: both keys pass", src: weakTwo, hand: weakSixSpades, passes: true},
		{name: "combinator: second item passes", src: multi, hand: eightSixSpades, passes: true},
		{name: "combinator: no item passes", src: multi, hand: eightFlat, passes: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			c, err := parseCondition(t, tc.src)
			is.NoErr(err)
			is.Equal(c.Eval(ctx(tc.hand)), tc.passes)
		})
	}
}

func TestParseCondition_Errors(t *testing.T) {
	tt := []struct {
		name    string
		src     string
		wantErr string
	}{
		{
			name:    "unknown key",
			src:     `suit_len: {hearts: [6, 13]}`,
			wantErr: "1: suit_len: unknown condition key (known keys: all, any, auction, balanced, hcp, no_prior_bids, partner_bid, partnership_hcp, rule_of_20, suit_length)",
		},
		{
			name:    "unknown key next to a known one",
			src:     "hcp: [15, 17]\nbalance: true",
			wantErr: "2: balance: unknown condition key (known keys: all, any, auction, balanced, hcp, no_prior_bids, partner_bid, partnership_hcp, rule_of_20, suit_length)",
		},
		{
			name:    "unknown key nested in any",
			src:     "any:\n  - hcp: [15, 17]\n  - hpc: [20, 21]",
			wantErr: "3: any.hpc: unknown condition key (known keys: all, any, auction, balanced, hcp, no_prior_bids, partner_bid, partnership_hcp, rule_of_20, suit_length)",
		},
		{
			name:    "error deep inside nested conditions",
			src:     "any:\n  - hcp: [15, 17]\n  - all:\n      - balanced: true\n      - suit_length:\n          clubs: [0, 3]\n          hearts: [0, 14]",
			wantErr: "7: any.all.suit_length.hearts: [0, 14] is outside 0 to 13",
		},
		{name: "empty conditions", src: `null`, wantErr: "1: a condition must be a mapping or a list"},
		{name: "empty mapping", src: `{}`, wantErr: "1: a condition mapping must not be empty"},
		{name: "empty list", src: `[]`, wantErr: "1: a condition list must not be empty"},
		{name: "empty all", src: `all: []`, wantErr: "1: all: a condition list must not be empty"},
		{name: "empty any", src: `any: []`, wantErr: "1: any: a condition list must not be empty"},
		{name: "empty mapping inside all", src: "all:\n  - {}", wantErr: "2: all: a condition mapping must not be empty"},
		{name: "scalar", src: `balanced`, wantErr: "1: a condition must be a mapping or a list"},
		{name: "scalar inside a list", src: "- hcp: [15, 17]\n- balanced", wantErr: "2: a condition must be a mapping or a list"},
		{name: "all is not a list", src: `all: {hcp: [15, 17]}`, wantErr: "1: all: must be a list of conditions"},
		{name: "any is not a list", src: `any: {hcp: [15, 17]}`, wantErr: "1: any: must be a list of conditions"},
		{name: "key repeated in a mapping", src: "hcp: [10, 12]\nhcp: [15, 17]", wantErr: "2: hcp: key appears more than once"},
		{
			name:    "suit repeated in suit_length",
			src:     "suit_length:\n  hearts: [6, 13]\n  H: [0, 3]",
			wantErr: "3: suit_length.H: suit appears more than once",
		},
		{
			name:    "notrump in suit_length",
			src:     `suit_length: {NT: [1, 2]}`,
			wantErr: "1: suit_length.NT: not a suit (clubs, diamonds, hearts or spades)",
		},
		{
			name:    "unknown suit in suit_length",
			src:     `suit_length: {stars: [6, 13]}`,
			wantErr: "1: suit_length.stars: not a suit (clubs, diamonds, hearts or spades)",
		},
		{name: "suit_length not a mapping", src: `suit_length: [6, 13]`, wantErr: "1: suit_length: must be a mapping of suits to [min, max]"},
		{name: "invalid partner_bid", src: `partner_bid: "1N"`, wantErr: `1: partner_bid: invalid bid suit in "1N"`},
		{name: "partner_bid not a call", src: `partner_bid: [1NT]`, wantErr: "1: partner_bid: must be a call such as 1NT or Pass"},
		{name: "hcp with one value", src: `hcp: [5]`, wantErr: "1: hcp: must be [min, max], got 1 values"},
		{name: "hcp with three values", src: `hcp: [5, 6, 7]`, wantErr: "1: hcp: must be [min, max], got 3 values"},
		{name: "hcp not a list", src: `hcp: 5`, wantErr: "1: hcp: must be [min, max] with two whole numbers"},
		{name: "hcp not integers", src: `hcp: [5, ten]`, wantErr: "1: hcp: must be [min, max] with two whole numbers"},
		{name: "hcp reversed", src: `hcp: [12, 10]`, wantErr: "1: hcp: min 12 is greater than max 10"},
		{name: "hcp above 40", src: `hcp: [0, 41]`, wantErr: "1: hcp: [0, 41] is outside 0 to 40"},
		{name: "hcp below 0", src: `hcp: [-1, 10]`, wantErr: "1: hcp: [-1, 10] is outside 0 to 40"},
		{name: "partnership_hcp above 40", src: `partnership_hcp: [25, 41]`, wantErr: "1: partnership_hcp: [25, 41] is outside 0 to 40"},
		{name: "partnership_hcp reversed", src: `partnership_hcp: [26, 25]`, wantErr: "1: partnership_hcp: min 26 is greater than max 25"},
		{name: "suit_length above 13", src: `suit_length: {hearts: [0, 14]}`, wantErr: "1: suit_length.hearts: [0, 14] is outside 0 to 13"},
		{name: "suit_length below 0", src: `suit_length: {hearts: [-1, 3]}`, wantErr: "1: suit_length.hearts: [-1, 3] is outside 0 to 13"},
		{name: "suit_length reversed", src: `suit_length: {spades: [6, 5]}`, wantErr: "1: suit_length.spades: min 6 is greater than max 5"},
		{name: "balanced false", src: `balanced: false`, wantErr: "1: balanced: only true is accepted"},
		{name: "balanced yes please", src: `balanced: yes please`, wantErr: "1: balanced: only true is accepted"},
		{name: "balanced yes", src: `balanced: yes`, wantErr: "1: balanced: only true is accepted"},
		{name: "balanced quoted true", src: `balanced: "true"`, wantErr: "1: balanced: only true is accepted"},
		{name: "balanced empty", src: `balanced:`, wantErr: "1: balanced: only true is accepted"},
		{name: "rule_of_20 false", src: `rule_of_20: false`, wantErr: "1: rule_of_20: only true is accepted"},
		{name: "no_prior_bids false", src: `no_prior_bids: false`, wantErr: "1: no_prior_bids: only true is accepted"},
		{name: "auction with an invalid element", src: `auction: [1H, Foo]`, wantErr: `1: auction: "Foo" is neither a call nor a wildcard (Any, Any*, Pass*)`},
		{
			name:    "auction element on its own line",
			src:     "auction:\n  - Pass*\n  - 1H\n  - pass*",
			wantErr: `4: auction: "pass*" is neither a call nor a wildcard (Any, Any*, Pass*)`,
		},
		{name: "auction not a list", src: `auction: 1H`, wantErr: "1: auction: must be a list of calls and wildcards (Any, Any*, Pass*)"},
		{name: "auction element not a call", src: `auction: [[1H]]`, wantErr: "1: auction: must be a call or a wildcard (Any, Any*, Pass*)"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			_, err := parseCondition(t, tc.src)
			is.True(err != nil)
			is.Equal(err.Error(), tc.wantErr)
		})
	}
}

// auctionCondition returns the auction pattern condition for elements.
func auctionCondition(t *testing.T, elements ...string) ppp.Condition {
	t.Helper()
	pattern, err := ppp.NewAuctionPattern(elements...)
	if err != nil {
		t.Fatal(err)
	}
	return ppp.NewAuctionPatternCondition(pattern)
}

package passpasspass_test

import (
	"strings"
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
		{name: "partnership_hcp", src: `partnership_hcp: [25, 40]`, want: ppp.NewPartnershipHCPCondition(25, 40)},
		{name: "balanced", src: `balanced: true`, want: balanced},
		{name: "rule_of_20", src: `rule_of_20: true`, want: ppp.NewRuleOfTwentyCondition()},
		{name: "no_prior_bids", src: `no_prior_bids: true`, want: ppp.NewNoPriorBidsCondition()},
		{name: "partner_bid", src: `partner_bid: "1NT"`, want: ppp.NewPartnerBidCondition("1NT")},
		{name: "partner_bid Pass", src: `partner_bid: Pass`, want: ppp.NewPartnerBidCondition("Pass")},
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
			wantErr: "suit_len: unknown condition key (known keys: all, any, balanced, hcp, no_prior_bids, partner_bid, partnership_hcp, rule_of_20, suit_length)",
		},
		{
			name:    "unknown key next to a known one",
			src:     `{hcp: [15, 17], balance: true}`,
			wantErr: "balance: unknown condition key",
		},
		{
			name:    "unknown key nested in any",
			src:     `any: [{hcp: [15, 17]}, {hpc: [20, 21]}]`,
			wantErr: "any: hpc: unknown condition key",
		},
		{name: "empty conditions", src: `null`, wantErr: "a condition must be a mapping or a list"},
		{name: "empty mapping", src: `{}`, wantErr: "a condition mapping must not be empty"},
		{name: "empty list", src: `[]`, wantErr: "a condition list must not be empty"},
		{name: "empty all", src: `all: []`, wantErr: "all: a condition list must not be empty"},
		{name: "empty any", src: `any: []`, wantErr: "any: a condition list must not be empty"},
		{name: "empty mapping inside all", src: `all: [{}]`, wantErr: "all: a condition mapping must not be empty"},
		{name: "scalar", src: `balanced`, wantErr: "a condition must be a mapping or a list"},
		{name: "scalar inside a list", src: `[balanced]`, wantErr: "a condition must be a mapping or a list"},
		{name: "all is not a list", src: `all: {hcp: [15, 17]}`, wantErr: "all: must be a list of conditions"},
		{name: "any is not a list", src: `any: {hcp: [15, 17]}`, wantErr: "any: must be a list of conditions"},
		{name: "key repeated in a mapping", src: "hcp: [10, 12]\nhcp: [15, 17]", wantErr: "hcp: condition key appears more than once"},
		{name: "suit repeated in suit_length", src: `suit_length: {hearts: [6, 13], H: [0, 3]}`, wantErr: `suit "H" appears more than once`},
		{name: "notrump in suit_length", src: `suit_length: {NT: [1, 2]}`, wantErr: `unknown suit "NT"`},
		{name: "invalid partner_bid", src: `partner_bid: "1N"`, wantErr: `"1N"`},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			_, err := parseCondition(t, tc.src)
			is.True(err != nil)
			is.True(strings.Contains(err.Error(), tc.wantErr)) // error names the problem
		})
	}
}

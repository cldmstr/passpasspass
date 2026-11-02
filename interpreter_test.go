package passpasspass_test

import (
	"io/fs"
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// colinSheet loads the built-in Colin card.
func colinSheet(t *testing.T) *ppp.ConventionSheet {
	t.Helper()
	dataFS, err := fs.Sub(ppp.DefaultData, "data")
	if err != nil {
		t.Fatal(err)
	}
	cs, err := ppp.NewLoader(dataFS).LoadCard("cards/colin.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return cs
}

// sheetOf returns a sheet with one convention per (bid, condition) pair, in order.
func sheetOf(t *testing.T, conventions ...[2]string) *ppp.ConventionSheet {
	t.Helper()
	cs := ppp.NewConventionSheet()
	for _, c := range conventions {
		bid, err := ppp.NewBidFromString(c[0])
		if err != nil {
			t.Fatal(err)
		}
		condition, err := parseCondition(t, c[1])
		if err != nil {
			t.Fatal(err)
		}
		cs.Add(ppp.NewConvention(bid, condition))
	}
	return cs
}

// callContext returns the context of the next call after the given calls, North dealing, with
// nothing shown yet.
func callContext(t *testing.T, before ...string) *ppp.CallContext {
	t.Helper()
	ctx := &ppp.CallContext{Caller: ppp.Direction(len(before) % 4), History: calls(t, before...)}
	for i := range ctx.Shown {
		ctx.Shown[i] = ppp.DefaultHandModel()
	}
	return ctx
}

// read returns what call shows after the given calls, North dealing, read from sheet.
func read(t *testing.T, sheet *ppp.ConventionSheet, call string, before ...string) ppp.HandModel {
	t.Helper()
	bid, err := ppp.NewBidFromString(call)
	if err != nil {
		t.Fatal(err)
	}
	return ppp.NewSheetInterpreter().Read(sheet, callContext(t, before...), bid)
}

// everySuit returns the same range for all four suits.
func everySuit(lo, hi int) [4]ppp.SuitRange {
	return [4]ppp.SuitRange{{lo, hi}, {lo, hi}, {lo, hi}, {lo, hi}}
}

func TestSheetInterpreter_Colin(t *testing.T) {
	colin := colinSheet(t)

	t.Run("1NT opening shows 15–17 balanced", func(t *testing.T) {
		is := is.New(t)
		got := read(t, colin, "1NT")
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].MinHCP, 15)
		is.Equal(got.Cases[0].MaxHCP, 17)
		is.Equal(got.Cases[0].SuitLength, everySuit(2, 5))
	})

	t.Run("2D response to 1NT shows 5+ hearts, not Multi", func(t *testing.T) {
		is := is.New(t)
		got := read(t, colin, "2D", "1NT", "Pass")
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].SuitLength[ppp.Hearts], ppp.SuitRange{Min: 5, Max: 13})
		is.Equal(got.Cases[0].MinHCP, 0) // a transfer can be weak
	})

	t.Run("Multi 2D opening shows four alternatives", func(t *testing.T) {
		is := is.New(t)
		got := read(t, colin, "2D")
		is.Equal(len(got.Cases), 4)

		is.Equal([2]int{got.Cases[0].MinHCP, got.Cases[0].MaxHCP}, [2]int{23, 24})
		is.Equal(got.Cases[0].SuitLength, everySuit(2, 5))

		is.Equal([2]int{got.Cases[1].MinHCP, got.Cases[1].MaxHCP}, [2]int{6, 10})
		is.Equal(got.Cases[1].SuitLength[ppp.Hearts], ppp.SuitRange{Min: 6, Max: 13})

		is.Equal([2]int{got.Cases[2].MinHCP, got.Cases[2].MaxHCP}, [2]int{6, 10})
		is.Equal(got.Cases[2].SuitLength[ppp.Spades], ppp.SuitRange{Min: 6, Max: 13})

		is.Equal([2]int{got.Cases[3].MinHCP, got.Cases[3].MaxHCP}, [2]int{27, 37})
	})

	t.Run("Pass leaves the model unchanged", func(t *testing.T) {
		is := is.New(t)
		is.Equal(read(t, colin, "Pass", "1NT"), ppp.DefaultHandModel())
	})

	t.Run("a call no convention explains shows nothing", func(t *testing.T) {
		is := is.New(t)
		is.Equal(read(t, colin, "7NT"), ppp.DefaultHandModel())
	})

	t.Run("a call whose conventions cannot apply shows nothing", func(t *testing.T) {
		is := is.New(t)
		// 1NT needs nobody to have opened; after 1C it has no meaning on this card.
		is.Equal(read(t, colin, "1NT", "1C"), ppp.DefaultHandModel())
	})
}

// TestSheetInterpreter_Not follows the negation scenarios of the auction spec's bid inference.
func TestSheetInterpreter_Not(t *testing.T) {
	t.Run("not hcp shows the rest of the range", func(t *testing.T) {
		is := is.New(t)
		got := read(t, sheetOf(t, [2]string{"1C", `not: {hcp: [0, 9]}`}), "1C")
		is.Equal(len(got.Cases), 1)
		is.Equal([2]int{got.Cases[0].MinHCP, got.Cases[0].MaxHCP}, [2]int{10, 37})
	})

	t.Run("not suit_length shows the short side", func(t *testing.T) {
		is := is.New(t)
		got := read(t, sheetOf(t, [2]string{"1C", `not: {suit_length: {hearts: [4, 13]}}`}), "1C")
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].SuitLength[ppp.Hearts], ppp.SuitRange{Min: 0, Max: 3})
	})

	t.Run("not balanced shows nothing", func(t *testing.T) {
		is := is.New(t)
		got := read(t, sheetOf(t, [2]string{"1C", `not: {balanced: true}`}), "1C")
		is.Equal(got, ppp.DefaultHandModel())
	})

	sheet := sheetOf(t,
		[2]string{"1C", `{not: {auction: [Pass*]}, hcp: [0, 9]}`},
		[2]string{"1C", `hcp: [12, 21]`},
	)

	t.Run("a negated auction guard that holds removes the convention", func(t *testing.T) {
		is := is.New(t)
		got := read(t, sheet, "1C")
		is.Equal(len(got.Cases), 1)
		is.Equal(got.Cases[0].MinHCP, 12)
	})

	t.Run("a negated auction guard that fails keeps the convention", func(t *testing.T) {
		is := is.New(t)
		got := read(t, sheet, "1C", "1D")
		is.Equal(len(got.Cases), 2)
		is.Equal(got.Cases[0].MaxHCP, 9)
		is.Equal(got.Cases[1].MinHCP, 12)
	})
}

func TestCondition_Shown(t *testing.T) {
	partnerOpened1NT := func(t *testing.T) *ppp.CallContext {
		t.Helper()
		ctx := callContext(t, "1NT", "Pass")
		ctx.Shown[ppp.North] = ppp.HandModel{Cases: []ppp.HandRange{handRange(15, 17, nil)}}
		return ctx
	}

	tt := []struct {
		name  string
		src   string
		holds bool
		ctx   func(t *testing.T) *ppp.CallContext
		want  []ppp.HandRange
	}{
		{
			name:  "all intersects its children",
			src:   `all: [{hcp: [6, 10]}, {suit_length: {hearts: [6, 13]}}]`,
			holds: true,
			want:  []ppp.HandRange{handRange(6, 10, map[ppp.Suit]ppp.SuitRange{ppp.Clubs: {0, 7}, ppp.Diamonds: {0, 7}, ppp.Hearts: {6, 13}, ppp.Spades: {0, 7}})},
		},
		{
			name:  "any keeps its children as alternatives",
			src:   `any: [{hcp: [6, 10]}, {hcp: [20, 22]}]`,
			holds: true,
			want:  []ppp.HandRange{handRange(6, 10, nil), handRange(20, 22, nil)},
		},
		{
			name:  "failing hcp breaks either bound",
			src:   `hcp: [6, 10]`,
			holds: false,
			want:  []ppp.HandRange{handRange(0, 5, nil), handRange(11, 37, nil)},
		},
		{
			name:  "failing hcp with no lower bound",
			src:   `hcp: [0, 9]`,
			holds: false,
			want:  []ppp.HandRange{handRange(10, 37, nil)},
		},
		{
			name:  "failing an hcp range that covers every hand shows no hand",
			src:   `hcp: [0, 40]`,
			holds: false,
			want:  nil,
		},
		{
			name:  "failing suit_length breaks either bound of each suit",
			src:   `suit_length: {hearts: [4, 6], spades: [0, 3]}`,
			holds: false,
			want: []ppp.HandRange{
				handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {0, 3}}),
				handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Clubs: {0, 6}, ppp.Diamonds: {0, 6}, ppp.Hearts: {7, 13}, ppp.Spades: {0, 6}}),
				handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Clubs: {0, 9}, ppp.Diamonds: {0, 9}, ppp.Hearts: {0, 9}, ppp.Spades: {4, 13}}),
			},
		},
		{
			name:  "not pushes through any by De Morgan's laws",
			src:   `not: {any: [{hcp: [0, 5]}, {hcp: [16, 37]}]}`,
			holds: true,
			want:  []ppp.HandRange{handRange(6, 15, nil)},
		},
		{
			name:  "not pushes through all by De Morgan's laws",
			src:   `not: {hcp: [0, 9], suit_length: {hearts: [4, 13]}}`,
			holds: true,
			want:  []ppp.HandRange{handRange(10, 37, nil), handRange(0, 37, map[ppp.Suit]ppp.SuitRange{ppp.Hearts: {0, 3}})},
		},
		{
			name:  "double negation",
			src:   `not: {not: {hcp: [6, 10]}}`,
			holds: true,
			want:  []ppp.HandRange{handRange(6, 10, nil)},
		},
		{
			name:  "balanced is 2–5 in every suit",
			src:   `balanced: true`,
			holds: true,
			want:  []ppp.HandRange{{MinHCP: 0, MaxHCP: 37, SuitLength: everySuit(2, 5)}},
		},
		{name: "failing balanced shows every hand", src: `balanced: true`, holds: false, want: []ppp.HandRange{ppp.DefaultHandRange()}},
		{name: "rule_of_20 shows every hand", src: `rule_of_20: true`, holds: true, want: []ppp.HandRange{ppp.DefaultHandRange()}},
		{name: "failing rule_of_20 shows every hand", src: `rule_of_20: true`, holds: false, want: []ppp.HandRange{ppp.DefaultHandRange()}},
		{
			name:  "partnership_hcp becomes own HCP",
			src:   `partnership_hcp: [25, 40]`,
			holds: true,
			ctx:   partnerOpened1NT,
			want:  []ppp.HandRange{handRange(10, 23, nil)},
		},
		{
			name:  "failing partnership_hcp shows every hand",
			src:   `partnership_hcp: [25, 40]`,
			holds: false,
			ctx:   partnerOpened1NT,
			want:  []ppp.HandRange{ppp.DefaultHandRange()},
		},
		{
			name:  "partnership_hcp cannot hold while partner has shown nothing",
			src:   `partnership_hcp: [25, 40]`,
			holds: true,
			want:  nil,
		},
		{name: "auction guard that holds", src: `auction: [Pass*]`, holds: true, want: []ppp.HandRange{ppp.DefaultHandRange()}},
		{
			name:  "auction guard that fails",
			src:   `auction: [Pass*]`,
			holds: true,
			ctx:   func(t *testing.T) *ppp.CallContext { t.Helper(); return callContext(t, "1C") },
			want:  nil,
		},
		{name: "negated auction guard that holds", src: `not: {auction: [Pass*]}`, holds: true, want: nil},
		{name: "partner_bid that holds", src: `partner_bid: "1NT"`, holds: true, ctx: partnerOpened1NT, want: []ppp.HandRange{ppp.DefaultHandRange()}},
		{name: "partner_bid that fails", src: `partner_bid: "1NT"`, holds: true, want: nil},
		{name: "no_prior_bids that holds", src: `no_prior_bids: true`, holds: true, want: []ppp.HandRange{ppp.DefaultHandRange()}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			c, err := parseCondition(t, tc.src)
			is.NoErr(err)
			ctx := callContext(t)
			if tc.ctx != nil {
				ctx = tc.ctx(t)
			}
			is.Equal(c.Shown(ctx, tc.holds).Cases, tc.want)
		})
	}
}

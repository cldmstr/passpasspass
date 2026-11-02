package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// calls returns the history of the given calls, the first made by North as dealer.
func calls(t *testing.T, bids ...string) ppp.History {
	t.Helper()
	h := make(ppp.History, len(bids))
	for i, s := range bids {
		bid, err := ppp.NewBidFromString(s)
		if err != nil {
			t.Fatal(err)
		}
		h[i] = ppp.BidAction{Direction: ppp.Direction(i % 4), Bid: bid}
	}
	return h
}

func TestNewAuctionPatternCondition(t *testing.T) {
	tt := []struct {
		name    string
		pattern []string
		calls   []string
		want    bool
	}{
		{name: "empty pattern before the dealer's call", pattern: nil, calls: nil, want: true},
		{name: "empty pattern after the dealer passed", pattern: nil, calls: []string{"Pass"}, want: false},

		{name: "dealer opening", pattern: []string{"1H", "Pass"}, calls: []string{"1H", "Pass"}, want: true},
		{name: "third-seat opening is not a dealer opening", pattern: []string{"1H", "Pass"}, calls: []string{"Pass", "Pass", "1H", "Pass"}, want: false},

		{name: "Pass* with a first-seat opening", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"1H", "Pass"}, want: true},
		{name: "Pass* with a second-seat opening", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"Pass", "1H", "Pass"}, want: true},
		{name: "Pass* with a third-seat opening", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"Pass", "Pass", "1H", "Pass"}, want: true},
		{name: "Pass* with a fourth-seat opening", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"Pass", "Pass", "Pass", "1H", "Pass"}, want: true},

		// After Pass 1H Pass, West is to call and has not passed. After Pass Pass 1H Pass, North is
		// to call and passed as dealer.
		{name: "unpassed responder to a dealer opening", pattern: []string{"1H", "Pass"}, calls: []string{"1H", "Pass"}, want: true},
		{name: "unpassed responder to a second-seat opening", pattern: []string{"Pass", "1H", "Pass"}, calls: []string{"Pass", "1H", "Pass"}, want: true},
		{name: "passed responder is not an unpassed one", pattern: []string{"Pass", "1H", "Pass"}, calls: []string{"Pass", "Pass", "1H", "Pass"}, want: false},

		{name: "Any matches one call", pattern: []string{"1H", "Pass", "Any", "Pass"}, calls: []string{"1H", "Pass", "2H", "Pass"}, want: true},
		{name: "Any matches a pass", pattern: []string{"1H", "Pass", "Any", "Pass"}, calls: []string{"1H", "Pass", "Pass", "Pass"}, want: true},
		{name: "Any does not match zero calls", pattern: []string{"1H", "Pass", "Any", "Pass"}, calls: []string{"1H", "Pass", "Pass"}, want: false},
		{name: "Any does not match two calls", pattern: []string{"1H", "Pass", "Any", "Pass"}, calls: []string{"1H", "Pass", "2H", "3H", "Pass"}, want: false},

		{name: "Any* matches zero calls", pattern: []string{"Pass*", "1H", "Pass", "Any*"}, calls: []string{"1H", "Pass"}, want: true},
		{name: "Any* matches one call", pattern: []string{"Pass*", "1H", "Pass", "Any*"}, calls: []string{"1H", "Pass", "2H"}, want: true},
		{name: "Any* matches several calls", pattern: []string{"Pass*", "1H", "Pass", "Any*"}, calls: []string{"1H", "Pass", "2H", "Pass", "3H", "Pass"}, want: true},
		{name: "Any* in the middle", pattern: []string{"1H", "Any*", "4H"}, calls: []string{"1H", "Pass", "2H", "Pass", "4H"}, want: true},
		{name: "Any* in the middle needs the last call", pattern: []string{"1H", "Any*", "4H"}, calls: []string{"1H", "Pass", "4S"}, want: false},

		{name: "Pass* does not match a bid", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"1C", "Pass", "1H", "Pass"}, want: false},
		{name: "Pass* matches no calls", pattern: []string{"Pass*"}, calls: nil, want: true},
		{name: "Pass* matches only passes", pattern: []string{"Pass*"}, calls: []string{"Pass", "Pass", "Pass"}, want: true},
		{name: "Pass* after an opening", pattern: []string{"Pass*"}, calls: []string{"Pass", "1C"}, want: false},
		{name: "Pass* gives a pass back to a following Pass", pattern: []string{"Pass*", "Pass"}, calls: []string{"Pass", "Pass"}, want: true},
		{name: "Pass* followed by Pass needs one pass", pattern: []string{"Pass*", "Pass"}, calls: nil, want: false},

		{name: "longer history does not match", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"1H", "Pass", "2H", "Pass"}, want: false},
		{name: "opponent's bid in place of a pass", pattern: []string{"Pass*", "1H", "Pass"}, calls: []string{"1H", "1S"}, want: false},
		{name: "different call", pattern: []string{"1H", "Pass"}, calls: []string{"1S", "Pass"}, want: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			pattern, err := ppp.NewAuctionPattern(tc.pattern...)
			is.NoErr(err)
			c := ppp.NewAuctionPatternCondition(pattern)
			is.Equal(c.Eval(&ppp.AuctionContext{History: calls(t, tc.calls...)}), tc.want)
		})
	}
}

func TestNewAuctionPattern_InvalidElement(t *testing.T) {
	for _, element := range []string{"Foo", "any", "Pass**", "*", "8H", ""} {
		t.Run(element, func(t *testing.T) {
			is := is.New(t)
			_, err := ppp.NewAuctionPattern("1H", element)
			is.True(err != nil)
			is.Equal(err.Error(), `"`+element+`" is neither a call nor a wildcard (Any, Any*, Pass*)`)
		})
	}
}

// TestAuctionCondition_PassedResponder follows the "Passed responder" scenario of the
// bidding-conditions spec: two unpassed patterns under `any` reject a passed responder.
func TestAuctionCondition_PassedResponder(t *testing.T) {
	is := is.New(t)
	c, err := parseCondition(t, `any: [{auction: [1H, Pass]}, {auction: [Pass, 1H, Pass]}]`)
	is.NoErr(err)

	is.True(c.Eval(&ppp.AuctionContext{History: calls(t, "1H", "Pass")}))                  // South, unpassed
	is.True(c.Eval(&ppp.AuctionContext{History: calls(t, "Pass", "1H", "Pass")}))          // West, unpassed
	is.True(!c.Eval(&ppp.AuctionContext{History: calls(t, "Pass", "Pass", "1H", "Pass")})) // North, passed as dealer
}

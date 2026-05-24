package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// mustHand builds a Hand from card strings, panicking on error.
func mustHand(cards []string) *ppp.Hand {
	h, err := ppp.NewHandFromStrings(cards)
	if err != nil {
		panic(err)
	}
	return h
}

// oneSuitedDeal returns four non-overlapping 13-card hands, one suit each.
// North=clubs, East=diamonds, South=hearts, West=spades.
func oneSuitedDeal() []*ppp.Hand {
	return []*ppp.Hand{
		mustHand([]string{"2♣️", "3♣️", "4♣️", "5♣️", "6♣️", "7♣️", "8♣️", "9♣️", "10♣️", "J♣️", "Q♣️", "K♣️", "A♣️"}),
		mustHand([]string{"2♦️", "3♦️", "4♦️", "5♦️", "6♦️", "7♦️", "8♦️", "9♦️", "10♦️", "J♦️", "Q♦️", "K♦️", "A♦️"}),
		mustHand([]string{"2❤️", "3❤️", "4❤️", "5❤️", "6❤️", "7❤️", "8❤️", "9❤️", "10❤️", "J❤️", "Q❤️", "K❤️", "A❤️"}),
		mustHand([]string{"2♠️", "3♠️", "4♠️", "5♠️", "6♠️", "7♠️", "8♠️", "9♠️", "10♠️", "J♠️", "Q♠️", "K♠️", "A♠️"}),
	}
}

func TestAuction_Run_PassedOut(t *testing.T) {
	is := is.New(t)

	cs := ppp.NewConventionSheet() // empty — every player passes
	a := ppp.NewAuction(ppp.North, cs, oneSuitedDeal())
	a.Run()

	is.True(a.Contract == nil)
	is.Equal(4, len(a.Bids))
	for _, ba := range a.Bids {
		is.Equal(ppp.Pass, ba.Bid)
	}
}

func TestAuction_Run_SingleOpener(t *testing.T) {
	is := is.New(t)

	// North: balanced 4-3-3-3, 16 HCP (A♣K♣Q♣ + A♦K♦).
	// East/South: 0–2 HCP. West: 22 HCP (above the 1NT ceiling).
	// Deck is fully partitioned across the four hands.
	hands := []*ppp.Hand{
		mustHand([]string{"A♣️", "K♣️", "Q♣️", "2♣️", "A♦️", "K♦️", "2♦️", "2❤️", "3❤️", "4❤️", "2♠️", "3♠️", "4♠️"}),
		mustHand([]string{"3♣️", "4♣️", "5♣️", "3♦️", "4♦️", "5♦️", "5❤️", "6❤️", "7❤️", "8❤️", "5♠️", "6♠️", "7♠️"}),
		mustHand([]string{"6♣️", "7♣️", "8♣️", "6♦️", "7♦️", "8♦️", "9❤️", "10❤️", "J❤️", "8♠️", "9♠️", "10♠️", "J♠️"}),
		mustHand([]string{"9♣️", "10♣️", "J♣️", "9♦️", "10♦️", "J♦️", "Q♦️", "Q❤️", "K❤️", "A❤️", "Q♠️", "K♠️", "A♠️"}),
	}

	cs := ppp.NewConventionSheet()
	cs.Add(ppp.NewConvention(
		ppp.NewBid(1, ppp.NoTrump),
		ppp.NewBalancedDistributionCondition(),
		ppp.NewHCPIntervalCondition(15, 17),
	))

	a := ppp.NewAuction(ppp.North, cs, hands)
	a.Run()

	is.Equal("1NT", a.Contract.String())
	is.Equal(ppp.North, a.Declarer)
	is.Equal(4, len(a.Bids))
	is.Equal("1NT", a.Bids[0].Bid.String()) // North opens
	is.Equal(ppp.Pass, a.Bids[1].Bid)       // East: 0 HCP
	is.Equal(ppp.Pass, a.Bids[2].Bid)       // South: 2 HCP
	is.Equal(ppp.Pass, a.Bids[3].Bid)       // West: 22 HCP > ceiling
}

func TestAuction_Run_BidLegalityEnforcement(t *testing.T) {
	is := is.New(t)

	// All four players have 10 HCP and qualify for 1♥.
	// After North opens, East/South/West try the same bid — it must be coerced to Pass.
	cs := ppp.NewConventionSheet()
	cs.Add(ppp.NewConvention(
		ppp.NewBid(1, ppp.Hearts),
		ppp.NewHCPIntervalCondition(10, 37),
	))

	a := ppp.NewAuction(ppp.North, cs, oneSuitedDeal())
	a.Run()

	is.Equal("1❤️", a.Contract.String())
	is.Equal(ppp.North, a.Declarer)
	is.Equal(4, len(a.Bids))
	is.Equal("1❤️", a.Bids[0].Bid.String()) // North opens
	is.Equal(ppp.Pass, a.Bids[1].Bid)        // East: 1♥ ≤ contract → coerced
	is.Equal(ppp.Pass, a.Bids[2].Bid)        // South: same
	is.Equal(ppp.Pass, a.Bids[3].Bid)        // West: same
}

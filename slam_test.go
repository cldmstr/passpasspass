package passpasspass_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

// TestAuction_GrandSlam runs a deal on which East–West hold thirteen top tricks and should bid a
// grand slam. West, the dealer, has 21 HCP with A–K–Q–J to seven hearts, the black aces and
// ♠️K, and ♦️A; East has 12 HCP with A–K–Q–J to six clubs, ♠️Q and a diamond void. Seven heart
// tricks, three spades, ♦️A and the top clubs make 13 tricks in hearts, and more in notrump.
func TestAuction_GrandSlam(t *testing.T) {
	t.Skip("no responses to 2♣️ and no slam conventions yet: West opens 2♣️ and everyone passes")

	is := is.New(t)
	hands := []*ppp.Hand{
		// North: 5 HCP
		mustHand([]string{"6♠️", "J♠️", "2❤️", "3❤️", "6❤️", "2♦️", "5♦️", "6♦️", "10♦️", "J♦️", "K♦️", "3♣️", "5♣️"}),
		// East: 12 HCP, void in diamonds
		mustHand([]string{"8♠️", "9♠️", "10♠️", "Q♠️", "4❤️", "5❤️", "7❤️", "7♣️", "8♣️", "J♣️", "Q♣️", "K♣️", "A♣️"}),
		// South: 2 HCP, void in hearts
		mustHand([]string{"2♠️", "4♠️", "5♠️", "7♠️", "3♦️", "4♦️", "8♦️", "Q♦️", "2♣️", "4♣️", "6♣️", "9♣️", "10♣️"}),
		// West: 21 HCP, void in clubs
		mustHand([]string{"3♠️", "K♠️", "A♠️", "8❤️", "9❤️", "10❤️", "J❤️", "Q❤️", "K❤️", "A❤️", "7♦️", "9♦️", "A♦️"}),
	}

	a := newAuction(ppp.West, colinSheet(t), hands)
	a.Run()

	is.True(a.Contract != nil)                           // someone bid
	is.True(strings.HasPrefix(a.Contract.String(), "7")) // at the seven level
	is.True(a.Declarer == ppp.East || a.Declarer == ppp.West)
}

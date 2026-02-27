package passpasspass_test

import (
	"fmt"
	"testing"

	ppp "github.com/cldmstr/passpasspass"
)

func TestStandardDeck(t *testing.T) {
	sut := ppp.StandardDeck()

	fmt.Printf("%v", sut)

	for r := range 13 {
		want := fmt.Sprintf("%s♣", ppp.Rank(r+2).String())
		if sut[int64(r)].String() != want {
			t.Errorf("got %s; want %s", sut[0], want)
		}
	}
	for r := range 13 {
		want := fmt.Sprintf("%s♦", ppp.Rank(r+2).String())
		if sut[int64(13+r)].String() != want {
			t.Errorf("got %s; want %s", sut[0], want)
		}
	}
	for r := range 13 {
		want := fmt.Sprintf("%s♥", ppp.Rank(r+2).String())
		if sut[int64(26+r)].String() != want {
			t.Errorf("got %s; want %s", sut[0], want)
		}
	}
	for r := range 13 {
		want := fmt.Sprintf("%s♠", ppp.Rank(r+2).String())
		if sut[int64(39+r)].String() != want {
			t.Errorf("got %s; want %s", sut[0], want)
		}
	}
}

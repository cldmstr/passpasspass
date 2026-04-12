package passpasspass_test

import (
	"testing"

	"github.com/matryer/is"

	ppp "github.com/cldmstr/passpasspass"
)

func TestNewSuitFromString(t *testing.T) {
	tt := []struct {
		name    string
		input   string
		want    ppp.Suit
		wantErr bool
	}{
		{name: "clubs symbol", input: "♣️", want: ppp.Clubs},
		{name: "diamonds symbol", input: "♦️", want: ppp.Diamonds},
		{name: "hearts symbol", input: "❤️", want: ppp.Hearts},
		{name: "spades symbol", input: "♠️", want: ppp.Spades},
		{name: "notrump symbol", input: "NT", want: ppp.NoTrump},
		{name: "clubs letter", input: "C", want: ppp.Clubs},
		{name: "diamonds letter", input: "D", want: ppp.Diamonds},
		{name: "hearts letter", input: "H", want: ppp.Hearts},
		{name: "spades letter", input: "S", want: ppp.Spades},
		{name: "notrump letter", input: "NT", want: ppp.NoTrump},
		{name: "empty string", input: "", wantErr: true},
		{name: "unknown letter", input: "X", wantErr: true},
		{name: "lowercase clubs", input: "clubs", wantErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			suit, err := ppp.NewSuitFromString(tc.input)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(tc.want, suit)
		})
	}
}

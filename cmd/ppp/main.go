// Command ppp deals a random board and prints the auction bid by a convention card: the
// embedded Colin card by default, or the card at <data dir> <card path> when both are given.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	ppp "github.com/cldmstr/passpasspass"
)

func main() {
	var cs *ppp.ConventionSheet
	var err error
	if len(os.Args) >= 3 {
		cs, err = ppp.NewLoader(os.DirFS(os.Args[1])).LoadCard(os.Args[2])
	} else {
		dataFS, _ := fs.Sub(ppp.DefaultData, "data")
		cs, err = ppp.NewLoader(dataFS).LoadCard("cards/colin.yaml")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading convention card: %v\n", err)
		os.Exit(1)
	}

	dealer := ppp.NewRandomDealer(time.Now().UnixNano())
	table := ppp.NewTable(cs, dealer)
	auction := table.RunAuction()
	fmt.Print(auction)
}

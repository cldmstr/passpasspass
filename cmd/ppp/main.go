package main

import (
	"fmt"
	"os"
	"time"

	ppp "github.com/cldmstr/passpasspass"
)

func main() {
	dealer := ppp.NewRandomDealer(time.Now().UnixNano())

	var cs *ppp.ConventionSheet
	if len(os.Args) >= 3 {
		cardPath := os.Args[1]
		conventionsDir := os.Args[2]
		var err error
		cs, err = ppp.NewLoader(conventionsDir).LoadCard(cardPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading convention card: %v\n", err)
			os.Exit(1)
		}
	} else {
		cs = ppp.Twooverone
	}

	table := ppp.NewTable(cs, dealer.Deal())
	fmt.Printf("%s", table)
}

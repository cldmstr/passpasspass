package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	ppp "github.com/cldmstr/passpasspass"
)

func main() {
	dealer := ppp.NewRandomDealer(time.Now().UnixNano())

	var cs *ppp.ConventionSheet
	var err error
	if len(os.Args) >= 3 {
		cs, err = ppp.NewLoader(os.DirFS(os.Args[1])).LoadCard(os.Args[2])
	} else {
		dataFS, _ := fs.Sub(ppp.DefaultData, "data")
		cs, err = ppp.NewLoader(dataFS).LoadCard("cards/two_over_one.yaml")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading convention card: %v\n", err)
		os.Exit(1)
	}

	table := ppp.NewTable(cs, dealer.Deal())
	fmt.Printf("%s", table)
}

package main

import (
	"fmt"
	"time"

	ppp "github.com/cldmstr/passpasspass"
)

func main() {
	dealer := ppp.NewRandomDealer(time.Now().UnixNano())
	table := ppp.NewTable(ppp.Twooverone, dealer.Deal())

	fmt.Printf("%s", table)
}

package main

import (
	"fmt"
	"time"

	ppp "github.com/cldmstr/passpasspass"
)

func main() {
	dealer := ppp.NewDealer(time.Now().UnixNano())
	table := dealer.Deal()

	fmt.Printf("%v", table)
}

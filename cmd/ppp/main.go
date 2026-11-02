// Command ppp deals a random board and prints the auction bid by a convention card: the
// embedded Colin card by default, or the card at <data dir> <card path> when both are given.
// With -v or --verbose it also prints the context every call was made in.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	ppp "github.com/cldmstr/passpasspass"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run parses args, loads the card, and prints one auction to stdout. It returns the exit status.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("ppp", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var verbose bool
	flags.BoolVar(&verbose, "v", false, "print the context of every call")
	flags.BoolVar(&verbose, "verbose", false, "print the context of every call")
	flags.Usage = func() {
		// A failed write to stderr leaves nowhere to report it.
		_, _ = fmt.Fprintln(stderr, "usage: ppp [-v|--verbose] [<data dir> <card path>]")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	cs, err := loadCard(flags.Args())
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error loading convention card: %v\n", err)
		return 1
	}

	dealer := ppp.NewRandomDealer(time.Now().UnixNano())
	table := ppp.NewTable(cs, dealer)
	auction := table.RunAuction()
	if _, err := io.WriteString(stdout, format(auction, verbose)); err != nil {
		return 1
	}
	return 0
}

// format shows the auction: the hands, then the context of every call when verbose, then the
// bidding.
func format(a *ppp.Auction, verbose bool) string {
	if !verbose {
		return a.String()
	}
	var sb strings.Builder
	sb.WriteString(a.HandsString())
	for i, step := range a.Steps {
		fmt.Fprintf(&sb, "%d. %s calls %s\n", i+1, step.Context.Position, step.Call)
		for line := range strings.Lines(step.Context.String()) {
			sb.WriteString("   " + line)
		}
		sb.WriteString("\n")
	}
	sb.WriteString(a.BiddingString())
	return sb.String()
}

// loadCard loads the card at <data dir> <card path> when both are given, else the embedded Colin
// card.
func loadCard(args []string) (*ppp.ConventionSheet, error) {
	if len(args) >= 2 {
		return ppp.NewLoader(os.DirFS(args[0])).LoadCard(args[1])
	}
	dataFS, err := fs.Sub(ppp.DefaultData, "data")
	if err != nil {
		return nil, err
	}
	return ppp.NewLoader(dataFS).LoadCard("cards/colin.yaml")
}

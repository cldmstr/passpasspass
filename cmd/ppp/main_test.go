package main

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// entry matches the first line of a verbose entry, such as "3. South calls 2♦️".
var entry = regexp.MustCompile(`(?m)^(\d+)\. (North|East|South|West) calls \S+$`)

func TestRun_Verbose(t *testing.T) {
	for _, flag := range []string{"-v", "--verbose"} {
		t.Run(flag, func(t *testing.T) {
			is := is.New(t)
			var stdout, stderr bytes.Buffer
			is.Equal(run([]string{flag}, &stdout, &stderr), 0)
			out := stdout.String()

			entries := entry.FindAllStringSubmatch(out, -1)
			is.True(len(entries) >= 4) // every auction has at least four calls
			for i, e := range entries {
				is.Equal(e[1], strconv.Itoa(i+1)) // numbered in order
			}

			// The entries sit between the hands and the bidding grid.
			first, last := strings.Index(out, entries[0][0]), strings.LastIndex(out, "Partnership ")
			is.True(strings.Index(out, "West (") < first)
			is.True(last < strings.Index(out, "Dealer: "))
		})
	}
}

func TestRun_Plain(t *testing.T) {
	is := is.New(t)
	var stdout, stderr bytes.Buffer
	is.Equal(run(nil, &stdout, &stderr), 0)
	is.True(!entry.MatchString(stdout.String()))
	is.True(strings.Contains(stdout.String(), "Dealer: "))
	is.Equal(stderr.String(), "")
}

func TestRun_MissingCard(t *testing.T) {
	is := is.New(t)
	var stdout, stderr bytes.Buffer
	is.Equal(run([]string{"-v", t.TempDir(), "cards/missing.yaml"}, &stdout, &stderr), 1)
	is.True(strings.HasPrefix(stderr.String(), "error loading convention card: "))
	is.Equal(stdout.String(), "")
}

func TestRun_UnknownFlag(t *testing.T) {
	is := is.New(t)
	var stdout, stderr bytes.Buffer
	is.Equal(run([]string{"-x"}, &stdout, &stderr), 2)
	is.True(strings.Contains(stderr.String(), "usage: ppp [-v|--verbose] [<data dir> <card path>]"))
}

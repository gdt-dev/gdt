package cli

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

const (
	runeLightHorizontal rune = '\u2500' // ─
)

var (
	hbar string
)

func init() {
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 80
	}
	hbar = strings.Repeat(string(runeLightHorizontal), width)
}

func HorizontalBar() {
	fmt.Printf("%s\n", hbar)
}

func HorizontalSectionHeader(header string) {
	HorizontalBar()
	// Centers the header...
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 80
	}
	padding := strings.Repeat(" ", (width-len(header))/2)
	fmt.Printf("%s%s\n", padding, strings.ToUpper(header))
	HorizontalBar()
}

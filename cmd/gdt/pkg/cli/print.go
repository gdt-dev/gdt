package cli

import (
	"fmt"
	"os"
	"strings"
)

// V prints the supplied message followed by a newline if the -v|--verbose CLI
// flag is true.
func V(msg string) {
	if CommonOptions.Verbose {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// Vf print-formats the supplied format string and format args if the
// -v|--verbose CLI flag is true.
func Vf(msg string, args ...any) {
	if CommonOptions.Verbose {
		out := fmt.Sprintf(msg, args...)
		fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
	}
}

// D prints the supplied message followed by a newline if the -d|--debug CLI
// flag is true.
func D(msg string) {
	if CommonOptions.Debug {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// Df print-formats the supplied format string and format args if the
// -d|--debug CLI flag is true.
func Df(msg string, args ...any) {
	if CommonOptions.Debug {
		out := fmt.Sprintf(msg, args...)
		fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
	}
}

// Ellipsis cuts the supplied string at some max length, adding an ellipsis to
// the end of any string that is longer than the max length.
func Ellipsis(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen < 3 {
		maxLen = 3
	}
	return string(runes[0:maxLen-3]) + "..."
}

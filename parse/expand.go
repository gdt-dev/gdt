// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package parse

import (
	"fmt"
	"os"
	"strings"
)

const (
	// hopefully nobody actually has an environment variable with this key!
	dollarSignReplacementToken = "oiuqdfjhaso7t213041"
)

// ExpandWithFixedDoubleDollar expands the given string using os.ExpandEnv,
// however unlike the default behaviour of replacing a string "$$VALUE" with
// "VALUE", it replaces the "$$" witha single "$". This allows test authors to
// use the dollar symbol in their test contents (they need to escape with
// '$$').
func ExpandWithFixedDoubleDollar(subject string) string {
	os.Setenv(dollarSignReplacementToken, "$")
	replaceStr := fmt.Sprintf("${%s}", dollarSignReplacementToken)
	return os.ExpandEnv(strings.Replace(subject, "$$", replaceStr, -1))
}

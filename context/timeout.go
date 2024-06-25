// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"
	"errors"
	"strings"

	"github.com/gdt-dev/gdt/api"
)

// TimedOut evaluates a context and an arbitrary other error and returns true
// if the context timed out or the error indicates a deadline exceeded or
// signal was killed.
func TimedOut(
	ctx context.Context,
	err error,
) bool {
	// Handle context timeouts before addressing other errors
	if ctx.Err() != nil {
		cerr := ctx.Err()
		return errors.Is(cerr, context.DeadlineExceeded)
	}
	if err != nil {
		if errors.Is(err, api.ErrTimeoutExceeded) {
			return true
		}
		return strings.Contains(err.Error(), "signal: killed")
	}
	return false
}

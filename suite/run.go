// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite

import (
	"context"
	"testing"
)

// Run executes the tests in the test case
func (s *Suite) Run(ctx context.Context, t *testing.T) error {
	for _, sc := range s.Scenarios {
		if err := sc.Run(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

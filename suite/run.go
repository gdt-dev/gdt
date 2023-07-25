// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite

import (
	"context"
	"testing"

	gdterrors "github.com/gdt-dev/gdt/errors"
)

// Run executes the tests in the test case
func (s *Suite) Run(ctx context.Context, t *testing.T) error {
	errs := gdterrors.NewRuntimeErrors()
	for _, unit := range s.Scenarios {
		errs.AppendIf(unit.Run(ctx, t))
	}
	if errs.Empty() {
		return nil
	}
	return errs
}

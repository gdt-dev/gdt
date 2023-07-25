// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"context"
	"strings"
	"testing"
	"time"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/result"
	gdttypes "github.com/gdt-dev/gdt/types"
)

// Run executes the tests in the test scenario
func (s *Scenario) Run(ctx context.Context, t *testing.T) error {
	if len(s.Require) > 0 {
		fixtures := gdtcontext.Fixtures(ctx)
		for _, fname := range s.Require {
			lookup := strings.ToLower(fname)
			fix, found := fixtures[lookup]
			if !found {
				return gdterrors.RequiredFixtureMissing(fname)
			}
			fix.Start()
			defer fix.Stop()
		}
	}
	errs := gdterrors.NewRuntimeErrors()
	var scDefaults *Defaults
	scDefaultsAny, found := s.Defaults[DefaultsKey]
	if found {
		scDefaults = scDefaultsAny.(*Defaults)
	}
	t.Run(s.Title(), func(t *testing.T) {
		for _, spec := range s.Tests {
			// Create a brand new context that inherits the top-level context's
			// cancel func. We want to set deadlines for each test spec and if
			// we mutate the single supplied top-level context, then only the
			// first deadline/timeout will be used.
			specCtx, specCancel := context.WithCancel(ctx)
			defer specCancel()

			sb := spec.Base()
			wait := sb.Wait
			if wait != nil && wait.Before != "" {
				debug.Println(ctx, t, "wait: %s before", wait.Before)
				time.Sleep(wait.BeforeDuration())
			}

			to := specTimeout(ctx, t, sb.Timeout, scDefaults)
			if to != nil {
				var cancel context.CancelFunc
				specCtx, cancel = context.WithTimeout(specCtx, to.Duration())
				defer cancel()
			}
			err := spec.Run(specCtx, t)
			if gdtcontext.TimedOut(specCtx, err) {
				if to != nil && !to.Expected {
					t.Fatal(gdterrors.TimeoutExceeded(to.After))
				}
				// Swallow the error since it's not a runtime error but rather
				// an assertion failure.
				err = nil
			}
			if res, ok := err.(*result.Result); ok {
				// Results can have arbitrary run data stored in them and we
				// save this prior run data in the top-level context (and pass
				// that context to the next Run invocation).
				if res.HasData() {
					ctx = gdtcontext.StorePriorRun(ctx, res.Data())
				}
				errs.AppendIf(res.Unwrap())
			} else {
				errs.AppendIf(err)
			}
			if wait != nil && wait.After != "" {
				debug.Println(ctx, t, "wait: %s after", wait.After)
				time.Sleep(wait.AfterDuration())
			}
		}
	})
	if errs.Empty() {
		return nil
	}
	return errs
}

// specTimeout returns the timeout value for the test spec. If the spec has a
// timeout override, we use that. Otherwise, we inspect the scenario's defaults
// and, if present, use that timeout.
func specTimeout(
	ctx context.Context,
	t *testing.T,
	specTimeout *gdttypes.Timeout,
	scenDefaults *Defaults,
) *gdttypes.Timeout {
	if specTimeout != nil {
		debug.Println(
			ctx, t, "using timeout of %s (expected: %t)",
			specTimeout.After, specTimeout.Expected,
		)
		return specTimeout
	}
	if scenDefaults != nil && scenDefaults.Timeout != nil {
		debug.Println(
			ctx, t, "using timeout of %s (expected: %t) [scenario default]",
			scenDefaults.Timeout.After, scenDefaults.Timeout.Expected,
		)
		return scenDefaults.Timeout
	}
	return nil
}

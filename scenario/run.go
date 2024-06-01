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
	gdttypes "github.com/gdt-dev/gdt/types"
)

// Run executes the scenario. The error that is returned will always be derived
// from `gdterrors.RuntimeError` and represents an *unrecoverable* error.
//
// Test assertion failures are *not* considered errors. The Scenario.Run()
// method controls whether `testing.T.Fail()` or `testing.T.Skip()` is called
// which will mark the test units failed or skipped if a test unit evaluates to
// false.
func (s *Scenario) Run(ctx context.Context, t *testing.T) error {
	if len(s.Fixtures) > 0 {
		fixtures := gdtcontext.Fixtures(ctx)
		for _, fname := range s.Fixtures {
			lookup := strings.ToLower(fname)
			fix, found := fixtures[lookup]
			if !found {
				return gdterrors.RequiredFixtureMissing(fname)
			}
			fix.Start()
			defer fix.Stop()
		}
	}
	var rterr error
	var scDefaults *Defaults
	scDefaultsAny, found := s.Defaults[DefaultsKey]
	if found {
		scDefaults = scDefaultsAny.(*Defaults)
	}
	// If the test author has specified any pre-flight checks in the `skip-if`
	// collection, evaluate those first and if any failed, skip the scenario's
	// tests.
	for _, skipIf := range s.SkipIf {
		res := skipIf.Eval(ctx, t)
		if res.HasRuntimeError() {
			return res.RuntimeError()
		}
		if len(res.Failures()) == 0 {
			t.Skipf(
				"skip-if: %s passed. skipping test.",
				skipIf.Base().Title(),
			)
			return nil
		}
	}
	t.Run(s.Title(), func(t *testing.T) {
		ctx = gdtcontext.PushTrace(ctx, s.Title())
		defer func() {
			ctx = gdtcontext.PopTrace(ctx)
		}()
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
				debug.Println(ctx, "wait: %s before", wait.Before)
				time.Sleep(wait.BeforeDuration())
			}

			to := specTimeout(ctx, t, sb.Timeout, scDefaults)
			if to != nil {
				var cancel context.CancelFunc
				specCtx, cancel = context.WithTimeout(specCtx, to.Duration())
				defer cancel()
			}
			res := spec.Eval(specCtx, t)
			if res.HasRuntimeError() {
				rterr = res.RuntimeError()
				t.Fatal(rterr)
				break
			}
			// Results can have arbitrary run data stored in them and we
			// save this prior run data in the top-level context (and pass
			// that context to the next Run invocation).
			if res.HasData() {
				ctx = gdtcontext.StorePriorRun(ctx, res.Data())
			}
			if wait != nil && wait.After != "" {
				debug.Println(ctx, "wait: %s after", wait.After)
				time.Sleep(wait.AfterDuration())
			}
		}
	})
	return rterr
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
			ctx, "using timeout of %s (expected: %t)",
			specTimeout.After, specTimeout.Expected,
		)
		return specTimeout
	}
	if scenDefaults != nil && scenDefaults.Timeout != nil {
		debug.Println(
			ctx, "using timeout of %s (expected: %t) [scenario default]",
			scenDefaults.Timeout.After, scenDefaults.Timeout.Expected,
		)
		return scenDefaults.Timeout
	}
	return nil
}

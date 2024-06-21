// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/result"
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
			if err := fix.Start(ctx); err != nil {
				return err
			}
			defer fix.Stop(ctx)
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
		res, err := skipIf.Eval(ctx)
		if err != nil {
			return err
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
		for idx, spec := range s.Tests {
			sb := spec.Base()
			wait := sb.Wait
			if wait != nil && wait.Before != "" {
				debug.Println(ctx, "wait: %s before", wait.Before)
				time.Sleep(wait.BeforeDuration())
			}
			plugin := s.evalPlugins[idx]
			pinfo := plugin.Info()
			pretry := pinfo.Retry
			ptimeout := pinfo.Timeout

			rt := getRetry(ctx, sb.Retry, scDefaults, pretry)

			// Create a brand new context that inherits the top-level context's
			// cancel func. We want to set deadlines for each test spec and if
			// we mutate the single supplied top-level context, then only the
			// first deadline/timeout will be used.
			specCtx, specCancel := context.WithCancel(ctx)
			defer specCancel()

			to := getTimeout(ctx, sb.Timeout, ptimeout, scDefaults)
			if to != nil {
				var cancel context.CancelFunc
				specCtx, cancel = context.WithTimeout(specCtx, to.Duration())
				defer cancel()
			}

			res, rterr := s.runSpec(specCtx, rt, to, idx, spec)
			if rterr != nil {
				break
			}
			if wait != nil && wait.After != "" {
				debug.Println(ctx, "wait: %s after", wait.After)
				time.Sleep(wait.AfterDuration())
			}
			// Results can have arbitrary run data stored in them and we
			// save this prior run data in the top-level context (and pass
			// that context to the next Run invocation).
			if res.HasData() {
				ctx = gdtcontext.StorePriorRun(ctx, res.Data())
			}
			for _, fail := range res.Failures() {
				t.Error(fail)
			}
		}
	})
	return rterr
}

// runSpec executes an individual test spec
func (s *Scenario) runSpec(
	ctx context.Context,
	retry *gdttypes.Retry,
	timeout *gdttypes.Timeout,
	idx int,
	spec gdttypes.Evaluable,
) (*result.Result, error) {
	sb := spec.Base()
	specTraceMsg := strconv.Itoa(idx)
	if sb.Name != "" {
		specTraceMsg += ":" + sb.Name
	}
	ctx = gdtcontext.PushTrace(ctx, specTraceMsg)
	defer func() {
		ctx = gdtcontext.PopTrace(ctx)
	}()
	if retry == nil {
		// Just evaluate the test spec once
		res, err := spec.Eval(ctx)
		if err != nil {
			return nil, err
		}
		debug.Println(
			ctx, "run: single-shot (no retries) ok: %v",
			!res.Failed(),
		)
		return res, nil
	}

	// retry the action and test the assertions until they succeed,
	// there is a terminal failure, or the timeout expires.
	var bo backoff.BackOff
	var res *result.Result
	var err error

	if retry.Exponential {
		bo = backoff.WithContext(
			backoff.NewExponentialBackOff(),
			ctx,
		)
	} else {
		interval := gdttypes.DefaultRetryConstantInterval
		if retry.Interval != "" {
			interval = retry.IntervalDuration()
		}
		bo = backoff.WithContext(
			backoff.NewConstantBackOff(interval),
			ctx,
		)
	}
	ticker := backoff.NewTicker(bo)
	maxAttempts := 0
	if retry.Attempts != nil {
		maxAttempts = *retry.Attempts
	}
	attempts := 1
	start := time.Now().UTC()
	success := false
	for tick := range ticker.C {
		if (maxAttempts > 0) && (attempts > maxAttempts) {
			debug.Println(
				ctx, "run: exceeded max attempts %d. stopping.",
				maxAttempts,
			)
			ticker.Stop()
			break
		}
		after := tick.Sub(start)

		res, err = spec.Eval(ctx)
		if err != nil {
			return nil, err
		}
		success = !res.Failed()
		debug.Println(
			ctx, "run: attempt %d after %s ok: %v",
			attempts, after, success,
		)
		if success {
			ticker.Stop()
			break
		}
		for _, f := range res.Failures() {
			debug.Println(
				ctx, "run: attempt %d after %s failure: %s",
				attempts, after, f,
			)
		}
		attempts++
	}
	return res, nil
}

// getTimeout returns the timeout value for the test spec. If the spec has a
// timeout override, we use that. Otherwise, we inspect the scenario's defaults
// and, if present, use that timeout. If the scenario's defaults for not
// indicate a timeout configuration, we ask the plugin if it has timeout
// defaults and use that.
func getTimeout(
	ctx context.Context,
	specTimeout *gdttypes.Timeout,
	pluginTimeout *gdttypes.Timeout,
	scenDefaults *Defaults,
) *gdttypes.Timeout {
	if specTimeout != nil {
		debug.Println(
			ctx, "using timeout of %s",
			specTimeout.After,
		)
		return specTimeout
	}
	if scenDefaults != nil && scenDefaults.Timeout != nil {
		debug.Println(
			ctx, "using timeout of %s [scenario default]",
			scenDefaults.Timeout.After,
		)
		return scenDefaults.Timeout
	}
	if pluginTimeout != nil {
		debug.Println(
			ctx, "using timeout of %s [plugin default]",
			pluginTimeout.After,
		)
		return pluginTimeout
	}
	return nil
}

// getRetry returns the retry configuration for the test spec. If the spec has a
// retry override, we use that. Otherwise, we inspect the scenario's defaults
// and, if present, use that timeout. If the scenario's defaults do not
// indicate a retry configuration, we ask the plugin if it has retry defaults
// and use that.
func getRetry(
	ctx context.Context,
	specRetry *gdttypes.Retry,
	scenDefaults *Defaults,
	pluginRetry *gdttypes.Retry,
) *gdttypes.Retry {
	if specRetry != nil {
		msg := "using retry"
		if specRetry.Attempts != nil {
			msg += fmt.Sprintf(" (attempts: %d)", *specRetry.Attempts)
		}
		if specRetry.Interval != "" {
			msg += fmt.Sprintf(" (interval: %s)", specRetry.Interval)
		}
		msg += fmt.Sprintf(" (exponential: %t)", specRetry.Exponential)
		debug.Println(ctx, msg)
		return specRetry
	}
	if scenDefaults != nil && scenDefaults.Retry != nil {
		scenRetry := scenDefaults.Retry
		msg := "using retry"
		if scenRetry.Attempts != nil {
			msg += fmt.Sprintf(" (attempts: %d)", *scenRetry.Attempts)
		}
		if scenRetry.Interval != "" {
			msg += fmt.Sprintf(" (interval: %s)", scenRetry.Interval)
		}
		msg += fmt.Sprintf(" (exponential: %t) [scenario default]", scenRetry.Exponential)
		debug.Println(ctx, msg)
		return scenRetry
	}
	if pluginRetry != nil {
		msg := "using retry"
		if pluginRetry.Attempts != nil {
			msg += fmt.Sprintf(" (attempts: %d)", *pluginRetry.Attempts)
		}
		if pluginRetry.Interval != "" {
			msg += fmt.Sprintf(" (interval: %s)", pluginRetry.Interval)
		}
		msg += fmt.Sprintf(" (exponential: %t) [plugin default]", pluginRetry.Exponential)
		debug.Println(ctx, msg)
		return pluginRetry
	}
	return nil
}

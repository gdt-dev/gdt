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

	"github.com/gdt-dev/gdt/api"
	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
)

// Run executes the scenario. The error that is returned will always be derived
// from `api.RuntimeError` and represents an *unrecoverable* error.
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
				return api.RequiredFixtureMissing(fname)
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

			// Create a brand new context that inherits the top-level context's
			// cancel func. We want to set deadlines for each test spec and if
			// we mutate the single supplied top-level context, then only the
			// first deadline/timeout will be used.
			specCtx, specCancel := context.WithCancel(ctx)

			specTraceMsg := strconv.Itoa(idx)
			if sb.Name != "" {
				specTraceMsg += ":" + sb.Name
			}
			specCtx = gdtcontext.PushTrace(specCtx, specTraceMsg)
			popTracer := func() {
				specCtx = gdtcontext.PopTrace(specCtx)
			}

			plugin := s.evalPlugins[idx]

			rt := getRetry(specCtx, scDefaults, plugin, spec)

			to := getTimeout(specCtx, scDefaults, plugin, spec)

			var res *api.Result
			ch := make(chan runSpecRes, 1)

			wait := sb.Wait
			if wait != nil && wait.Before != "" {
				debug.Println(specCtx, "wait: %s before", wait.Before)
				time.Sleep(wait.BeforeDuration())
			}

			if to != nil {
				specCtx, specCancel = context.WithTimeout(specCtx, to.Duration())
			}

			go s.runSpec(specCtx, ch, rt, idx, spec)

			select {
			case <-specCtx.Done():
				t.Fatalf("assertion failed: timeout exceeded (%s)", to.After)
				popTracer()
				specCancel()
				break
			case runres := <-ch:
				res = runres.r
				rterr = runres.err
			}
			if rterr != nil {
				popTracer()
				specCancel()
				break
			}

			if wait != nil && wait.After != "" {
				debug.Println(specCtx, "wait: %s after", wait.After)
				time.Sleep(wait.AfterDuration())
			}

			// Results can have arbitrary run data stored in them and we
			// save this prior run data in the top-level context (and pass
			// that context to the next Run invocation).
			if res.HasData() {
				ctx = gdtcontext.StorePriorRun(ctx, res.Data())
			}
			for _, fail := range res.Failures() {
				t.Fatal(fail)
			}
			popTracer()
			specCancel()
		}
	})
	return rterr
}

type runSpecRes struct {
	r   *api.Result
	err error
}

// runSpec executes an individual test spec
func (s *Scenario) runSpec(
	ctx context.Context,
	ch chan runSpecRes,
	retry *api.Retry,
	idx int,
	spec api.Evaluable,
) {
	if retry == nil || retry == api.NoRetry {
		// Just evaluate the test spec once
		res, err := spec.Eval(ctx)
		if err != nil {
			ch <- runSpecRes{nil, err}
			return
		}
		debug.Println(
			ctx, "run: single-shot (no retries) ok: %v",
			!res.Failed(),
		)
		ch <- runSpecRes{res, nil}
		return
	}

	// retry the action and test the assertions until they succeed,
	// there is a terminal failure, or the timeout expires.
	var bo backoff.BackOff
	var res *api.Result
	var err error

	if retry.Exponential {
		bo = backoff.WithContext(
			backoff.NewExponentialBackOff(),
			ctx,
		)
	} else {
		interval := api.DefaultRetryConstantInterval
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
			ch <- runSpecRes{nil, err}
			return
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
				ctx, "run: attempt %d failure: %s",
				attempts, f,
			)
		}
		attempts++
	}
	ch <- runSpecRes{res, nil}
}

// getTimeout returns the timeout configuration for the test spec. We check for
// overrides in timeout configuration using the following precedence:
//
// * Spec (Evaluable) override
// * Spec's Base override
// * Scenario's default
// * Plugin's default
func getTimeout(
	ctx context.Context,
	scenDefaults *Defaults,
	plugin api.Plugin,
	eval api.Evaluable,
) *api.Timeout {
	evalTimeout := eval.Timeout()
	if evalTimeout != nil {
		debug.Println(
			ctx, "using timeout of %s",
			evalTimeout.After,
		)
		return evalTimeout
	}

	sb := eval.Base()
	baseTimeout := sb.Timeout
	if baseTimeout != nil {
		debug.Println(
			ctx, "using timeout of %s",
			baseTimeout.After,
		)
		return baseTimeout
	}

	if scenDefaults != nil && scenDefaults.Timeout != nil {
		debug.Println(
			ctx, "using timeout of %s [scenario default]",
			scenDefaults.Timeout.After,
		)
		return scenDefaults.Timeout
	}

	pluginInfo := plugin.Info()
	pluginTimeout := pluginInfo.Timeout

	if pluginTimeout != nil {
		debug.Println(
			ctx, "using timeout of %s [plugin default]",
			pluginTimeout.After,
		)
		return pluginTimeout
	}
	return nil
}

// getRetry returns the retry configuration for the test spec. We check for
// overrides in retry configuration using the following precedence:
//
// * Spec (Evaluable) override
// * Spec's Base override
// * Scenario's default
// * Plugin's default
func getRetry(
	ctx context.Context,
	scenDefaults *Defaults,
	plugin api.Plugin,
	eval api.Evaluable,
) *api.Retry {
	evalRetry := eval.Retry()
	if evalRetry != nil {
		if evalRetry == api.NoRetry {
			return evalRetry
		}
		msg := "using retry"
		if evalRetry.Attempts != nil {
			msg += fmt.Sprintf(" (attempts: %d)", *evalRetry.Attempts)
		}
		if evalRetry.Interval != "" {
			msg += fmt.Sprintf(" (interval: %s)", evalRetry.Interval)
		}
		msg += fmt.Sprintf(" (exponential: %t)", evalRetry.Exponential)
		debug.Println(ctx, msg)
		return evalRetry
	}

	sb := eval.Base()
	baseRetry := sb.Retry
	if baseRetry != nil {
		if baseRetry == api.NoRetry {
			return baseRetry
		}
		msg := "using retry"
		if baseRetry.Attempts != nil {
			msg += fmt.Sprintf(" (attempts: %d)", *baseRetry.Attempts)
		}
		if baseRetry.Interval != "" {
			msg += fmt.Sprintf(" (interval: %s)", baseRetry.Interval)
		}
		msg += fmt.Sprintf(" (exponential: %t)", baseRetry.Exponential)
		debug.Println(ctx, msg)
		return baseRetry
	}

	if scenDefaults != nil && scenDefaults.Retry != nil {
		scenRetry := scenDefaults.Retry
		if scenRetry == api.NoRetry {
			return scenRetry
		}
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

	pluginInfo := plugin.Info()
	pluginRetry := pluginInfo.Retry

	if pluginRetry != nil {
		if pluginRetry == api.NoRetry {
			return pluginRetry
		}
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

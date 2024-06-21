// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package result

// Result is returned from a `Evaluable.Eval` execution. It serves two
// purposes:
//
// 1) to return an error, if any, from the Eval execution. This error will
// *always* be a `gdterrors.RuntimeError`. Failed assertions are not errors.
// 2) to pass back information about the Run that can be injected into the
// context's `PriorRun` cache. Some plugins, e.g. the gdt-http plugin, use
// cached data from a previous run in order to construct current Run fixtures.
// In the case of the gdt=http plugin, the previous `nethttp.Response` is
// returned in the Result and the `Scenario.Run` method injects that
// information into the context that is supplied to the next Spec's `Run`.
type Result struct {
	// failures is the collection of error messages from assertion failures
	// that occurred during Eval(). These are *not* `gdterrors.RuntimeError`.
	failures []error
	// data is a map, keyed by plugin name, of data about the spec run. Plugins
	// can place anything they want in here and grab it from the context with
	// the `gdtcontext.PriorRunData()` function. Plugins are responsible for
	// clearing and setting any used prior run data.
	data map[string]interface{}
}

// HasData returns true if any of the run data has been set, false otherwise.
func (r *Result) HasData() bool {
	return r.data != nil
}

// Data returns the raw run data saved in the result
func (r *Result) Data() map[string]interface{} {
	return r.data
}

// Failed returns true if any assertion failed during Eval(), false otherwise.
func (r *Result) Failed() bool {
	return len(r.failures) > 0
}

// Failures returns the collection of assertion failures that occurred during
// Eval().
func (r *Result) Failures() []error {
	return r.failures
}

// SetData sets a value in the result's run data cache.
func (r *Result) SetData(
	key string,
	val interface{},
) {
	if r.data == nil {
		r.data = map[string]interface{}{}
	}
	r.data[key] = val
}

// SetFailures sets the result's collection of assertion failures.
func (r *Result) SetFailures(failures ...error) {
	r.failures = failures
}

type ResultModifier func(*Result)

// WithData modifies the Result with the supplied run data key and value
func WithData(key string, val interface{}) ResultModifier {
	return func(r *Result) {
		r.SetData(key, val)
	}
}

// WithFailures modifies the Result the supplied collection of assertion
// failures
func WithFailures(failures ...error) ResultModifier {
	return func(r *Result) {
		r.SetFailures(failures...)
	}
}

// New returns a new Result
func New(mods ...ResultModifier) *Result {
	r := &Result{}
	for _, mod := range mods {
		mod(r)
	}
	return r
}

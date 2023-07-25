// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package result

// Result is returned from a `Spec.Run` execution. It serves two purposes: 1)
// to return the error, if any, from the Run execution and 2) to pass back
// information about the Run that can be injected into the context's `PriorRun`
// cache. Some plugins, e.g. the gdt-http plugin, use cached data from a
// previous run in order to construct current Run fixtures. In the case of the
// gdt=http plugin, the previous `nethttp.Response` is returned in the Result
// and the `Scenario.Run` method injects that information into the context that
// is supplied to the next Spec's `Run`.
type Result struct {
	// err is any error that was returned from the spec's execution
	err error
	// data is a map, keyed by plugin name, of data about the spec run. Plugins
	// can place anything they want in here and grab it from the context with
	// the `gdtcontext.PriorRunData()` function. Plugins are responsible for
	// clearing and setting any used prior run data.
	data map[string]interface{}
}

// Unwrap returns the wrapped error
func (r *Result) Unwrap() error {
	return r.err
}

// Error implements the error interface.
func (r *Result) Error() string {
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

// HasData returns true if any of the run data has been set, false otherwise.
func (r *Result) HasData() bool {
	return r.data != nil
}

// Data returns the raw run data saved in the result
func (r *Result) Data() map[string]interface{} {
	return r.data
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

type ResultModifier func(*Result)

// WithError modifies the Result with the supplied error
func WithError(err error) ResultModifier {
	return func(r *Result) {
		r.err = err
	}
}

// WithData modifies the Result with the supplied run data key and value
func WithData(key string, val interface{}) ResultModifier {
	return func(r *Result) {
		r.SetData(key, val)
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

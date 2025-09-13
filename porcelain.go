// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"io"
	"os"

	"github.com/gdt-dev/core/api"
	gdtcontext "github.com/gdt-dev/core/context"
	jsonfix "github.com/gdt-dev/core/fixture/json"
	"github.com/gdt-dev/core/plugin"
	_ "github.com/gdt-dev/core/plugin/exec"
	"github.com/gdt-dev/core/scenario"
	"github.com/gdt-dev/core/suite"
)

var (
	// RegisterPlugin registers a plugin with gdt's set of known plugins.
	//
	// Generally only plugin authors will ever need to call this function. It is
	// not required for normal use of gdt or any known plugin.
	RegisterPlugin = plugin.Register
	// RegisterFixture registers a named fixtures with the context
	RegisterFixture = gdtcontext.RegisterFixture
	// NewContext returns a new `context.Context` that can be passed to a
	// `scenario.Run` or `suite.Run` pointer receivers.
	NewContext = gdtcontext.New
	// WithFixtures sets a context's Fixtures
	WithFixtures = gdtcontext.WithFixtures
	// WithDebug informs gdt to output extra debugging information. You can
	// supply zero or more `io.Writer` objects to the function.
	//
	// If no `io.Writer` objects are supplied, gdt will output debug messages
	// using the `testing.T.Log[f]()` function. This means that you will only
	// get these debug messages if you call the `go test` tool with the `-v`
	// option (either as `go test -v` or with `go test -v=test2json`.
	//
	// ```go
	//	func TestExample(t *testing.T) {
	//      require := require.New(t)
	//      fp := filepath.Join("testdata", "example.yaml")
	//      f, err := os.Open(fp)
	//      require.Nil(err)
	//
	//      ctx := gdt.NewContext(gdt.WithDebug())
	//      s, err := scenario.FromReader(
	//          f,
	//          scenario.WithPath(fp),
	//          scenario.WithContext(ctx),
	//      )
	//      require.Nil(err)
	//      require.NotNil(s)
	//
	//      err = s.Run(ctx, t)
	//      require.Nil(err)
	// }
	// ```
	//
	// If you want gdt to log extra debugging information about tests and
	// assertions to a different file or collecting buffer, pass it a context
	// with a debug `io.Writer`:
	//
	// ```go
	// f := ioutil.TempFile("", "mytest*.log")
	// ctx := gdt.NewContext(gdt.WithDebug(f))
	// ```
	//
	// ```go
	// var b bytes.Buffer
	// w := bufio.NewWriter(&b)
	// ctx := gdt.NewContext(gdt.WithDebug(w))
	// ```
	//
	// you can then inspect the debug "log" and do whatever you'd like with it.
	WithDebug = gdtcontext.WithDebug
	// SetDebug sets gdt's debug logging to the supplied `io.Writer`.
	//
	// The `writers` parameters is optional. If no `io.Writer` objects are
	// supplied, gdt will output debug messages using the `testing.T.Log[f]()`
	// function. This means that you will only get these debug messages if you
	// call the `go test` tool with the `-v` option (either as `go test -v` or
	// with `go test -v=test2json`.
	SetDebug = gdtcontext.SetDebug
	// NewJSONFixture takes a string, some bytes or an io.Reader and returns a
	// new api.Fixture that can have its state queried via JSONPath
	NewJSONFixture = jsonfix.New
)

// From returns a new `api.Runnable` from an `io.Reader`, a string file or
// directory path, or the raw bytes of YAML content describing a scenario or
// suite.
func From(source any) (api.Runnable, error) {
	switch src := source.(type) {
	case io.Reader:
		s, err := scenario.FromReader(src)
		if err != nil {
			return nil, err
		}
		return suite.FromScenario(s), nil
	case string:
		f, err := os.Open(src)
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			return suite.FromDir(src)
		}
		s, err := scenario.FromReader(f, scenario.WithPath(src))
		if err != nil {
			return nil, err
		}
		return suite.FromScenario(s), nil
	case []byte:
		s, err := scenario.FromBytes(src)
		if err != nil {
			return nil, err
		}
		return suite.FromScenario(s), nil
	default:
		return nil, api.UnknownSourceType(source)
	}
}

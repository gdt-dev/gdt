// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"
	"io"
	"os"

	"github.com/samber/lo"

	"github.com/gdt-dev/gdt/api"
)

type ContextKey string

var (
	debugKey    = ContextKey("gdt.debug")
	traceKey    = ContextKey("gdt.trace")
	pluginsKey  = ContextKey("gdt.plugins")
	fixturesKey = ContextKey("gdt.fixtures")
	priorRunKey = ContextKey("gdt.run.prior")
)

// ContextModifier sets some value on the context
type ContextModifier func(context.Context) context.Context

// WithDebug informs gdt to output extra debugging information. You can supply
// zero or more `io.Writer` objects to the function.
//
// If no `io.Writer` objects are supplied, gdt will output debug messages using
// the `fmt.Printf` function. The `fmt.Printf` function is *unbuffered* however
// unless you call `go test` with the `-v` argument, `go test` swallows output
// to stdout and does not display it unless a test fails.
//
// This means that you will only get these debug messages if you call the `go
// test` tool with the `-v` option (either as `go test -v` or with `go test
// -v=test2json`.
//
// ```go
//
//	func TestExample(t *testing.T) {
//		   require := require.New(t)
//		   fp := filepath.Join("testdata", "example.yaml")
//		   f, err := os.Open(fp)
//		   require.Nil(err)
//
//		   ctx := gdtcontext.New(gdtcontext.WithDebug())
//	       s, err := scenario.FromReader(
//	          f,
//	          scenario.WithPath(fp),
//	          scenario.WithContext(ctx),
//	       )
//		   require.Nil(err)
//		   require.NotNil(s)
//
//		   err = s.Run(ctx, t)
//		   require.Nil(err)
//	}
//
// ```
//
// If you want gdt to log extra debugging information about tests and
// assertions to a different file or collecting buffer, pass it a context with
// a debug `io.Writer`:
//
// ```go
// f := ioutil.TempFile("", "mytest*.log")
// ctx := gdtcontext.New(gdtcontext.WithDebug(f))
// ```
//
// ```go
// var b bytes.Buffer
// w := bufio.NewWriter(&b)
// ctx := gdtcontext.New(gdtcontext.WithDebug(w))
// ```
//
// you can then inspect the debug "log" and do whatever you'd like with it.
func WithDebug(writers ...io.Writer) ContextModifier {
	return func(ctx context.Context) context.Context {
		if len(writers) == 0 {
			// Write to stdout when WithDebug() is called with no parameters
			writers = []io.Writer{os.Stdout}
		}
		return context.WithValue(ctx, debugKey, writers)
	}
}

// WithPlugins sets a context's Plugins
func WithPlugins(plugins []api.Plugin) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pluginsKey, plugins)
	}
}

// WithFixtures sets a context's Fixtures
func WithFixtures(fixtures map[string]api.Fixture) ContextModifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, fixturesKey, fixtures)
	}
}

// SetDebug sets gdt's debug logging to the supplied `io.Writer`.
//
// The `writers` parameters is optional. If no `io.Writer` objects are
// supplied, gdt will output debug messages to stdout.
func SetDebug(
	ctx context.Context,
	writers ...io.Writer,
) context.Context {
	if len(writers) == 0 {
		// This triggers writes to stdout when WithDebug() is called with no
		// parameters...
		writers = []io.Writer{os.Stdout}
	}
	return context.WithValue(ctx, debugKey, writers)
}

// RegisterFixture registers a named fixtures with the context
func RegisterFixture(
	ctx context.Context,
	name string,
	f api.Fixture,
) context.Context {
	fixtures := Fixtures(ctx)
	fixtures[name] = f
	return context.WithValue(ctx, fixturesKey, fixtures)
}

// RegisterPlugin registers a plugin with the context
func RegisterPlugin(
	ctx context.Context,
	p api.Plugin,
) context.Context {
	plugins := Plugins(ctx)
	for _, plug := range plugins {
		if plug.Info().Name == p.Info().Name {
			// No need to register... already known.
			return ctx
		}
	}
	plugins = append(plugins, p)
	return context.WithValue(ctx, pluginsKey, plugins)
}

// StorePriorRun saves prior run data in the context. If there is already prior
// run data cached in the supplied context, the existing data is merged with
// the supplied data.
func StorePriorRun(
	ctx context.Context,
	data map[string]interface{},
) context.Context {
	existing := PriorRun(ctx)
	merged := lo.Assign(existing, data)
	return context.WithValue(ctx, priorRunKey, merged)
}

// PushTrace pushes a debug/trace name onto the debug/trace stack. It is used
// by plugins to track where in the processing of a test or assertion the
// plugin is and gets output at the start of a debug.Printf/Println message.
func PushTrace(
	ctx context.Context,
	name string,
) context.Context {
	stack := TraceStack(ctx)
	stack = append(stack, name)
	return context.WithValue(ctx, traceKey, stack)
}

// PopTrace pops the last name off the debug/trace stack. It is used by plugins
// to track where in the processing of a test or assertion the plugin is and
// gets output at the start of a debug.Printf/Println message.
func PopTrace(
	ctx context.Context,
) context.Context {
	stack := TraceStack(ctx)
	stack = stack[:len(stack)-1]
	return context.WithValue(ctx, traceKey, stack)
}

// New returns a new Context
func New(mods ...ContextModifier) context.Context {
	ctx := context.TODO()
	for _, mod := range mods {
		ctx = mod(ctx)
	}
	return ctx
}

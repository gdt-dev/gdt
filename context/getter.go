// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"
	"io"

	gdttypes "github.com/gdt-dev/gdt/types"
)

// Debug gets a context's Debug writer
func Debug(ctx context.Context) []io.Writer {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(debugKey); v != nil {
		return v.([]io.Writer)
	}
	return nil
}

// Plugins gets a context's Plugins
func Plugins(ctx context.Context) []gdttypes.Plugin {
	if ctx == nil {
		return []gdttypes.Plugin{}
	}
	if v := ctx.Value(pluginsKey); v != nil {
		return v.([]gdttypes.Plugin)
	}
	return []gdttypes.Plugin{}
}

// Fixtures gets a context's Fixtures
func Fixtures(ctx context.Context) map[string]gdttypes.Fixture {
	if ctx == nil {
		return map[string]gdttypes.Fixture{}
	}
	if v := ctx.Value(fixturesKey); v != nil {
		return v.(map[string]gdttypes.Fixture)
	}
	return map[string]gdttypes.Fixture{}
}

// PriorRun gets a context's prior run data
func PriorRun(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return map[string]interface{}{}
	}
	if v := ctx.Value(priorRunKey); v != nil {
		return v.(map[string]interface{})
	}
	return map[string]interface{}{}
}

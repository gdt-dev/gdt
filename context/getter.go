// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gdt-dev/gdt/api"
)

const (
	traceDelimiter = "/"
)

// Trace gets a context's trace name stack joined together with
func Trace(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(traceKey); v != nil {
		return strings.Join(v.([]string), traceDelimiter)
	}
	return ""
}

// TraceStack gets a context's trace name stack
func TraceStack(ctx context.Context) []string {
	if ctx == nil {
		return []string{}
	}
	if v := ctx.Value(traceKey); v != nil {
		return v.([]string)
	}
	return []string{}
}

// Debug gets a context's Debug writer
func Debug(ctx context.Context) []io.Writer {
	if ctx == nil {
		return []io.Writer{}
	}
	if v := ctx.Value(debugKey); v != nil {
		return v.([]io.Writer)
	}
	return []io.Writer{}
}

// Plugins gets a context's Plugins
func Plugins(ctx context.Context) []api.Plugin {
	if ctx == nil {
		return []api.Plugin{}
	}
	if v := ctx.Value(pluginsKey); v != nil {
		return v.([]api.Plugin)
	}
	return []api.Plugin{}
}

// Fixtures gets a context's Fixtures
func Fixtures(ctx context.Context) map[string]api.Fixture {
	if ctx == nil {
		return map[string]api.Fixture{}
	}
	if v := ctx.Value(fixturesKey); v != nil {
		return v.(map[string]api.Fixture)
	}
	return map[string]api.Fixture{}
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

// ReplaceVariables replaces all occurrences of any of the variables in the
// prior run data with their stored variable values
func ReplaceVariables(
	ctx context.Context,
	subject string,
) string {
	data := PriorRun(ctx)
	for dataKey, dataVal := range data {
		var dataValStr string
		switch dataVal := dataVal.(type) {
		case string:
			dataValStr = dataVal
		case []byte:
			dataValStr = string(dataVal)
		case int, uint, int8, int16, int32, int64:
			dataValStr = strconv.Itoa(dataVal.(int))
		case float32, float64:
			dataValStr = strconv.FormatFloat(dataVal.(float64), 'f', -1, 64)
		default:
			continue
		}
		subject = strings.ReplaceAll(
			subject,
			fmt.Sprintf("$%s", dataKey),
			dataValStr,
		)
	}
	return subject
}

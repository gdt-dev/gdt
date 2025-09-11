// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	gdtcontext "github.com/gdt-dev/gdt/context"
)

const (
	varFromStdout = "stdout"
	varFromStderr = "stderr"
	varFromRC     = "returncode"
)

type VarEntry struct {
	// From is a string that indicates where the value of the variable will be
	// sourced from. `stdout`, `stderr` and `returncode` indicate to source the
	// value of the variable from the output buffer for stdout, stderr or the
	// returncode value. All other strings indicate the value of the variable
	// should be sourced from an envvar of the same name.
	From string `yaml:"from"`
}

// Variables allows the test author to save arbitrary data to the test scenario,
// facilitating the passing of variables between test specs potentially
// provided by different gdt Plugins.
type Variables map[string]VarEntry

// Replace replaces all occurrences of any of the variables in the supplied
// string with their stored variable values
func (v Variables) Replace(
	ctx context.Context,
	subject string,
) string {
	data := gdtcontext.PriorRun(ctx)
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

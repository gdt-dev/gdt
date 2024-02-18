// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package debug

import (
	"context"
	"fmt"
	"strings"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
)

// Printf writes a message with optional message arguments to the context's
// Debug output.
func Printf(
	ctx context.Context,
	t *testing.T,
	format string,
	args ...interface{},
) {
	t.Helper()
	writers := gdtcontext.Debug(ctx)
	if len(writers) == 0 {
		return
	}
	t.Logf(format, args...)

	if !strings.HasPrefix(format, "[gdt] ") {
		format = "[gdt] " + t.Name() + " " + format
	}
	msg := fmt.Sprintf(format, args...)
	for _, w := range writers {
		w.Write([]byte(msg))
	}
}

// Println writes a message with optional message arguments to the context's
// Debug output, ensuring there is a newline in the message line.
func Println(
	ctx context.Context,
	t *testing.T,
	format string,
	args ...interface{},
) {
	t.Helper()
	writers := gdtcontext.Debug(ctx)
	if len(writers) == 0 {
		return
	}
	// NOTE(jaypipes): T.Logf() automatically adds newlines...
	t.Logf(format, args...)

	if !strings.HasPrefix(format, "[gdt] ") {
		format = "[gdt] " + t.Name() + " " + format
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	msg := fmt.Sprintf(format, args...)
	for _, w := range writers {
		w.Write([]byte(msg))
	}
}

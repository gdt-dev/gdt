// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package debug

import (
	"context"
	"fmt"
	"strings"

	gdtcontext "github.com/gdt-dev/gdt/context"
)

// Printf writes a message with optional message arguments to the context's
// Debug output.
func Printf(
	ctx context.Context,
	format string,
	args ...interface{},
) {
	writers := gdtcontext.Debug(ctx)
	if len(writers) == 0 {
		return
	}

	trace := gdtcontext.Trace(ctx)

	if !strings.HasPrefix(format, "[gdt] ") {
		format = "[gdt] [" + trace + "] " + format
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
	format string,
	args ...interface{},
) {
	writers := gdtcontext.Debug(ctx)
	if len(writers) == 0 {
		return
	}

	trace := gdtcontext.Trace(ctx)

	if !strings.HasPrefix(format, "[gdt] ") {
		format = "[gdt] [" + trace + "] " + format
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	msg := fmt.Sprintf(format, args...)
	for _, w := range writers {
		w.Write([]byte(msg))
	}
}

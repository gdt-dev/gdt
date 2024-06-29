// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errstarter

import (
	"context"
	"fmt"

	"github.com/gdt-dev/gdt/fixture"
)

var (
	errStarter = func(_ context.Context) error {
		return fmt.Errorf("error starting fixture!")
	}

	Fixture = fixture.New(
		fixture.WithStarter(errStarter),
	)
)

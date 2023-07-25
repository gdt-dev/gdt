// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors_test

import (
	"testing"

	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/stretchr/testify/assert"
)

func TestUnknownSourceType(t *testing.T) {
	assert := assert.New(t)

	err := gdterrors.UnknownSourceType(1)
	assert.ErrorContains(err, "int")

	source := []string{"foo", "bar"}
	err = gdterrors.UnknownSourceType(source)
	assert.ErrorContains(err, "[]string")
}

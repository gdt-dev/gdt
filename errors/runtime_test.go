// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors_test

import (
	"testing"

	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeErrorsHas(t *testing.T) {
	assert := assert.New(t)

	re := gdterrors.NewRuntimeErrors()
	re.AppendIf(gdterrors.RequiredFixtureMissing("fixture"))

	assert.True(re.Has(gdterrors.ErrRequiredFixture))
}

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite_test

import (
	"context"
	"testing"

	"github.com/gdt-dev/gdt/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunExecSuite(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	s, err := suite.FromDir("testdata/exec")
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	assert.Nil(err)
}

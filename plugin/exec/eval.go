// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"context"
	"testing"

	"github.com/gdt-dev/gdt/debug"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/result"
)

// Eval performs an action and evaluates the results of that action, returning
// a Result that informs the Scenario about what failed or succeeded about the
// Evaluable's conditions.
func (s *Spec) Eval(ctx context.Context, t *testing.T) *result.Result {
	outbuf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}

	var ec int

	if err := s.Do(ctx, t, outbuf, errbuf, &ec); err != nil {
		if err == gdterrors.ErrTimeoutExceeded {
			return result.New(result.WithFailures(gdterrors.ErrTimeoutExceeded))
		}
		return result.New(result.WithRuntimeError(ExecRuntimeError(err)))
	}
	a := newAssertions(s.Assert, ec, outbuf, errbuf)
	if !a.OK() {
		for _, fail := range a.Failures() {
			t.Error(fail)
		}
		if s.On != nil {
			if s.On.Fail != nil {
				outbuf.Reset()
				errbuf.Reset()
				err := s.On.Fail.Do(ctx, t, outbuf, errbuf, nil)
				if err != nil {
					debug.Println(ctx, t, "error in on.fail.exec: %s", err)
				}
			}
		}
	}
	return result.New(result.WithFailures(a.Failures()...))
}

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"context"

	"github.com/gdt-dev/gdt/api"
	"github.com/gdt-dev/gdt/debug"
)

// Eval performs an action and evaluates the results of that action, returning
// a Result that informs the Scenario about what failed or succeeded about the
// Evaluable's conditions.
//
// Errors returned by Eval() are **RuntimeErrors**, not failures in assertions.
func (s *Spec) Eval(
	ctx context.Context,
) (*api.Result, error) {
	outbuf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}

	var ec int

	if err := s.Do(ctx, outbuf, errbuf, &ec); err != nil {
		if err == api.ErrTimeoutExceeded {
			return api.NewResult(api.WithFailures(api.ErrTimeoutExceeded)), nil
		}
		return nil, ExecRuntimeError(err)
	}
	a := newAssertions(s.Assert, ec, outbuf, errbuf)
	if !a.OK(ctx) {
		if s.On != nil {
			if s.On.Fail != nil {
				outbuf.Reset()
				errbuf.Reset()
				err := s.On.Fail.Do(ctx, outbuf, errbuf, nil)
				if err != nil {
					debug.Println(ctx, "error in on.fail.exec: %s", err)
				}
			}
		}
	}
	return api.NewResult(api.WithFailures(a.Failures()...)), nil
}

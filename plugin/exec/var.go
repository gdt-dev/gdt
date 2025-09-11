// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

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

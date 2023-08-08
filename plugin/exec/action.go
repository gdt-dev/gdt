// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

// Action describes a single execution of one or more commands via the
// operating system's `exec` family of functions.
type Action struct {
	// Exec is the exact command to execute.
	//
	// You may execute more than one command but must include the `shell` field
	// to indicate that the command should be run in a shell. It is best
	// practice, however, to simply use multiple `exec` specs instead of
	// executing multiple commands in a single shell call.
	Exec string `yaml:"exec"`
	// Shell is the specific shell to use in executing the command. If empty
	// (the default), no shell is used to execute the command and instead the
	// operating system's `exec` family of calls is used.
	Shell string `yaml:"shell,omitempty"`
}

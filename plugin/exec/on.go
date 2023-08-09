// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

// On describes actions that can be taken upon certain conditions.
type On struct {
	// Fail contains one or more actions to take if any of a Spec's assertions
	// fail.
	//
	// For example, if you wanted to grep a log file in the event that no
	// connectivity on a particular IP:PORT combination could be made you might
	// do this:
	//
	// ```yaml
	// tests:
	//  - exec: nc -z $HOST $PORT
	//    on:
	//      fail:
	//        exec: grep ERROR /var/log/myapp.log
	// ```
	//
	// The `grep ERROR /var/log/myapp.log` command will only be executed if
	// there is no connectivity to $HOST:$PORT and the results of that grep
	// will be directed to the test's output. You can use the `gdt.WithDebug()`
	// function to configure additional `io.Writer`s to direct this output to.
	Fail *Action `yaml:"fail,omitempty"`
}

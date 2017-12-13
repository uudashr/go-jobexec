package jobexec_test

import "github.com/uudashr/go-jobexec"

func ExampleExec() {
	var ctx jobexec.Context
	// TODO: initialize the ctx

	exec := func() error {
		// TODO: do something
		return nil
	}

	foo := jobexec.ExecutorID("foo")
	if err := jobexec.Exec(ctx, foo, exec); err != jobexec.ErrNotPermitted {
		// TODO: handle the error
	}
}

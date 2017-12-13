package jobexec_test

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/uudashr/go-jobexec"
)

func TestInMemContext(t *testing.T) {
	execInterval := 5 * time.Second

	foo := jobexec.ExecutorID("foo")
	bar := jobexec.ExecutorID("bar")

	clock := clockwork.NewFakeClockAt(time.Now())
	ctx, err := jobexec.NewInMemContextWithClock(execInterval, clock)
	if err != nil {
		t.Fatal("should success instantiating in-mem context:", err)
	}

	// foo is the first excutor, should be able to acquire permission succesfully
	if got, want := ctx.LastExecutor(), jobexec.EmptyExecutorID; got != want {
		t.Error("got:", got, "want:", want)
	}

	if err = ctx.Acquire(foo); err != nil {
		t.Fatal("should success acquire:", err)
	}

	if got, want := ctx.LastExecutor(), foo; got != want {
		t.Error("got:", got, "want:", want)
	}

	// bar is should not allowed, it require to wait 5 seconds after last acquire
	if err = ctx.Acquire(bar); err == nil {
		t.Fatal("expect to fail acquire")
	}

	// foo should be granted event without having to wait 5 seconds
	if err = ctx.Acquire(foo); err != nil {
		t.Fatal("should success acquire:", err)
	}

	// bar should be able to acquire after more than 5 seconds
	clock.Advance(execInterval + 1*time.Millisecond)
	if err = ctx.Acquire(bar); err != nil {
		t.Fatal("Should succeed acquire")
	}

	if got, want := ctx.LastExecutor(), bar; got != want {
		t.Error("got:", got, "want:", want)
	}
}

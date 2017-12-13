package mysql_test

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/uudashr/go-jobexec"
)

func TestContext(t *testing.T) {
	if testing.Short() {
		t.Skip("Require non-short mode to run test")
	}

	suite := Setup(t)
	defer suite.TearDown()

	execInterval := 2 * time.Second

	foo := jobexec.ExecutorID("foo")
	bar := jobexec.ExecutorID("bar")

	clock := clockwork.NewFakeClockAt(time.Now())
	ctx, err := jobexec.NewSQLContextWithClock(suite.DB, "DefaultCron", execInterval, clock)
	if err != nil {
		t.Fatal("err:", err)
	}

	if err = ctx.Init(); err != nil {
		t.Fatal("err:", err)
	}

	// foo is the first excutor, should be able to acquire permission succesfully
	lastExecutor, err := ctx.LastExecutor()
	if err != nil {
		t.Fatal("err:", err)
	}

	if got, want := lastExecutor, jobexec.EmptyExecutorID; got != want {
		t.Error("got:", got, "want:", want)
	}

	if err = ctx.Acquire(foo); err != nil {
		t.Fatal("should success acquire:", err)
	}

	lastExecutor, err = ctx.LastExecutor()
	if err != nil {
		t.Fatal("err:", err)
	}

	if got, want := lastExecutor, foo; got != want {
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

	lastExecutor, err = ctx.LastExecutor()
	if err != nil {
		t.Fatal("err:", err)
	}

	if got, want := lastExecutor, bar; got != want {
		t.Error("got:", got, "want:", want)
	}
}

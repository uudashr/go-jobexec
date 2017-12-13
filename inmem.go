package jobexec

import (
	"errors"
	"time"

	"github.com/jonboulle/clockwork"
)

// InMemContext is the in memory context.
type InMemContext struct {
	interval      time.Duration
	nextExecution time.Time
	lastExecutor  ExecutorID

	clock clockwork.Clock
}

// Acquire an execution permission.
func (c *InMemContext) Acquire(execID ExecutorID) error {
	now := c.clock.Now()
	if c.lastExecutor != EmptyExecutorID && c.lastExecutor != execID && !now.After(c.nextExecution) {
		return ErrNotPermitted
	}

	c.lastExecutor = execID
	c.nextExecution = now.Add(c.interval)
	return nil
}

// LastExecutor permitted to do an execution.
func (c *InMemContext) LastExecutor() ExecutorID {
	return c.lastExecutor
}

// NewInMemContext constructs new InMemContext.
func NewInMemContext(interval time.Duration) (*InMemContext, error) {
	return NewInMemContextWithClock(interval, clockwork.NewRealClock())
}

// NewInMemContextWithClock constructs new InMemContext with clock.
func NewInMemContextWithClock(interval time.Duration, clock clockwork.Clock) (*InMemContext, error) {
	if interval == 0 {
		return nil, errors.New("0 interval not allowed")
	}

	return &InMemContext{
		interval: interval,
		clock:    clock,
	}, nil
}

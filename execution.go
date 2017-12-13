package jobexec

// Exec a function based on given context.
func Exec(ctx Context, executorID ExecutorID, fn func() error) error {
	if err := ctx.Acquire(executorID); err != nil {
		return err
	}

	return fn()
}

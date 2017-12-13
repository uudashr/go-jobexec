# go-jobexec

jobexec provides mechanism to control permission for distributed job execution by providing context.

## Example
```go
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
```

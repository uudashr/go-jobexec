package jobexec

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// ErrNotPermitted is error when the acquire is not permitted.
var ErrNotPermitted = errors.New("jobexec: not permitted")

// Context of the job execution.
type Context interface {
	Acquire(ExecutorID) error
}

// EmptyExecutorID is empty ExecutorID
const EmptyExecutorID ExecutorID = ExecutorID("")

// ExecutorID is the executor indentifier
type ExecutorID string

// OK check the validity.
func (id ExecutorID) OK() bool {
	return strings.TrimSpace(string(id)) != ""
}

// Scan implements the sql.Scanner interface.
func (id *ExecutorID) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		*id = ExecutorID(string(v))
		return nil
	case string:
		*id = ExecutorID(v)
		return nil
	default:
		return fmt.Errorf("jobexec: unable to scan executor id %#v", v)
	}
}

// Value implements the driver.Value interface.
func (id ExecutorID) Value() (driver.Value, error) {
	return string(id), nil
}
